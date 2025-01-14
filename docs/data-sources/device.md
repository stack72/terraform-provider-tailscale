---
page_title: "device Data Source - terraform-provider-tailscale"
subcategory: ""
description: |-
The device data source describes a single device in a tailnet.
---

# Data Source `device`

The device data source describes a single device in a tailnet.

## Example Usage

```terraform
data "tailscale_device" "sample_device" {
  name = "user1-device.example.com"
  wait_for = "60s"
}

```

## Argument Reference

- `name` - (Required) The name of the tailnet device to obtain the attributes of.
- `wait_for` - (Optional) If specified, the provider will retry obtaining the device data every second until the specified duration has been reached. Must be a value greater than 1 second

## Attributes Reference

The following attributes are exported.

- `id` - The unique identifier for the device
- `user` - The user associated with the device
- `addresses` - Tailscale IPs for the device
- `tags` - Tags applied to the device
