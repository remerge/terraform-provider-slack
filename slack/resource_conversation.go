package slack

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/slack-go/slack"
)

const (
	conversationActionOnDestroyNone    = "none"
	conversationActionOnDestroyArchive = "archive"
)

var validateConversationActionOnDestroyValue = validation.StringInSlice([]string{
	conversationActionOnDestroyNone,
	conversationActionOnDestroyArchive,
}, false)

func resourceSlackConversation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSlackConversationCreate,
		ReadContext:   resourceSlackConversationRead,
		UpdateContext: resourceSlackConversationUpdate,
		DeleteContext: resourceSlackConversationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_private": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"topic": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"purpose": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"action_on_destroy": {
				Type:         schema.TypeString,
				Description:  "Either of none or archive",
				Required:     true,
				ValidateFunc: validateConversationActionOnDestroyValue,
			},
			"is_shared": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_ext_shared": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_org_shared": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"creator": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSlackConversationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var client = meta.(*slack.Client)

	name := d.Get("name").(string)
	isPrivate := d.Get("is_private").(bool)

	log.Printf("[DEBUG] Creating conversation %q: %#v", d.Id(), name)
	channel, err := client.CreateConversationContext(ctx, name, isPrivate)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(channel.ID)

	log.Printf("[DEBUG] Finished creating conversation %q: %#v", d.Id(), d.Get("name").(string))

	return resourceSlackConversationRead(ctx, d, meta)
}

func resourceSlackConversationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var client = meta.(*slack.Client)

	log.Printf("[DEBUG] Reading conversation: %q: %#v", d.Id(), d.Id())

	channel, err := client.GetConversationInfoContext(ctx, d.Id(), false)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(channel.ID)

	d.Set("name", channel.Name)
	d.Set("topic", channel.Topic.Value)
	d.Set("purpose", channel.Purpose.Value)
	d.Set("is_ext_shared", channel.IsExtShared)
	d.Set("is_private", channel.IsPrivate)
	d.Set("is_shared", channel.IsShared)
	d.Set("is_org_shared", channel.IsOrgShared)
	d.Set("created", channel.Created)
	d.Set("creator", channel.Creator)

	log.Printf("[DEBUG] Finished reading conversation %q: %#v", d.Id(), d.Get("name").(string))

	return diags
}

func resourceSlackConversationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var client = meta.(*slack.Client)

	log.Printf("[DEBUG] Updating conversation %q: %#v", d.Id(), d.Get("name").(string))

	if d.HasChange("name") {
		if _, err := client.RenameConversationContext(ctx, d.Id(), d.Get("name").(string)); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("topic") {
		if topic, ok := d.GetOk("topic"); ok {
			if _, err := client.SetTopicOfConversationContext(ctx, d.Id(), topic.(string)); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("purpose") {
		if purpose, ok := d.GetOk("purpose"); ok {
			if _, err := client.SetPurposeOfConversationContext(ctx, d.Id(), purpose.(string)); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceSlackConversationRead(ctx, d, meta)
}

func resourceSlackConversationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var client = meta.(*slack.Client)

	action := d.Get("action_on_destroy").(string)

	switch action {
	case conversationActionOnDestroyNone:
		log.Printf("[DEBUG] Doing nothing about conversation %q: %#v", d.Id(), d.Get("name").(string))
	case conversationActionOnDestroyArchive:
		log.Printf("[DEBUG] Archiving conversation %q: %#v", d.Id(), d.Get("name").(string))
		if err := client.ArchiveConversationContext(ctx, d.Id()); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}
