package services

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/magodo/terraform-provider-outlook/outlook/clients"
	"github.com/magodo/terraform-provider-outlook/outlook/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

func ResourceMessageRule() *schema.Resource {
	predicateSchema := &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		MinItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"body_contains": {
					Type:     schema.TypeSet,
					MinItems: 1,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"body_or_subject_contains": {
					Type:     schema.TypeSet,
					MinItems: 1,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"categories": {
					Type:     schema.TypeSet,
					MinItems: 1,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"from_addresses": {
					Type:     schema.TypeSet,
					MinItems: 1,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"has_attachments": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"header_contains": {
					Type:     schema.TypeSet,
					MinItems: 1,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"importance": {
					Type:     schema.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice(
						[]string{
							string(msgraph.ImportanceVLow),
							string(msgraph.ImportanceVNormal),
							string(msgraph.ImportanceVHigh),
						},
						false,
					),
				},
				"is_approval_request": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"is_automatic_forward": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"is_automatic_reply": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"is_encrypted": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"is_meeting_request": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"is_meeting_response": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"is_non_delivery_report": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"is_permission_controlled": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"is_read_receipt": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"is_signed": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"is_voicemail": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"message_action_flag": {
					Type:     schema.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice(
						[]string{
							string(msgraph.MessageActionFlagVAny),
							string(msgraph.MessageActionFlagVCall),
							string(msgraph.MessageActionFlagVDoNotForward),
							string(msgraph.MessageActionFlagVFollowUp),
							string(msgraph.MessageActionFlagVFyi),
							string(msgraph.MessageActionFlagVForward),
							string(msgraph.MessageActionFlagVNoResponseNecessary),
							string(msgraph.MessageActionFlagVRead),
							string(msgraph.MessageActionFlagVReply),
							string(msgraph.MessageActionFlagVReplyToAll),
							string(msgraph.MessageActionFlagVReview),
						},
						false,
					),
				},
				"not_sent_to_me": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"recipient_contains": {
					Type:     schema.TypeSet,
					MinItems: 1,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"sender_contains": {
					Type:     schema.TypeSet,
					MinItems: 1,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"sensitivity": {
					Type:     schema.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice(
						[]string{
							string(msgraph.SensitivityVNormal),
							string(msgraph.SensitivityVPrivate),
							string(msgraph.SensitivityVPersonal),
							string(msgraph.SensitivityVConfidential),
						},
						false,
					),
				},
				"sent_cc_me": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"sent_only_to_me": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"sent_to_addresses": {
					Type:     schema.TypeSet,
					MinItems: 1,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"sent_to_me": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"sent_to_or_cc_me": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"subject_contains": {
					Type:     schema.TypeSet,
					MinItems: 1,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"within_size_range": {
					Type:     schema.TypeList,
					MaxItems: 1,
					MinItems: 1,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"max_size": {
								Type:     schema.TypeInt,
								Required: true,
							},
							"min_size": {
								Type:     schema.TypeInt,
								Required: true,
							},
						},
					},
				},
			},
		},
	}

	actionList := []string{"action.0.assign_categories", "action.0.copy_to_folder", "action.0.delete", "action.0.forward_as_attachment_to", "action.0.forward_to",
		"action.0.mark_as_read", "action.0.mark_importance", "action.0.move_to_folder", "action.0.permanent_delete", "action.0.redirect_to"}

	return &schema.Resource{
		CreateContext: resourceMessageRuleCreate,
		ReadContext:   resourceMessageRuleRead,
		UpdateContext: resourceMessageRuleUpdate,
		DeleteContext: resourceMessageRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sequence": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"condition": predicateSchema,
			"exception": predicateSchema,
			"action": {
				Type:     schema.TypeList,
				MaxItems: 1,
				MinItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"assign_categories": {
							Type:         schema.TypeSet,
							MinItems:     1,
							Optional:     true,
							Elem:         &schema.Schema{Type: schema.TypeString},
							AtLeastOneOf: actionList,
						},
						"copy_to_folder": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: actionList,
						},
						"delete": {
							Type:         schema.TypeBool,
							Optional:     true,
							AtLeastOneOf: actionList,
						},
						"forward_as_attachment_to": {
							Type:         schema.TypeSet,
							MinItems:     1,
							Optional:     true,
							Elem:         &schema.Schema{Type: schema.TypeString},
							AtLeastOneOf: actionList,
						},
						"forward_to": {
							Type:         schema.TypeSet,
							MinItems:     1,
							Optional:     true,
							Elem:         &schema.Schema{Type: schema.TypeString},
							AtLeastOneOf: actionList,
						},
						"mark_as_read": {
							Type:         schema.TypeBool,
							Optional:     true,
							AtLeastOneOf: actionList,
						},
						"mark_importance": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice(
								[]string{
									string(msgraph.ImportanceVLow),
									string(msgraph.ImportanceVNormal),
									string(msgraph.ImportanceVHigh),
								},
								false,
							),
							AtLeastOneOf: actionList,
						},
						"move_to_folder": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: actionList,
						},
						"permanent_delete": {
							Type:         schema.TypeBool,
							Optional:     true,
							AtLeastOneOf: actionList,
						},
						"redirect_to": {
							Type:         schema.TypeSet,
							MinItems:     1,
							Optional:     true,
							Elem:         &schema.Schema{Type: schema.TypeString},
							AtLeastOneOf: actionList,
						},
						"stop_processing_rules": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceMessageRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).MessageRules
	name := d.Get("name").(string)

	if d.IsNewResource() {
		req := client.Request()
		req.Filter(fmt.Sprintf(`displayName eq '%s'`, name))
		objs, err := req.Get(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(objs) != 0 {
			return utils.ImportAsExistsError("outlook_message_rule", *(objs[0].ID))
		}
	}

	sequence := d.Get("sequence").(int)
	if sequence == 0 {
		sequence = math.MaxInt16
	}
	param := &msgraph.MessageRule{
		DisplayName: utils.String(name),
		Sequence:    utils.Int(sequence),
		IsEnabled:   utils.Bool(d.Get("enabled").(bool)),
		Conditions:  expandMessageRulePredicate(d.Get("condition").([]interface{})),
		Exceptions:  expandMessageRulePredicate(d.Get("exception").([]interface{})),
		Actions:     expandMessageRuleAction(d.Get("action").([]interface{})),
	}

	resp, err := client.Request().Add(ctx, param)
	if err != nil {
		return diag.Errorf("creating Message Rule %q: %+v", name, err)
	}

	if resp.ID == nil {
		return diag.Errorf("nil ID for Message Rule %+v", name)
	}
	d.SetId(*resp.ID)

	return resourceMessageRuleRead(ctx, d, meta)
}

func resourceMessageRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).MessageRules

	resp, err := client.ID(d.Id()).Request().Get(ctx)
	if err != nil {
		if utils.MessageResponseErrorWasNotFound(err) {
			log.Printf("[WARN] Message Rule %q doesn't exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("name", resp.DisplayName)
	d.Set("sequence", resp.Sequence)
	d.Set("enabled", resp.IsEnabled)
	if err := d.Set("condition", flattenMessageRulePredicate(resp.Conditions)); err != nil {
		return diag.Errorf(`setting "condition": %+v"`, err)
	}
	if err := d.Set("exception", flattenMessageRulePredicate(resp.Exceptions)); err != nil {
		return diag.Errorf(`setting "exception": %+v"`, err)
	}
	if err := d.Set("action", flattenMessageRuleAction(resp.Actions)); err != nil {
		return diag.Errorf(`setting "action": %+v"`, err)
	}

	return nil
}

func resourceMessageRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).MessageRules

	var param msgraph.MessageRule

	if d.HasChange("sequence") {
		param.Sequence = utils.Int(d.Get("sequence").(int))
	}
	if d.HasChange("enabled") {
		param.IsEnabled = utils.Bool(d.Get("enabled").(bool))
	}
	if d.HasChange("condition") {
		param.Conditions = expandMessageRulePredicate(d.Get("condition").([]interface{}))
		// NOTE: When optional attribute is changed to be absent, we should force zero it in request body.
		//       Otherwise, it will be omit in request body.
		if param.Conditions == nil {
			param.Conditions = &msgraph.MessageRulePredicates{}
		}
	}
	if d.HasChange("exception") {
		// NOTE: When optional attribute is changed to be absent, we should force zero it in request body.
		//       Otherwise, it will be omit in request body.
		param.Exceptions = expandMessageRulePredicate(d.Get("exception").([]interface{}))
		if param.Exceptions == nil {
			param.Exceptions = &msgraph.MessageRulePredicates{}
		}
	}
	if d.HasChange("action") {
		param.Actions = expandMessageRuleAction(d.Get("action").([]interface{}))
	}

	if err := client.ID(d.Id()).Request().Update(ctx, &param); err != nil {
		return diag.FromErr(err)
	}

	return resourceMessageRuleRead(ctx, d, meta)
}

func resourceMessageRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).MessageRules
	if err := client.ID(d.Id()).Request().Delete(ctx); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func expandMessageRulePredicate(input []interface{}) *msgraph.MessageRulePredicates {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	output := &msgraph.MessageRulePredicates{
		BodyContains:          *utils.ExpandSlice(raw["body_contains"].(*schema.Set).List(), "", nil).(*[]string),
		BodyOrSubjectContains: *utils.ExpandSlice(raw["body_or_subject_contains"].(*schema.Set).List(), "", nil).(*[]string),
		Categories:            *utils.ExpandSlice(raw["categories"].(*schema.Set).List(), "", nil).(*[]string),
		FromAddresses: *utils.ExpandSlice(raw["from_addresses"].(*schema.Set).List(), msgraph.Recipient{}, func(i interface{}) interface{} {
			return msgraph.Recipient{
				EmailAddress: &msgraph.EmailAddress{
					Address: utils.String(i.(string)),
				},
			}
		}).(*[]msgraph.Recipient),
		HasAttachments:         utils.ToPtrOrNil(raw["has_attachments"].(bool)).(*bool),
		HeaderContains:         *utils.ExpandSlice(raw["header_contains"].(*schema.Set).List(), "", nil).(*[]string),
		Importance:             utils.ToPtrOrNil(msgraph.Importance(raw["importance"].(string))).(*msgraph.Importance),
		IsApprovalRequest:      utils.ToPtrOrNil(raw["is_approval_request"].(bool)).(*bool),
		IsAutomaticForward:     utils.ToPtrOrNil(raw["is_automatic_forward"].(bool)).(*bool),
		IsAutomaticReply:       utils.ToPtrOrNil(raw["is_automatic_reply"].(bool)).(*bool),
		IsEncrypted:            utils.ToPtrOrNil(raw["is_encrypted"].(bool)).(*bool),
		IsMeetingRequest:       utils.ToPtrOrNil(raw["is_meeting_request"].(bool)).(*bool),
		IsMeetingResponse:      utils.ToPtrOrNil(raw["is_meeting_response"].(bool)).(*bool),
		IsNonDeliveryReport:    utils.ToPtrOrNil(raw["is_non_delivery_report"].(bool)).(*bool),
		IsPermissionControlled: utils.ToPtrOrNil(raw["is_permission_controlled"].(bool)).(*bool),
		IsReadReceipt:          utils.ToPtrOrNil(raw["is_read_receipt"].(bool)).(*bool),
		IsSigned:               utils.ToPtrOrNil(raw["is_signed"].(bool)).(*bool),
		IsVoicemail:            utils.ToPtrOrNil(raw["is_voicemail"].(bool)).(*bool),
		MessageActionFlag:      utils.ToPtrOrNil(msgraph.MessageActionFlag(raw["message_action_flag"].(string))).(*msgraph.MessageActionFlag),
		NotSentToMe:            utils.ToPtrOrNil(raw["not_sent_to_me"].(bool)).(*bool),
		RecipientContains:      *utils.ExpandSlice(raw["recipient_contains"].(*schema.Set).List(), "", nil).(*[]string),
		SenderContains:         *utils.ExpandSlice(raw["sender_contains"].(*schema.Set).List(), "", nil).(*[]string),
		Sensitivity:            utils.ToPtrOrNil(msgraph.Sensitivity(raw["sensitivity"].(string))).(*msgraph.Sensitivity),
		SentCcMe:               utils.ToPtrOrNil(raw["sent_cc_me"].(bool)).(*bool),
		SentOnlyToMe:           utils.ToPtrOrNil(raw["sent_only_to_me"].(bool)).(*bool),
		SentToAddresses: *utils.ExpandSlice(raw["sent_to_addresses"].(*schema.Set).List(), msgraph.Recipient{}, func(i interface{}) interface{} {
			return msgraph.Recipient{
				EmailAddress: &msgraph.EmailAddress{
					Address: utils.String(i.(string)),
				},
			}
		}).(*[]msgraph.Recipient),
		SentToMe:        utils.ToPtrOrNil(raw["sent_to_me"].(bool)).(*bool),
		SentToOrCcMe:    utils.ToPtrOrNil(raw["sent_to_or_cc_me"].(bool)).(*bool),
		SubjectContains: *utils.ExpandSlice(raw["subject_contains"].(*schema.Set).List(), "", nil).(*[]string),
		WithinSizeRange: expandMessageSizeRange(raw["within_size_range"].([]interface{})),
	}

	return output
}

func expandMessageRuleAction(input []interface{}) *msgraph.MessageRuleActions {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	output := &msgraph.MessageRuleActions{
		AssignCategories: *utils.ExpandSlice(raw["assign_categories"].(*schema.Set).List(), "", nil).(*[]string),
		CopyToFolder:     utils.ToPtrOrNil(raw["copy_to_folder"].(string)).(*string),
		Delete:           utils.ToPtrOrNil(raw["delete"].(bool)).(*bool),
		ForwardAsAttachmentTo: *utils.ExpandSlice(raw["forward_as_attachment_to"].(*schema.Set).List(), msgraph.Recipient{}, func(i interface{}) interface{} {
			return msgraph.Recipient{
				EmailAddress: &msgraph.EmailAddress{
					Address: utils.String(i.(string)),
				},
			}
		}).(*[]msgraph.Recipient),
		ForwardTo: *utils.ExpandSlice(raw["forward_to"].(*schema.Set).List(), msgraph.Recipient{}, func(i interface{}) interface{} {
			return msgraph.Recipient{
				EmailAddress: &msgraph.EmailAddress{
					Address: utils.String(i.(string)),
				},
			}
		}).(*[]msgraph.Recipient),
		MarkAsRead:      utils.ToPtrOrNil(raw["mark_as_read"].(bool)).(*bool),
		MarkImportance:  utils.ToPtrOrNil(msgraph.Importance(raw["mark_importance"].(string))).(*msgraph.Importance),
		MoveToFolder:    utils.ToPtrOrNil(raw["move_to_folder"].(string)).(*string),
		PermanentDelete: utils.ToPtrOrNil(raw["permanent_delete"].(bool)).(*bool),
		RedirectTo: *utils.ExpandSlice(raw["forward_to"].(*schema.Set).List(), msgraph.Recipient{}, func(i interface{}) interface{} {
			return msgraph.Recipient{
				EmailAddress: &msgraph.EmailAddress{
					Address: utils.String(i.(string)),
				},
			}
		}).(*[]msgraph.Recipient),
		StopProcessingRules: utils.ToPtrOrNil(raw["stop_processing_rules"].(bool)).(*bool),
	}

	return output
}

func expandMessageSizeRange(input []interface{}) *msgraph.SizeRange {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	output := &msgraph.SizeRange{
		MaximumSize: utils.Int(raw["max_size"].(int)),
		MinimumSize: utils.Int(raw["min_size"].(int)),
	}
	return output
}

func flattenMessageRulePredicate(input *msgraph.MessageRulePredicates) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"body_contains":            utils.FlattenSlicePtr(utils.ToPtr(input.BodyContains).(*[]string), nil),
			"body_or_subject_contains": utils.FlattenSlicePtr(utils.ToPtr(input.BodyOrSubjectContains).(*[]string), nil),
			"categories":               utils.FlattenSlicePtr(utils.ToPtr(input.Categories).(*[]string), nil),
			"from_addresses": utils.FlattenSlicePtr(utils.ToPtr(input.FromAddresses).(*[]msgraph.Recipient), func(i interface{}) interface{} {
				addr := i.(msgraph.Recipient).EmailAddress
				if addr == nil {
					return ""
				}
				return utils.SafeDeref(addr.Address)
			}),
			"has_attachments":          utils.SafeDeref(input.HasAttachments),
			"header_contains":          utils.FlattenSlicePtr(utils.ToPtr(input.HeaderContains).(*[]string), nil),
			"importance":               string(utils.SafeDeref(input.Importance).(msgraph.Importance)),
			"is_approval_request":      utils.SafeDeref(input.IsApprovalRequest),
			"is_automatic_forward":     utils.SafeDeref(input.IsAutomaticForward),
			"is_automatic_reply":       utils.SafeDeref(input.IsAutomaticReply),
			"is_encrypted":             utils.SafeDeref(input.IsEncrypted),
			"is_meeting_request":       utils.SafeDeref(input.IsMeetingRequest),
			"is_meeting_response":      utils.SafeDeref(input.IsMeetingResponse),
			"is_non_delivery_report":   utils.SafeDeref(input.IsNonDeliveryReport),
			"is_permission_controlled": utils.SafeDeref(input.IsPermissionControlled),
			"is_read_receipt":          utils.SafeDeref(input.IsReadReceipt),
			"is_signed":                utils.SafeDeref(input.IsSigned),
			"is_voicemail":             utils.SafeDeref(input.IsVoicemail),
			"message_action_flag":      string(utils.SafeDeref(input.MessageActionFlag).(msgraph.MessageActionFlag)),
			"not_sent_to_me":           utils.SafeDeref(input.NotSentToMe),
			"recipient_contains":       utils.FlattenSlicePtr(utils.ToPtr(input.RecipientContains).(*[]string), nil),
			"sender_contains":          utils.FlattenSlicePtr(utils.ToPtr(input.SenderContains).(*[]string), nil),
			"sensitivity":              string(utils.SafeDeref(input.Sensitivity).(msgraph.Sensitivity)),
			"sent_cc_me":               utils.SafeDeref(input.SentCcMe),
			"sent_only_to_me":          utils.SafeDeref(input.SentOnlyToMe),
			"sent_to_addresses": utils.FlattenSlicePtr(utils.ToPtr(input.SentToAddresses).(*[]msgraph.Recipient), func(i interface{}) interface{} {
				addr := i.(msgraph.Recipient).EmailAddress
				if addr == nil {
					return ""
				}
				return utils.SafeDeref(addr.Address)
			}),
			"sent_to_me":        utils.SafeDeref(input.SentToMe),
			"sent_to_or_cc_me":  utils.SafeDeref(input.SentToOrCcMe),
			"subject_contains":  utils.FlattenSlicePtr(utils.ToPtr(input.SubjectContains).(*[]string), nil),
			"within_size_range": flattenMessageSizeRange(input.WithinSizeRange),
		},
	}
}

