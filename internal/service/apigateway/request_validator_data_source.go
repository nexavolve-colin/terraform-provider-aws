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

// @FrameworkDataSource("aws_api_gateway_request_validator", name="Request Validator")
func newDataSourceRequestValidator(context.Context) (datasource.DataSourceWithConfigure, error) {
	return &dataSourceRequestValidator{}, nil
}

const (
	DSNameRequestValidator = "Request Validator"
)

type dataSourceRequestValidator struct {
	framework.DataSourceWithConfigure
}

func (d *dataSourceRequestValidator) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"rest_api_id": schema.StringAttribute{
				Required: true,
			},
			names.AttrID: schema.StringAttribute{
				Required: true,
			},
			names.AttrName: schema.StringAttribute{
				Computed: true,
			},
			"validate_request_body": schema.BoolAttribute{
				Computed: true,
			},
			"validate_request_parameters": schema.BoolAttribute{
				Computed: true,
			},
		},
	}
}

func (d *dataSourceRequestValidator) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data dataSourceRequestValidatorModel

	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	conn := d.Meta().APIGatewayClient(ctx)
	input := apigateway.GetRequestValidatorInput{
		RestApiId:          flex.StringFromFramework(ctx, data.RestApiId),
		RequestValidatorId: flex.StringFromFramework(ctx, data.ID),
	}

	output, err := conn.GetRequestValidator(ctx, &input)

	if err != nil {
		response.Diagnostics.AddError(
			create.ProblemStandardMessage(names.APIGateway, create.ErrActionReading, DSNameRequestValidator, data.ID.ValueString(), err),
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

type dataSourceRequestValidatorModel struct {
	RestApiId types.String `tfsdk:"rest_api_id"`
	requestValidatorModel
}

type requestValidatorModel struct {
	ID                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	ValidateRequestBody       types.Bool   `tfsdk:"validate_request_body"`
	ValidateRequestParameters types.Bool   `tfsdk:"validate_request_parameters"`
}
