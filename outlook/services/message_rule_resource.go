package services

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

func ResourceMessageRule() *schema.Resource {
	conditionSchema := &schema.Schema{
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
			State: schema.ImportStatePassthrough,
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
			"condition": conditionSchema,
			"exception": conditionSchema,
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
}

func resourceMessageRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
}

func resourceMessageRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
}

func resourceMessageRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
}
