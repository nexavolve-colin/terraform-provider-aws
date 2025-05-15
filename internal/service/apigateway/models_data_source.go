// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apigateway

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	awstypes "github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/framework"
	"github.com/hashicorp/terraform-provider-aws/internal/framework/flex"
	fwtypes "github.com/hashicorp/terraform-provider-aws/internal/framework/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @FrameworkDataSource("aws_api_gateway_models", name="Models")
func newDataSourceModels(context.Context) (datasource.DataSourceWithConfigure, error) {
	return &dataSourceModels{}, nil
}

const (
	DSNameModels = "Models"
)

type dataSourceModels struct {
	framework.DataSourceWithConfigure
}

func (d *dataSourceModels) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"rest_api_id": schema.StringAttribute{
				Required: true,
			},
			"models": framework.DataSourceComputedListOfObjectAttribute[modelModel](ctx),
		},
	}
}

func (d *dataSourceModels) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data dataSourceModelsModel

	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	conn := d.Meta().APIGatewayClient(ctx)
	input := apigateway.GetModelsInput{
		RestApiId: flex.StringFromFramework(ctx, data.RestApiId),
	}

	items, err := findModels(ctx, conn, &input)
	if err != nil {
		response.Diagnostics.AddError(
			create.ProblemStandardMessage(names.APIGateway, create.ErrActionReading, DSNameModels, data.RestApiId.ValueString(), err),
			err.Error(),
		)
		return
	}

	response.Diagnostics.Append(flex.Flatten(ctx, items, &data.Models)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func findModels(ctx context.Context, conn *apigateway.Client, input *apigateway.GetModelsInput) ([]awstypes.Model, error) {
	var items []awstypes.Model

	err := getModelsPages(ctx, conn, input, func(page *apigateway.GetModelsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.Items {
			items = append(items, v)
		}

		return !lastPage
	})

	if err != nil {
		return nil, err
	}

	return items, nil
}

type dataSourceModelsModel struct {
	RestApiId types.String                                `tfsdk:"rest_api_id"`
	Models    fwtypes.ListNestedObjectValueOf[modelModel] `tfsdk:"models"`
}
