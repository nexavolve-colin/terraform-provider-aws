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

// @FrameworkDataSource("aws_api_gateway_request_validators", name="Request Validators")
func newDataSourceRequestValidators(context.Context) (datasource.DataSourceWithConfigure, error) {
	return &dataSourceRequestValidators{}, nil
}

const (
	DSNameRequestValidators = "Request Validators"
)

type dataSourceRequestValidators struct {
	framework.DataSourceWithConfigure
}

func (d *dataSourceRequestValidators) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"rest_api_id": schema.StringAttribute{
				Required: true,
			},
			"request_validators": framework.DataSourceComputedListOfObjectAttribute[requestValidatorModel](ctx),
		},
	}
}

func (d *dataSourceRequestValidators) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data dataSourceRequestValidatorsModel

	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	conn := d.Meta().APIGatewayClient(ctx)
	input := apigateway.GetRequestValidatorsInput{
		RestApiId: flex.StringFromFramework(ctx, data.RestApiId),
	}

	items, err := findRequestValidators(ctx, conn, &input)
	if err != nil {
		response.Diagnostics.AddError(
			create.ProblemStandardMessage(names.APIGateway, create.ErrActionReading, DSNameRequestValidators, data.RestApiId.ValueString(), err),
			err.Error(),
		)
		return
	}

	response.Diagnostics.Append(flex.Flatten(ctx, items, &data.RequestValidators)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func findRequestValidators(ctx context.Context, conn *apigateway.Client, input *apigateway.GetRequestValidatorsInput) ([]awstypes.RequestValidator, error) {
	var items []awstypes.RequestValidator

	err := getRequestValidatorsPages(ctx, conn, input, func(page *apigateway.GetRequestValidatorsOutput, lastPage bool) bool {
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

type dataSourceRequestValidatorsModel struct {
	RestApiId         types.String                                           `tfsdk:"rest_api_id"`
	RequestValidators fwtypes.ListNestedObjectValueOf[requestValidatorModel] `tfsdk:"request_validators"`
}
