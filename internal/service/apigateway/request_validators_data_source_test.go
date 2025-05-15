// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apigateway_test

import (
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccAPIGatewayRequestValidatorsDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_api_gateway_request_validator.test"
	dataSourceName := "data.aws_api_gateway_request_validators.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.APIGatewayServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRequestValidatorDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccRequestValidatorsDataSourceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckResourceAttrGreaterThanOrEqualValue(dataSourceName, "request_validators.#", 1),
					resource.TestCheckTypeSetElemAttrPair(dataSourceName, "request_validators.*.name", resourceName, names.AttrName),
					resource.TestCheckTypeSetElemAttrPair(dataSourceName, "request_validators.*.validate_request_body", resourceName, "validate_request_body"),
					resource.TestCheckTypeSetElemAttrPair(dataSourceName, "request_validators.*.validate_request_parameters", resourceName, "validate_request_parameters"),
					resource.TestCheckNoResourceAttr(dataSourceName, "request_validators.*.value"),
				),
			},
		},
	})
}

func TestAccAPIGatewayRequestValidatorsDataSource_manyKeys(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	dataSourceName := "data.aws_api_gateway_request_validators.test"
	keyCount := 3

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.APIGatewayServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRequestValidatorDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccRequestValidatorsDataSourceConfig_manyKeys(rName, keyCount),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckResourceAttrGreaterThanOrEqualValue(dataSourceName, "request_validators.#", keyCount),
				),
			},
		},
	})
}

func testAccRequestValidatorsDataSourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_rest_api" "test" {
  name = %[1]q
}

resource "aws_api_gateway_request_validator" "test" {
  name                        = "test"
  rest_api_id                 = aws_api_gateway_rest_api.test.id
  validate_request_body       = true
  validate_request_parameters = true
  depends_on                  = [aws_api_gateway_rest_api.test]
}

data "aws_api_gateway_request_validators" "test" {
  rest_api_id = aws_api_gateway_rest_api.test.id
  depends_on  = [aws_api_gateway_request_validator.test]
}
`, rName)
}

func testAccRequestValidatorsDataSourceConfig_manyKeys(rName string, count int) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_rest_api" "test" {
  name = %[1]q
}

resource "aws_api_gateway_request_validator" "test" {
  count                       = %[2]d
  name                        = "test-${count.index}"
  rest_api_id                 = aws_api_gateway_rest_api.test.id
  validate_request_body       = true
  validate_request_parameters = true
  depends_on                  = [aws_api_gateway_rest_api.test]
}

data "aws_api_gateway_request_validators" "test" {
  rest_api_id = aws_api_gateway_rest_api.test.id
  depends_on  = [aws_api_gateway_request_validator.test]
}
`, rName, count)
}
