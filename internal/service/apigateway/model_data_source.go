// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apigateway

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/framework"
	"github.com/hashicorp/terraform-provider-aws/internal/framework/flex"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @FrameworkDataSource("aws_api_gateway_model", name="Model")
func newDataSourceModel(context.Context) (datasource.DataSourceWithConfigure, error) {
	return &dataSourceModel{}, nil
}

const (
	DSNameModel = "Model"
)

type dataSourceModel struct {
	framework.DataSourceWithConfigure
}

func (d *dataSourceModel) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"rest_api_id": schema.StringAttribute{
				Required: true,
			},
			names.AttrID: schema.StringAttribute{
				Computed: true,
			},
			"model_name": schema.StringAttribute{
				Required: true,
			},
			"content_type": schema.StringAttribute{
				Computed: true,
			},
			names.AttrSchema: schema.StringAttribute{
				Computed: true,
			},
			names.AttrDescription: schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *dataSourceModel) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data dataSourceModelModel

	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	conn := d.Meta().APIGatewayClient(ctx)
	input := apigateway.GetModelInput{
		RestApiId: flex.StringFromFramework(ctx, data.RestApiId),
		ModelName: flex.StringFromFramework(ctx, data.Name),
	}

	output, err := conn.GetModel(ctx, &input)

	if err != nil {
		response.Diagnostics.AddError(
			create.ProblemStandardMessage(names.APIGateway, create.ErrActionReading, DSNameModel, data.ID.ValueString(), err),
			err.Error(),
		)
		return
	}

	response.Diagnostics.Append(flex.Flatten(ctx, output, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

type dataSourceModelModel struct {
	RestApiId types.String `tfsdk:"rest_api_id"`
	modelModel
}

type modelModel struct {
	ID          types.String `tfsdk:"id"`
	ContentType types.String `tfsdk:"content_type"`
	Description types.String `tfsdk:"description"`
	Name        types.String `tfsdk:"model_name"`
	Schema      types.String `tfsdk:"schema"`
}
