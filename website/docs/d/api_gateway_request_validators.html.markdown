---
subcategory: "API Gateway"
layout: "aws"
page_title: "AWS: aws_api_gateway_request_validators"
description: |-
  Terraform data source for retrieving details about Request Validators for an AWS REST API Gateway.
---

# Data Source: aws_api_gateway_request_validators

Retrieves details about Request Validators for an AWS REST API Gateway.

## Example Usage

```terraform
data "aws_api_gateway_request_validators" "example" {
  rest_api_id = aws_api_gateway_rest_api.example.id
}
```

## Argument Reference

This data source supports the following arguments:

* `rest_api_id` - (Required) ID of the associated AWS REST API Gateway.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `request_validators` - List of objects containing Request Validators information. See below.

### `request_validators`

* `id` - ID of the Request Validator.
* `name` - Name of the Request Validator.
* `validate_request_body` - Whether the request body is validated.
* `validate_request_parameters` - Whether the request parameters are validated.
