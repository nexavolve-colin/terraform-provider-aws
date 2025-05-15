---
subcategory: "API Gateway"
layout: "aws"
page_title: "AWS: aws_api_gateway_request_validator"
description: |-
  Terraform data source for retrieving details about a specific AWS REST API Gateway Request Validator.
---

# Data Source: aws_api_gateway_request_validator

Provides details about a specific AWS REST API Gateway Request Validator.

## Example Usage

```terraform
data "aws_api_gateway_request_validator" "example" {
  rest_api_id = aws_api_gateway_rest_api.example.id
  id          = data.aws_api_gateway_request_validators.example.request_validators[0].id
}
```

## Argument Reference

This data source supports the following arguments:

* `rest_api_id` - (Required) ID of the associated AWS REST API Gateway.
* `id` - (Required) ID of the request validator.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `name` - Name of the Request Validator.
* `validate_request_body` - Whether the request body is validated.
* `validate_request_parameters` - Whether the request parameters are validated.
