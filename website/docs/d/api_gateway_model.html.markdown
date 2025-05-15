---
subcategory: "API Gateway"
layout: "aws"
page_title: "AWS: aws_api_gateway_model"
description: |-
  Terraform data source for retrieving details about a specific AWS REST API Gateway Model.
---

# Data Source: aws_api_gateway_model

Provides details about a specific AWS REST API Gateway Model.

## Example Usage

```terraform
data "aws_api_gateway_model" "example" {
  rest_api_id = aws_api_gateway_rest_api.example.id
  id          = data.aws_api_gateway_models.example.models[0].id
}
```

## Argument Reference

This data source supports the following arguments:

* `rest_api_id` - (Required) ID of the associated AWS REST API Gateway
* `model_name` - (Required) Name of the model

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `id` - ID of the model
* `content_type` - Content type of the model
* `schema` - Schema of the model in JSON
* `description` - Description of the model
