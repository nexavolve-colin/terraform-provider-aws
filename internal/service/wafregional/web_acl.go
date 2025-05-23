// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package wafregional

import (
	"context"
	"log"
	"slices"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/wafregional"
	awstypes "github.com/aws/aws-sdk-go-v2/service/wafregional/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/enum"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_wafregional_web_acl", name="Web ACL")
// @Tags(identifierAttribute="arn")
func resourceWebACL() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceWebACLCreate,
		ReadWithoutTimeout:   resourceWebACLRead,
		UpdateWithoutTimeout: resourceWebACLUpdate,
		DeleteWithoutTimeout: resourceWebACLDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			names.AttrARN: {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrDefaultAction: {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						names.AttrType: {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: enum.Validate[awstypes.WafActionType](),
						},
					},
				},
			},
			names.AttrLoggingConfiguration: {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log_destination": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: verify.ValidARN,
						},
						"redacted_fields": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_to_match": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"data": {
													Type:     schema.TypeString,
													Optional: true,
												},
												names.AttrType: {
													Type:             schema.TypeString,
													Required:         true,
													ValidateDiagFunc: enum.Validate[awstypes.MatchFieldType](),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			names.AttrMetricName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			names.AttrName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			names.AttrRule: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						names.AttrAction: {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									names.AttrType: {
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: enum.Validate[awstypes.WafActionType](),
									},
								},
							},
						},
						"override_action": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									names.AttrType: {
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: enum.Validate[awstypes.WafOverrideActionType](),
									},
								},
							},
						},
						names.AttrPriority: {
							Type:     schema.TypeInt,
							Required: true,
						},
						names.AttrType: {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          awstypes.WafRuleTypeRegular,
							ValidateDiagFunc: enum.Validate[awstypes.WafRuleType](),
						},
						"rule_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},
	}
}

func resourceWebACLCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).WAFRegionalClient(ctx)
	region := meta.(*conns.AWSClient).Region(ctx)

	name := d.Get(names.AttrName).(string)
	output, err := newRetryer(conn, region).RetryWithToken(ctx, func(token *string) (any, error) {
		input := &wafregional.CreateWebACLInput{
			ChangeToken:   token,
			DefaultAction: expandAction(d.Get(names.AttrDefaultAction).([]any)),
			MetricName:    aws.String(d.Get(names.AttrMetricName).(string)),
			Name:          aws.String(name),
			Tags:          getTagsIn(ctx),
		}

		return conn.CreateWebACL(ctx, input)
	})

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating WAF Regional Web ACL (%s): %s", name, err)
	}

	d.SetId(aws.ToString(output.(*wafregional.CreateWebACLOutput).WebACL.WebACLId))

	if loggingConfiguration := d.Get(names.AttrLoggingConfiguration).([]any); len(loggingConfiguration) == 1 {
		arn := arn.ARN{
			Partition: meta.(*conns.AWSClient).Partition(ctx),
			Service:   "waf-regional",
			Region:    meta.(*conns.AWSClient).Region(ctx),
			AccountID: meta.(*conns.AWSClient).AccountID(ctx),
			Resource:  "webacl/" + d.Id(),
		}.String()

		input := &wafregional.PutLoggingConfigurationInput{
			LoggingConfiguration: expandLoggingConfiguration(loggingConfiguration, arn),
		}

		_, err := conn.PutLoggingConfiguration(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "putting WAF Regional Web ACL (%s) logging configuration: %s", d.Id(), err)
		}
	}

	if rules := d.Get(names.AttrRule).(*schema.Set).List(); len(rules) > 0 {
		_, err := newRetryer(conn, region).RetryWithToken(ctx, func(token *string) (any, error) {
			input := &wafregional.UpdateWebACLInput{
				ChangeToken:   token,
				DefaultAction: expandAction(d.Get(names.AttrDefaultAction).([]any)),
				Updates:       diffWebACLRules([]any{}, rules),
				WebACLId:      aws.String(d.Id()),
			}

			return conn.UpdateWebACL(ctx, input)
		})

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating WAF Regional Web ACL (%s) rules: %s", d.Id(), err)
		}
	}

	return append(diags, resourceWebACLRead(ctx, d, meta)...)
}

func resourceWebACLRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).WAFRegionalClient(ctx)

	webACL, err := findWebACLByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] WAF Regional Web ACL (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading WAF Regional Web ACL (%s): %s", d.Id(), err)
	}

	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition(ctx),
		Service:   "waf-regional",
		Region:    meta.(*conns.AWSClient).Region(ctx),
		AccountID: meta.(*conns.AWSClient).AccountID(ctx),
		Resource:  "webacl/" + d.Id(),
	}.String()
	d.Set(names.AttrARN, arn)
	if err := d.Set(names.AttrDefaultAction, flattenAction(webACL.DefaultAction)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting default_action: %s", err)
	}
	d.Set(names.AttrMetricName, webACL.MetricName)
	d.Set(names.AttrName, webACL.Name)
	if err := d.Set(names.AttrRule, flattenWebACLRules(webACL.Rules)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting rule: %s", err)
	}

	input := &wafregional.GetLoggingConfigurationInput{
		ResourceArn: aws.String(arn),
	}

	output, err := conn.GetLoggingConfiguration(ctx, input)

	loggingConfiguration := []any{}
	switch {
	case err == nil:
		loggingConfiguration = flattenLoggingConfiguration(output.LoggingConfiguration)
	case errs.IsA[*awstypes.WAFNonexistentItemException](err):
	default:
		return sdkdiag.AppendErrorf(diags, "reading WAF Regional Web ACL (%s) logging configuration: %s", d.Id(), err)
	}

	if err := d.Set(names.AttrLoggingConfiguration, loggingConfiguration); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting logging_configuration: %s", err)
	}

	return diags
}

func resourceWebACLUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).WAFRegionalClient(ctx)
	region := meta.(*conns.AWSClient).Region(ctx)

	if d.HasChanges(names.AttrDefaultAction, names.AttrRule) {
		o, n := d.GetChange(names.AttrRule)
		oldR, newR := o.(*schema.Set).List(), n.(*schema.Set).List()

		_, err := newRetryer(conn, region).RetryWithToken(ctx, func(token *string) (any, error) {
			input := &wafregional.UpdateWebACLInput{
				ChangeToken:   token,
				DefaultAction: expandAction(d.Get(names.AttrDefaultAction).([]any)),
				Updates:       diffWebACLRules(oldR, newR),
				WebACLId:      aws.String(d.Id()),
			}

			return conn.UpdateWebACL(ctx, input)
		})

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating WAF Regional Web ACL (%s): %s", d.Id(), err)
		}
	}

	if d.HasChange(names.AttrLoggingConfiguration) {
		if loggingConfiguration := d.Get(names.AttrLoggingConfiguration).([]any); len(loggingConfiguration) == 1 {
			input := &wafregional.PutLoggingConfigurationInput{
				LoggingConfiguration: expandLoggingConfiguration(loggingConfiguration, d.Get(names.AttrARN).(string)),
			}

			_, err := conn.PutLoggingConfiguration(ctx, input)

			if err != nil {
				return sdkdiag.AppendErrorf(diags, "putting WAF Regional Web ACL (%s) logging configuration: %s", d.Id(), err)
			}
		} else {
			input := &wafregional.DeleteLoggingConfigurationInput{
				ResourceArn: aws.String(d.Get(names.AttrARN).(string)),
			}

			_, err := conn.DeleteLoggingConfiguration(ctx, input)

			if err != nil {
				return sdkdiag.AppendErrorf(diags, "deleting WAF Regional Web ACL (%s) logging configuration: %s", d.Id(), err)
			}
		}
	}

	return append(diags, resourceWebACLRead(ctx, d, meta)...)
}

func resourceWebACLDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).WAFRegionalClient(ctx)
	region := meta.(*conns.AWSClient).Region(ctx)

	if rules := d.Get(names.AttrRule).(*schema.Set).List(); len(rules) > 0 {
		_, err := newRetryer(conn, region).RetryWithToken(ctx, func(token *string) (any, error) {
			input := &wafregional.UpdateWebACLInput{
				ChangeToken:   token,
				DefaultAction: expandAction(d.Get(names.AttrDefaultAction).([]any)),
				Updates:       diffWebACLRules(rules, []any{}),
				WebACLId:      aws.String(d.Id()),
			}

			return conn.UpdateWebACL(ctx, input)
		})

		if err != nil && !errs.IsA[*awstypes.WAFNonexistentItemException](err) && !errs.IsA[*awstypes.WAFNonexistentContainerException](err) {
			return sdkdiag.AppendErrorf(diags, "updating WAF Regional Web ACL (%s) rules: %s", d.Id(), err)
		}
	}

	log.Printf("[INFO] Deleting WAF Regional Web ACL: %s", d.Id())
	_, err := newRetryer(conn, region).RetryWithToken(ctx, func(token *string) (any, error) {
		input := &wafregional.DeleteWebACLInput{
			ChangeToken: token,
			WebACLId:    aws.String(d.Id()),
		}

		return conn.DeleteWebACL(ctx, input)
	})

	if errs.IsA[*awstypes.WAFNonexistentItemException](err) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting WAF Regional Web ACL (%s): %s", d.Id(), err)
	}

	return diags
}

func findWebACLByID(ctx context.Context, conn *wafregional.Client, id string) (*awstypes.WebACL, error) {
	input := &wafregional.GetWebACLInput{
		WebACLId: aws.String(id),
	}

	output, err := conn.GetWebACL(ctx, input)

	if errs.IsA[*awstypes.WAFNonexistentItemException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.WebACL == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.WebACL, nil
}

func expandLoggingConfiguration(l []any, resourceARN string) *awstypes.LoggingConfiguration {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]any)

	loggingConfiguration := &awstypes.LoggingConfiguration{
		LogDestinationConfigs: []string{
			m["log_destination"].(string),
		},
		RedactedFields: expandRedactedFields(m["redacted_fields"].([]any)),
		ResourceArn:    aws.String(resourceARN),
	}

	return loggingConfiguration
}

