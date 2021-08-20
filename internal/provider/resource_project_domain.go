package provider

import (
	"context"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel"
	pdomain "github.com/chronark/terraform-provider-vercel/pkg/vercel/project_domain"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProjectDomain() *schema.Resource {
	return &schema.Resource{
		Description: "https://vercel.com/docs/api#endpoints/projects/get-a-single-project-domain",

		CreateContext: resourceProjectDomainCreate,
		ReadContext:   resourceProjectDomainRead,
		DeleteContext: resourceProjectDomainDelete,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Description: "The unique Project identifier.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "The name of the project domain.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},
			"redirect": {
				Description: "Target destination domain for redirect",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"redirect_status_code": {
				Description: "The redirect status code (301, 302, 307, 308).",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"git_branch": {
				Description: "Git branch for the domain to be auto assigned to. The Project's production branch is the default (null).",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"created_at": {
				Description: "A number containing the project domain when the variable was created in milliseconds.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"updated_at": {
				Description: "A number containing the project domain when the variable was updated in milliseconds.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func resourceProjectDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*vercel.Client)

	domain, err := client.ProjectDomain.Read(d.Get("project_id").(string), d.Get("name").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(domain.Name)

	if err := d.Set("name", domain.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("git_branch", domain.GitBranch); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("redirect", domain.Redirect); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("redirect_status_code", domain.RedirectStatusCode); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("updated_at", domain.UpdatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", domain.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceProjectDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*vercel.Client)

	dto := pdomain.CreateOrUpdateProjectDomain{
		Name:     d.Get("name").(string),
		Redirect: d.Get("redirect").(string),
		// RedirectStatusCode: d.Get("redirect_status_code").(int),
	}

	if rsc := d.Get("redirect_status_code").(int); rsc > 0 {
		dto.RedirectStatusCode = &rsc
	}

	if b := d.Get("git_branch").(string); b != "" {
		dto.GitBranch = &b
	}

	if _, err := client.ProjectDomain.Create(d.Get("project_id").(string), dto); err != nil {
		return diag.FromErr(err)
	}

	return resourceProjectDomainRead(ctx, d, meta)
}

func resourceProjectDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*vercel.Client)

	if err := client.ProjectDomain.Delete(d.Get("project_id").(string), d.Get("name").(string)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diag.Diagnostics{}
}
