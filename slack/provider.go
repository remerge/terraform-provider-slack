package slack

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			desc = fmt.Sprintf("Defaults to `%v`. ", s.Default) + desc
		}
		return strings.TrimSpace(desc)
	}
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"token": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("SLACK_TOKEN", nil),
					Description: "The OAuth token used to connect to Slack.",
				},
			},

			DataSourcesMap: map[string]*schema.Resource{
				"slack_conversation": dataSourceConversation(),
				"slack_user":         dataSourceSlackUser(),
				"slack_usergroup":    dataSourceUserGroup(),
			},

			ResourcesMap: map[string]*schema.Resource{
				"slack_conversation":       resourceSlackConversation(),
				"slack_usergroup_channels": resourceSlackUserGroupChannels(),
				"slack_usergroup_members":  resourceSlackUserGroupMembers(),
				"slack_usergroup":          resourceSlackUserGroup(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var diags diag.Diagnostics
		var config = Config{}

		if token, ok := d.GetOk("token"); ok {
			config.Token = token.(string)
		}

		client, err := config.Client()
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Unable to create client from config: %v", err),
			})
		}

		return client, diags
	}
}
