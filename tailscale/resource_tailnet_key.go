package tailscale

import (
	"context"

	"github.com/davidsbond/tailscale-client-go/tailscale"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTailnetKey() *schema.Resource {
	return &schema.Resource{
		Description:   "The tailnet_key resource allows you to create pre-authentication keys that can register new nodes without needing to sign in via a web browser. See https://tailscale.com/kb/1085/auth-keys for more information",
		ReadContext:   resourceTailnetKeyRead,
		CreateContext: resourceTailnetKeyCreate,
		DeleteContext: resourceTailnetKeyDelete,
		Schema: map[string]*schema.Schema{
			"reusable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the key is reusable or single-use.",
				ForceNew:    true,
			},
			"ephemeral": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the key is ephemeral.",
				ForceNew:    true,
			},
			"tags": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "List of tags to apply to the machines authenticated by the key.",
				ForceNew:    true,
			},
			"key": {
				Type:        schema.TypeString,
				Description: "The authentication key",
				Computed:    true,
				Sensitive:   true,
			},
			"id": {
				Type:        schema.TypeString,
				Description: "The key's identifier",
				Computed:    true,
			},
		},
	}
}

func resourceTailnetKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*tailscale.Client)
	reusable := d.Get("reusable").(bool)
	ephemeral := d.Get("ephemeral").(bool)
	var tags []string
	for _, tag := range d.Get("tags").(*schema.Set).List() {
		tags = append(tags, tag.(string))
	}

	var capabilities tailscale.KeyCapabilities
	capabilities.Devices.Create.Reusable = reusable
	capabilities.Devices.Create.Ephemeral = ephemeral
	capabilities.Devices.Create.Tags = tags

	key, err := client.CreateKey(ctx, capabilities)
	if err != nil {
		return diagnosticsError(err, "Failed to create key")
	}

	d.SetId(key.ID)

	if err = d.Set("key", key.Key); err != nil {
		return diagnosticsError(err, "Failed to set key")
	}

	if err = d.Set("id", key.ID); err != nil {
		return diagnosticsError(err, "Failed to set id")
	}

	return nil
}

func resourceTailnetKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*tailscale.Client)
	id := d.Get("id").(string)

	err := client.DeleteKey(ctx, id)
	switch {
	case tailscale.IsNotFound(err):
		// Single-use keys may no longer be here, so we can ignore deletions that fail due to not-found errors.
		return nil
	case err != nil:
		return diagnosticsError(err, "Failed to delete key")
	default:
		return nil
	}
}

func resourceTailnetKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*tailscale.Client)
	id := d.Get("id").(string)
	key, err := client.GetKey(ctx, id)

	reusable := d.Get("reusable").(bool)

	switch {
	case tailscale.IsNotFound(err) && !reusable:
		// If we get a 404 on a one-off key, don't return an error here.
		return nil
	case err != nil:
		return diagnosticsError(err, "Failed to fetch key")
	}

	d.SetId(key.ID)
	if err = d.Set("reusable", key.Capabilities.Devices.Create.Reusable); err != nil {
		return diagnosticsError(err, "Failed to set reusable")
	}

	if err = d.Set("ephemeral", key.Capabilities.Devices.Create.Ephemeral); err != nil {
		return diagnosticsError(err, "failed to set ephemeral")
	}

	return nil
}