func expandRedactedFields(l []any) []awstypes.FieldToMatch {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]any)

	if m["field_to_match"] == nil {
		return nil
	}

	redactedFields := make([]awstypes.FieldToMatch, 0)

	for _, fieldToMatch := range m["field_to_match"].(*schema.Set).List() {
		if fieldToMatch == nil {
			continue
		}

		redactedFields = append(redactedFields, *expandFieldToMatch(fieldToMatch.(map[string]any)))
	}

	return redactedFields
}

func flattenLoggingConfiguration(loggingConfiguration *awstypes.LoggingConfiguration) []any {
	if loggingConfiguration == nil {
		return []any{}
	}

	m := map[string]any{
		"log_destination": "",
		"redacted_fields": flattenRedactedFields(loggingConfiguration.RedactedFields),
	}

	if len(loggingConfiguration.LogDestinationConfigs) > 0 {
		m["log_destination"] = loggingConfiguration.LogDestinationConfigs[0]
	}

	return []any{m}
}

func flattenRedactedFields(fieldToMatches []awstypes.FieldToMatch) []any {
	if len(fieldToMatches) == 0 {
		return []any{}
	}

	fieldToMatchResource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			names.AttrType: {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
	l := make([]any, len(fieldToMatches))

	for i, fieldToMatch := range fieldToMatches {
		l[i] = flattenFieldToMatch(&fieldToMatch)[0]
	}

	m := map[string]any{
		"field_to_match": schema.NewSet(schema.HashResource(fieldToMatchResource), l),
	}

	return []any{m}
}

func diffWebACLRules(oldR, newR []any) []awstypes.WebACLUpdate {
	updates := make([]awstypes.WebACLUpdate, 0)

	for _, or := range oldR {
		aclRule := or.(map[string]any)

		if idx, contains := sliceContainsMap(newR, aclRule); contains {
			newR = slices.Delete(newR, idx, idx+1)
			continue
		}
		updates = append(updates, expandWebACLUpdate(string(awstypes.ChangeActionDelete), aclRule))
	}

	for _, nr := range newR {
		aclRule := nr.(map[string]any)
		updates = append(updates, expandWebACLUpdate(string(awstypes.ChangeActionInsert), aclRule))
	}
	return updates
}

func expandAction(l []any) *awstypes.WafAction {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]any)

	return &awstypes.WafAction{
		Type: awstypes.WafActionType(m[names.AttrType].(string)),
	}
}

func expandOverrideAction(l []any) *awstypes.WafOverrideAction {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]any)

	return &awstypes.WafOverrideAction{
		Type: awstypes.WafOverrideActionType(m[names.AttrType].(string)),
	}
}

func expandWebACLUpdate(updateAction string, aclRule map[string]any) awstypes.WebACLUpdate {
	var rule *awstypes.ActivatedRule

	switch aclRule[names.AttrType].(string) {
	case string(awstypes.WafRuleTypeGroup):
		rule = &awstypes.ActivatedRule{
			OverrideAction: expandOverrideAction(aclRule["override_action"].([]any)),
			Priority:       aws.Int32(int32(aclRule[names.AttrPriority].(int))),
			RuleId:         aws.String(aclRule["rule_id"].(string)),
			Type:           awstypes.WafRuleType(aclRule[names.AttrType].(string)),
		}
	default:
		rule = &awstypes.ActivatedRule{
			Action:   expandAction(aclRule[names.AttrAction].([]any)),
			Priority: aws.Int32(int32(aclRule[names.AttrPriority].(int))),
			RuleId:   aws.String(aclRule["rule_id"].(string)),
			Type:     awstypes.WafRuleType(aclRule[names.AttrType].(string)),
		}
	}

	update := awstypes.WebACLUpdate{
		Action:        awstypes.ChangeAction(updateAction),
		ActivatedRule: rule,
	}

	return update
}

func flattenAction(n *awstypes.WafAction) []map[string]any {
	if n == nil {
		return nil
	}

	result := map[string]any{
		names.AttrType: string(n.Type),
	}

	return []map[string]any{result}
}

func flattenWebACLRules(ts []awstypes.ActivatedRule) []map[string]any {
	out := make([]map[string]any, len(ts))
	for i, r := range ts {
		m := make(map[string]any)

		switch r.Type {
		case awstypes.WafRuleTypeGroup:
			actionMap := map[string]any{
				names.AttrType: r.OverrideAction.Type,
			}
			m["override_action"] = []map[string]any{actionMap}
		default:
			actionMap := map[string]any{
				names.AttrType: r.Action.Type,
			}
			m[names.AttrAction] = []map[string]any{actionMap}
		}

		m[names.AttrPriority] = r.Priority
		m["rule_id"] = aws.ToString(r.RuleId)
		m[names.AttrType] = string(r.Type)
		out[i] = m
	}
	return out
}
