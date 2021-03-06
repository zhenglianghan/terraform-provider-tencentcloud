package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

func TencentMysqlSellType() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cdb_type": {
			Type:     schema.TypeString,
			Computed: true},
		"mem_size": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"min_volume_size": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"max_volume_size": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"volume_step": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"qps": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func TencentMysqlZoneConfig() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"is_default": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"is_support_disaster_recovery": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"is_support_vpc": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"engine_versions": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
		"pay_type": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeInt},
			Computed: true,
		},
		"hour_instance_sale_max_num": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"support_slave_sync_modes": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeInt},
			Computed: true,
		},
		"disaster_recovery_zones": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
		"slave_deploy_modes": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeInt},
			Computed: true,
		},
		"first_slave_zones": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
		"second_slave_zones": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
		"sells": {Type: schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: TencentMysqlSellType(),
			},
		},
	}
}

func dataSourceTencentMysqlZoneConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentMysqlZoneConfigRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(connectivity.MysqlSupportedRegions),
			},
			"result_output_file": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},

			// Computed values
			"list": {Type: schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: TencentMysqlZoneConfig(),
				},
			},
		},
	}
}

func dataSourceTencentMysqlZoneConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer LogElapsed("data_source.tencentcloud_mysql_zone_config.read")()

	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	region := meta.(*TencentCloudClient).apiV3Conn.Region
	if regionInterface, ok := d.GetOk("region"); ok {
		region = regionInterface.(string)
	} else {
		log.Printf("[INFO]%s region is not set,so we use [%s] from env\n ", logId, region)
	}

	sellConfigures, err := mysqlService.DescribeDBZoneConfig(ctx)
	if err != nil {
		return fmt.Errorf("api[DescribeBackups]fail, return %s", err.Error())
	}
	var regionItem *cdb.RegionSellConf
	for _, regionItem = range sellConfigures {
		if *regionItem.Region == region {
			break
		}
	}
	if regionItem == nil {
		return nil
	}
	var zoneConfigs []interface{}
	for _, sellItem := range regionItem.ZonesConf {
		if *sellItem.Status != ZONE_SELL_STATUS_ONLINE && *sellItem.Status != ZONE_SELL_STATUS_NEW {
			continue
		}
		var zoneConfig = make(map[string]interface{})
		zoneConfig["name"] = *sellItem.Zone
		zoneConfig["hour_instance_sale_max_num"] = *sellItem.HourInstanceSaleMaxNum

		if *sellItem.IsDefaultZone {
			zoneConfig["is_default"] = 1
		} else {
			zoneConfig["is_default"] = 0
		}

		if *sellItem.IsSupportDr {
			zoneConfig["is_support_disaster_recovery"] = 1
		} else {
			zoneConfig["is_support_disaster_recovery"] = 0
		}

		if *sellItem.IsSupportVpc {
			zoneConfig["is_support_vpc"] = 1
		} else {
			zoneConfig["is_support_vpc"] = 0
		}

		payTypes := make([]int, len(sellItem.PayType))
		for index, strPtr := range sellItem.PayType {
			if tempInt, err := strconv.ParseInt(*strPtr, 10, 64); err != nil {
				errRet := fmt.Errorf("api[DescribeDBZoneConfig]return PayType error,not int")
				log.Printf("[CRITAL]%s %s\n ", logId, errRet.Error())
				return errRet
			} else {
				payTypes[index] = int(tempInt)
			}
		}
		zoneConfig["pay_type"] = payTypes

		supportSlaveSyncModes := make([]string, len(sellItem.ProtectMode))
		for index, intPtr := range sellItem.ProtectMode {
			supportSlaveSyncModes[index] = *intPtr
		}
		zoneConfig["support_slave_sync_modes"] = payTypes

		disasterRecoveryZones := make([]string, len(sellItem.DrZone))
		for index, strPtr := range sellItem.DrZone {
			disasterRecoveryZones[index] = *strPtr
		}
		zoneConfig["disaster_recovery_zones"] = disasterRecoveryZones

		var (
			slaveDeployModes                                  []int
			firstSlaveZones, secondSlaveZones, engineVersions []string
			sells                                             []interface{}
		)
		if sellItem.ZoneConf != nil {
			for _, mode := range sellItem.ZoneConf.DeployMode {
				slaveDeployModes = append(slaveDeployModes, int(*mode))
			}
			for _, zoneName := range sellItem.ZoneConf.SlaveZone {
				firstSlaveZones = append(firstSlaveZones, *zoneName)
			}
			for _, zoneName := range sellItem.ZoneConf.BackupZone {
				secondSlaveZones = append(secondSlaveZones, *zoneName)
			}
		}
		zoneConfig["slave_deploy_modes"] = slaveDeployModes
		zoneConfig["first_slave_zones"] = firstSlaveZones
		zoneConfig["second_slave_zones"] = secondSlaveZones

		for _, mysqlConfigs := range sellItem.SellType {
			for _, strPtr := range mysqlConfigs.EngineVersion {
				engineVersions = append(engineVersions, *strPtr)
			}
			for _, mysqlConfig := range mysqlConfigs.Configs {
				var showConfigMap = make(map[string]interface{})
				showConfigMap["cdb_type"] = *mysqlConfig.CdbType
				showConfigMap["mem_size"] = int(*mysqlConfig.Memory)
				showConfigMap["max_volume_size"] = int(*mysqlConfig.VolumeMax)
				showConfigMap["min_volume_size"] = int(*mysqlConfig.VolumeMin)
				showConfigMap["volume_step"] = int(*mysqlConfig.VolumeStep)
				showConfigMap["qps"] = int(*mysqlConfig.Qps)
				sells = append(sells, showConfigMap)
			}
		}
		zoneConfig["engine_versions"] = engineVersions
		zoneConfig["sells"] = sells

		zoneConfigs = append(zoneConfigs, zoneConfig)
	}

	if err := d.Set("list", zoneConfigs); err != nil {
		log.Printf("[CRITAL]%s provider set zoneConfigs fail, reason:%s\n ", logId, err.Error())
	}
	d.SetId("zoneconfig" + region)

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), zoneConfigs); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
		}

	}
	return nil
}
