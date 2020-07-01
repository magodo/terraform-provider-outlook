package services

import (
	"context"
	"fmt"
	"log"
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
		Required: true,
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
				//TODO: support this after implement [categories](https://docs.microsoft.com/en-us/graph/api/resources/outlookcategory?view=graph-rest-1.0)
				//                 "categories": {
				//                     Type:     schema.TypeSet,
				//                     MinItems: 1,
				//                     Optional: true,
				//                     Elem:     &schema.Schema{Type: schema.TypeString},
				//                 },
				"from_addresses": {
					Type:     schema.TypeSet,
					MinItems: 1,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"has_attachments": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
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
					Default:  false,
				},
				"is_automatic_forward": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"is_automatic_reply": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"is_encrypted": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"is_meeting_request": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"is_meeting_response": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"is_non_delivery_report": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"is_permission_controlled": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"is_read_receipt": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"is_signed": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"is_voicemail": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
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
					Default:  false,
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
					Default:  false,
				},
				"sent_only_to_me": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
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
					Default:  false,
				},
				"sent_to_or_cc_me": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
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
				Required: true,
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
						//TODO: support this after implement [categories](https://docs.microsoft.com/en-us/graph/api/resources/outlookcategory?view=graph-rest-1.0)
						//                         "assign_categories": {
						//                             Type:     schema.TypeSet,
						//                             MinItems: 1,
						//                             Optional: true,
						//                             Elem:     &schema.Schema{Type: schema.TypeString},
						//                         },
						"copy_to_folder": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"delete": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"forward_as_attachment_to": {
							Type:     schema.TypeSet,
							MinItems: 1,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"forward_to": {
							Type:     schema.TypeSet,
							MinItems: 1,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"mark_as_read": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
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
						},
						"move_to_folder": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"permanent_delete": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"redirect_to": {
							Type:     schema.TypeSet,
							MinItems: 1,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"stop_processing_rules": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
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

	param := &msgraph.MessageRule{
		DisplayName: utils.String(name),
		Sequence:    utils.Int(d.Get("sequence").(int)),
		IsEnabled:   utils.Bool(d.Get("enabled").(bool)),
		Conditions:  expandMessageRulePredicate(d.Get("condition").([]interface{})),
		Exceptions:  expandMessageRulePredicate(d.Get("exception").([]interface{})),
		Actions:     expandMessageRuleAction(d.Get("action").([]interface{})),
	}

	resp, err := client.Request().Add(ctx, param)
	if err != nil {
		return diag.Errorf("creating Message Rule %q: %w", name, err)
	}

	if resp.ID == nil {
		return diag.Errorf("nil ID for Message Rule %q", name)
	}
	d.SetId(*resp.ID)

	return resourceMessageRuleRead(ctx, d, meta)
}

func resourceMessageRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).MessageRules

	resp, err := client.ID(d.Id()).Request().Get(ctx)
	if err != nil {
		if utils.ResponseErrorWasNotFound(err) {
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
		return diag.Errorf(`setting "condition": %w"`, err)
	}
	if err := d.Set("exception", flattenMessageRulePredicate(resp.Exceptions	)); err != nil {
		return diag.Errorf(`setting "exception": %w"`, err)
	}
	if err := d.Set("action", flattenMessageRuleAction(resp.Actions	)); err != nil {
		return diag.Errorf(`setting "action": %w"`, err)
	}

	return nil
}


func resourceMessageRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).MessageRules.ID(d.Id())

	var param msgraph.MessageRule

	if d.HasChange("sequence") {
		param.Sequence = utils.Int(d.Get("sequence").(int))
	}
	if d.HasChange("enabled") {
		param.IsEnabled	 = utils.Bool(d.Get("enabled").(bool))
	}
	if d.HasChange("condition") {
		param.Conditions	 = expandMessageRulePredicate(d.Get("condition").([]interface{})),
	}
	if d.HasChange("exception") {
		param.Exceptions	 = expandMessageRulePredicate(d.Get("exception").([]interface{})),
	}
	if d.HasChange("action") {
		param.Actions	 = expandMessageRuleAction(d.Get("action").([]interface{})),
	}


	if err := client.Request().Update(ctx, &param); err != nil {
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

}

func expandMessageRuleAction(input []interface{}) *msgraph.MessageRuleActions {
}

func flattenMessageRulePredicate(input *msgraph.MessageRulePredicates) interface{} {
}

func flattenMessageRuleAction(input *msgraph.MessageRuleActions) interface{} {
}

