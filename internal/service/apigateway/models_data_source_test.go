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

func TestAccAPIGatewayModelsDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_api_gateway_model.test"
	dataSourceName := "data.aws_api_gateway_models.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.APIGatewayServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckModelDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccModelsDataSourceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckResourceAttrGreaterThanOrEqualValue(dataSourceName, "models.#", 1),
					resource.TestCheckTypeSetElemAttrPair(dataSourceName, "models.*.model_name", resourceName, "model_name"),
					resource.TestCheckTypeSetElemAttrPair(dataSourceName, "models.*.description", resourceName, names.AttrDescription),
					resource.TestCheckTypeSetElemAttrPair(dataSourceName, "models.*.content_type", resourceName, "content_type"),
					resource.TestCheckTypeSetElemAttrPair(dataSourceName, "models.*.schema", resourceName, names.AttrSchema),
					resource.TestCheckNoResourceAttr(dataSourceName, "models.*.value"),
				),
			},
		},
	})
}

func TestAccAPIGatewayModelsDataSource_manyKeys(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	dataSourceName := "data.aws_api_gateway_models.test"
	keyCount := 2

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.APIGatewayServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckModelDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccModelsDataSourceConfig_manyKeys(rName, keyCount),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckResourceAttrGreaterThanOrEqualValue(dataSourceName, "models.#", keyCount),
				),
			},
		},
	})
}

func testAccModelsDataSourceConfig_basic(rName string) string {
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

data "aws_api_gateway_models" "test" {
  rest_api_id = aws_api_gateway_rest_api.test.id
  depends_on  = [aws_api_gateway_model.test]
}
`, rName)
}

func testAccModelsDataSourceConfig_manyKeys(rName string, count int) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_rest_api" "test" {
  name = %[1]q
}

resource "aws_api_gateway_model" "test" {
  count        = %[2]d
  name         = "TestModel${count.index}"
  rest_api_id  = aws_api_gateway_rest_api.test.id
  description  = "Test JSON schema"
  content_type = "application/json"

  schema = jsonencode({
    type = "object"
  })

  depends_on = [aws_api_gateway_rest_api.test]
}

data "aws_api_gateway_models" "test" {
  rest_api_id = aws_api_gateway_rest_api.test.id
  depends_on  = [aws_api_gateway_model.test]
}
`, rName, count)
}
