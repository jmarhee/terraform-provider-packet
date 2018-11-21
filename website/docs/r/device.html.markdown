---
layout: "packet"
page_title: "Packet: packet_device"
sidebar_current: "docs-packet-resource-device"
description: |-
  Provides a Packet device resource. This can be used to create, modify, and delete devices.
---

# packet\_device

Provides a Packet device resource. This can be used to create,
modify, and delete devices.

~> **Note:** All arguments including the root_password and user_data will be stored in
 the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).


## Example Usage

```hcl
# Create a device and add it to cool_project
resource "packet_device" "web1" {
  hostname         = "tf.coreos2"
  plan             = "baremetal_1"
  facility         = "ewr1"
  operating_system = "coreos_stable"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.cool_project.id}"
}
```

```hcl
# Specify more facilities. The device will be deployed to first avaiable.
resource "packet_device" "web1" {
  hostname         = "tf.coreos2"
  plan             = "baremetal_1"
  facility         = "ewr1,sjc1"
  operating_system = "coreos_stable"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.cool_project.id}"
}
```

```hcl
# Same as above, but boot via iPXE initially, using the Ignition Provider for provisioning
resource "packet_device" "pxe1" {
  hostname         = "tf.coreos2-pxe"
  plan             = "baremetal_1"
  facility         = "ewr1"
  operating_system = "custom_ipxe"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.cool_project.id}"
  ipxe_script_url  = "https://rawgit.com/cloudnativelabs/pxe/master/packet/coreos-stable-packet.ipxe"
  always_pxe       = "false"
  user_data        = "${data.ignition_config.example.rendered}"
  }
```

```hcl
# Deploy device on next-available reserved hardware and do custom partitioning.
resource "packet_device" "web1" {
  hostname         = "tftest"
  plan             = "baremetal_0"
  facility         = "sjc1"
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.cool_project.id}"
  hardware_reservation_id = "next-available"
  storage = <<EOS
{
  "disks": [
    {
      "device": "/dev/sda",
      "wipeTable": true,
      "partitions": [
        {
          "label": "BIOS",
          "number": 1,
          "size": 4096
        },
        {
          "label": "SWAP",
          "number": 2,
          "size": "3993600"
        },
        {
          "label": "ROOT",
          "number": 3,
          "size": 0
        }
      ]
    }
  ],
  "filesystems": [
    {
      "mount": {
        "device": "/dev/sda3",
        "format": "ext4",
        "point": "/",
        "create": {
          "options": [
            "-L",
            "ROOT"
          ]
        }
      }
    },
    {
      "mount": {
        "device": "/dev/sda2",
        "format": "swap",
        "point": "none",
        "create": {
          "options": [
            "-L",
            "SWAP"
          ]
        }
      }
    }
  ]
}
  EOS
}
```

## Argument Reference

The following arguments are supported:

* `hostname` - (Required) The device name
* `project_id` - (Required) The id of the project in which to create the device
* `operating_system` - (Required) The operating system slug. To find the slug, or visit [Operating Systems API docs](https://www.packet.net/developers/api/#operatingsystems), set your API auth token in the top of the page and see JSON from the API response.
* `facility` - (Required) The facility in which to create the device. It can be a single facility or a comma-separated listi, in which case the device will be deployed to first listed facility with free capacity. To find the facility codes, visit [Facilities API docs](https://www.packet.net/developers/api/#facilities), set your API auth token in the top of the page and see JSON from the API response. In addition to the listed facility codes, user can pass "any", in order to deploy to facility with most free capacity.
* `plan` - (Required) The device plan slug. To find the plan slug, visit [Device plans API docs](https://www.packet.net/developers/api/#plans), set your auth token in the top of the page and see JSON from the API response.
* `billing_cycle` - (Required) monthly or hourly
* `user_data` (Optional) - A string of the desired User Data for the device.
* `public_ipv4_subnet_size` (Optional) - Size of allocated subnet, more
  information is in the
  [Custom Subnet Size](https://help.packet.net/article/55-custom-subnet-size) doc.
* `ipxe_script_url` (Optional) - URL pointing to a hosted iPXE script. More
  information is in the
  [Custom iPXE](https://help.packet.net/article/26-custom-ipxe)
  doc.
* `always_pxe` (Optional) - If true, a device with OS `custom_ipxe` will
  continue to boot via iPXE on reboots.
* `hardware_reservation_id` (Optional) - The id of hardware reservation where you want this device deployed, or `next-available` if you want to pick your next available reservation automatically.
* `storage` (Optional) - JSON for custom partitioning. Only usable on reserved hardware. More information in in the [Custom Partitioning and RAID](https://help.packet.net/article/61-custom-partitioning-raid) doc.
* `tags` - Tags attached to the device
* `description` - Description string for the device

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the device
* `hostname`- The hostname of the device
* `project_id`- The ID of the project the device belongs to
* `facility` - The facility where the device is deployed.
* `plan` - The hardware config of the device
* `network` - The device's private and public IP (v4 and v6) network details
* `access_public_ipv6` - The ipv6 maintenance IP assigned to the device
* `access_public_ipv4` - The ipv4 maintenance IP assigned to the device
* `access_private_ipv4` - The ipv4 private IP assigned to the device
* `locked` - Whether the device is locked
* `billing_cycle` - The billing cycle of the device (monthly or hourly)
* `operating_system` - The operating system running on the device
* `state` - The status of the device
* `created` - The timestamp for when the device was created
* `updated` - The timestamp for the last time the device was updated
* `tags` - Tags attached to the device
* `description` - Description string for the device
* `hardware_reservation_id` - The id of hardware reservation which this device occupies
* `root_password` - Root password to the server (disabled after 24 hours)