func flattenMessageRuleAction(input *msgraph.MessageRuleActions) interface{} {
	if input == nil {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			"assign_categories": utils.FlattenSlicePtr(utils.ToPtr(input.AssignCategories).(*[]string), nil),
			"copy_to_folder":    utils.SafeDeref(input.CopyToFolder),
			"delete":            utils.SafeDeref(input.Delete),
			"forward_as_attachment_to": utils.FlattenSlicePtr(utils.ToPtr(input.ForwardAsAttachmentTo).(*[]msgraph.Recipient), func(i interface{}) interface{} {
				addr := i.(msgraph.Recipient).EmailAddress
				if addr == nil {
					return ""
				}
				return utils.SafeDeref(addr.Address)
			}),
			"forward_to": utils.FlattenSlicePtr(utils.ToPtr(input.ForwardTo).(*[]msgraph.Recipient), func(i interface{}) interface{} {
				addr := i.(msgraph.Recipient).EmailAddress
				if addr == nil {
					return ""
				}
				return utils.SafeDeref(addr.Address)
			}),
			"mark_as_read":     utils.SafeDeref(input.MarkAsRead),
			"mark_importance":  string(utils.SafeDeref(input.MarkImportance).(msgraph.Importance)),
			"move_to_folder":   utils.SafeDeref(input.MoveToFolder),
			"permanent_delete": utils.SafeDeref(input.PermanentDelete),
			"redirect_to": utils.FlattenSlicePtr(utils.ToPtr(input.RedirectTo).(*[]msgraph.Recipient), func(i interface{}) interface{} {
				addr := i.(msgraph.Recipient).EmailAddress
				if addr == nil {
					return ""
				}
				return utils.SafeDeref(addr.Address)
			}),
			"stop_processing_rules": utils.SafeDeref(input.StopProcessingRules),
		},
	}
}

func flattenMessageSizeRange(input *msgraph.SizeRange) interface{} {
	if input == nil {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			"max_size": utils.SafeDeref(input.MaximumSize),
			"min_size": utils.SafeDeref(input.MinimumSize),
		},
	}
}
