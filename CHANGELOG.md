## 1.4.0 (Unreleased)

ENHANCEMENTS:

* resource/tencentcloud_instance: data disk count limit now is upgrade from 1 to 10, as API has supported more disks.

## 1.3.0 (March 12, 2019)

FEATURES:

* **New Resource**: `tencentcloud_lb` ([#3](https://github.com/terraform-providers/terraform-provider-scaffolding/issues/3))

ENHANCEMENTS:

* resource/tencentcloud_instance: Add `user_data_raw` argument ([#4](https://github.com/terraform-providers/terraform-provider-scaffolding/issues/4))


## 1.2.2 (September 28, 2018)

BUG FIXES:

* resource/tencentcloud_cbs_storage: make name to be required ([#25](https://github.com/tencentyun/terraform-provider-tencentcloud/issues/25))
* resource/tencentcloud_instance: support user data and private ip

## 1.2.0 (April 3, 2018)

FEATURES:

* **New Resource**: `tencentcloud_container_cluster`
* **New Resource**: `tencentcloud_container_cluster_instance`
* **New Data Source**: `tencentcloud_container_clusters`
* **New Data Source**: `tencentcloud_container_cluster_instances`

## 1.1.0 (March 9, 2018)

FEATURES:

* **New Resource**: `tencentcloud_eip`
* **New Resource**: `tencentcloud_eip_association`
* **New Data Source**: `tencentcloud_eip`
* **New Resource**: `tencentcloud_nat_gateway`
* **New Resource**: `tencentcloud_dnat`
* **New Data Source**: `tencentcloud_nats`
* **New Resource**: `tencentcloud_cbs_snapshot`
* **New Resource**: `tencentcloud_alb_server_attachment`

## 1.0.0 (January 19, 2018)

FEATURES:

### CVM

RESOURCES:

* instance create
* instance read
* instance update
    * reset instance
    * reset password
    * update instance name
    * update security groups
* instance delete
* key pair create
* key pair read
* key pair delete

DATA SOURCES:

* image read
* instance\_type read
* zone read

### VPC

RESOURCES:

* vpc create
* vpc read
* vpc update (update name)
* vpc delete
* subnet create
* subnet read
* subnet update (update name)
* subnet delete
* security group create
* security group read
* security group update (update name, description)
* security group delete
* security group rule create
* security group rule read
* security group rule delete
* route table create
* route table read
* route table update (update name)
* route table delete
* route entry create
* route entry read
* route entry delete

DATA SOURCES:

* vpc read
* subnet read
* security group read
* route table read

### CBS

RESOURCES:

* storage create
* storage read
* storage update (update name)
* storage attach
* storage detach
