---
layout: "yandex"
page_title: "Provider: Yandex Cloud"
sidebar_current: "docs-yandex-index"
description: |-
  The Yandex Cloud provider is used to interact with Yandex Cloud services.
  The provider needs to be configured with the proper credentials before it can be used.
---

# Yandex Cloud Provider

The Yandex Cloud provider is used to interact with
[Yandex Cloud services](https://cloud.yandex.com/). The provider needs
to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
// Configure the Yandex Cloud provider
provider "yandex" {
  token     = "auth_token_here"
  cloud_id  = "cloud_id_here"
  folder_id = "folder_id_here"
  zone      = "ru-central1-a"
}

// Create a new instance
resource "yandex_compute_instance" "default" {
  ...
}
```

## Configuration Reference

The following keys can be used to configure the provider.

* `token` - (Required) Security token used for authentication in Yandex Cloud.

  This can also be set as the `YC_TOKEN` environment variable.

* `cloud_id` - (Required) The ID of the [cloud][yandex-cloud] to apply any resources to.

  This can also be set as the `YC_CLOUD_ID` environment variable.

* `folder_id` - (Required) The ID of the [folder][yandex-folder] to operate under, if not specified by a given resource.

  This can also be set as the `YC_FOLDER` environment variable.

* `zone` - (Optional) The default [accessibility zone][yandex-zone] to operate under, if not specified by a given resource.

  This can also be set as the `YC_ZONE` environment variable.


[yandex-cloud]: https://cloud.yandex.com/docs/resource-manager/concepts/resources-hierarchy#cloud
[yandex-folder]: https://cloud.yandex.com/docs/resource-manager/concepts/resources-hierarchy#folder
[yandex-zone]: https://cloud.yandex.com/docs/overview/concepts/geo-scope
