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

func TestAccAPIGatewayRequestValidatorDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_api_gateway_request_validator.test"
	dataSourceName := "data.aws_api_gateway_request_validator.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.APIGatewayServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRequestValidatorDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccRequestValidatorDataSourceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckResourceAttrGreaterThanOrEqualValue(dataSourceName, "request_validators.#", 1),
					resource.TestCheckTypeSetElemAttrPair(dataSourceName, "request_validators.*.name", resourceName, names.AttrName),
					resource.TestCheckNoResourceAttr(dataSourceName, "request_validators.*.value"),
				),
			},
		},
	})
}

func testAccRequestValidatorDataSourceConfig_basic(rName string) string {
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

data "aws_api_gateway_request_validator" "test" {
  rest_api_id = aws_api_gateway_rest_api.test.id
  id          = aws_api_gateway_request_validator.test.id
  depends_on  = [aws_api_gateway_request_validator.test]
}
`, rName)
}
