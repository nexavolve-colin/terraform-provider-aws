---
subcategory: "API Gateway"
layout: "aws"
page_title: "AWS: aws_api_gateway_models"
description: |-
  Terraform data source for retrieving details about all Models for an AWS REST API Gateway.
---

# Data Source: aws_api_gateway_models

Provides details about all Models for an AWS REST API Gateway.

## Example Usage

```terraform
data "aws_api_gateway_models" "example" {
  rest_api_id = aws_api_gateway_rest_api.example.id
}
```

## Argument Reference

This data source supports the following arguments:

* `rest_api_id` - (Required) ID of the associated AWS REST API Gateway

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `models` - List of objects containing Models information. See below.

### `models`

* `id` - ID of the model
* `model_name` - Name of the model
* `description` - Description of the model
* `content_type` - Content type of the model
* `schema` - Schema of the model in JSON
