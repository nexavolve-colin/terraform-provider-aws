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

func TestAccAPIGatewayModelDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_api_gateway_request_validator.test"
	dataSourceName := "data.aws_api_gateway_request_validator.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.APIGatewayServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckModelDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccModelDataSourceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckResourceAttrGreaterThanOrEqualValue(dataSourceName, "models.#", 1),
					resource.TestCheckTypeSetElemAttrPair(dataSourceName, "models.*.name", resourceName, names.AttrName),
					resource.TestCheckNoResourceAttr(dataSourceName, "models.*.value"),
				),
			},
		},
	})
}

func testAccModelDataSourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_rest_api" "test" {
  name = %[1]q
}

resource "aws_api_gateway_model" "test" {
  name         = "TestModel"
  rest_api_id  = aws_api_gateway_rest_api.test.id
  description  = "Test JSON schema"
  content_type = "application/json"

  schema = jsonencode({
    type = "object"
  })

  depends_on = [aws_api_gateway_rest_api.test]
}

data "aws_api_gateway_model" "test" {
  rest_api_id = aws_api_gateway_rest_api.test.id
  model_name  = aws_api_gateway_model.test.model_name
  depends_on  = [aws_api_gateway_model.test]
}
`, rName)
}
