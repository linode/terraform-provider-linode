package instancetype

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var dataSourceSchema = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeString,
		Description: "The unique ID assigned to this Instance type.",
		Required:    true,
	},
	"label": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The Linode Type's label is for display purposes only.",
		Computed:    true,
	},
	"disk": {
		Type:        schema.TypeInt,
		Description: "The Disk size, in MB, of the Linode Type.",
		Computed:    true,
	},
	"class": {
		Type: schema.TypeString,
		Description: "The class of the Linode Type. There are currently three classes of Linodes: nanode, " +
			"standard, highmem, dedicated",
		Computed: true,
	},
	"price": {
		Type:        schema.TypeList,
		Description: "Cost in US dollars, broken down into hourly and monthly charges.",
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"hourly": {
					Type:        schema.TypeFloat,
					Description: "Cost (in US dollars) per hour.",
					Computed:    true,
				},
				"monthly": {
					Type:        schema.TypeFloat,
					Description: "Cost (in US dollars) per month.",
					Computed:    true,
				},
			},
		},
	},
	"addons": {
		Type:        schema.TypeList,
		Description: "",
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"backups": {
					Type:        schema.TypeList,
					Description: "Information about the optional Backup service offered for Linodes.",
					Computed:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"price": {
								Type:        schema.TypeList,
								Description: "Cost of enabling Backups for this Linode Type.",
								Computed:    true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"hourly": {
											Type:        schema.TypeFloat,
											Description: "The cost (in US dollars) per hour to add Backups service.",
											Computed:    true,
										},
										"monthly": {
											Type:        schema.TypeFloat,
											Description: "The cost (in US dollars) per month to add Backups service.",
											Computed:    true,
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
	"network_out": {
		Type:        schema.TypeInt,
		Description: "The Mbits outbound bandwidth allocation.",
		Computed:    true,
	},
	"memory": {
		Type:        schema.TypeInt,
		Description: "Amount of RAM included in this Linode Type.",
		Computed:    true,
	},
	"transfer": {
		Type:        schema.TypeInt,
		Description: "The monthly outbound transfer amount, in MB.",
		Computed:    true,
	},
	"vcpus": {
		Type:        schema.TypeInt,
		Description: "The number of VCPU cores this Linode Type offers.",
		Computed:    true,
	},
}
