package instance

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/linode/linodego"
)

// linodeInterfacesSchema defines the schema for the `linode_interfaces` block
// on the instance resource. These are the new "Linode Interfaces" (as opposed to
// legacy config interfaces), and they additionally support `rdma_vpc` for
// GPUDirect RDMA capable Linode plans.
//
// NOTE: These interfaces can only be specified at instance creation time.
// Modifying them in-place requires using the standalone `linode_interface`
// resource (which does not support creating RDMA VPC interfaces).
var linodeInterfacesSchema = &schema.Schema{
	Type: schema.TypeList,
	Description: "An array of new-generation Linode Interfaces to attach to this Linode at " +
		"creation. Supports `public`, `vlan`, `vpc`, and `rdma_vpc` interface types. " +
		"At most one of `public`, `vlan`, `vpc`, or `rdma_vpc` can be specified per interface entry. " +
		"NOTE: This option requires `interface_generation = \"linode\"`.",
	Optional:      true,
	ForceNew:      true,
	ConflictsWith: []string{"interface", "config", "disk"},
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"firewall_id": {
				Type:        schema.TypeInt,
				Description: "The ID of an enabled firewall to attach to this interface. Not allowed for VLAN interfaces.",
				Optional:    true,
				ForceNew:    true,
			},
			"default_route": {
				Type:        schema.TypeList,
				Description: "Default route configuration for the interface.",
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"ipv6": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"public": {
				Type:        schema.TypeList,
				Description: "Configuration for a Linode public interface.",
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"addresses": {
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
												"primary": {
													Type:     schema.TypeBool,
													Optional: true,
													ForceNew: true,
												},
											},
										},
									},
								},
							},
						},
						"ipv6": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ranges": {
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"range": {
													Type:     schema.TypeString,
													Required: true,
													ForceNew: true,
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
			"vlan": {
				Type:        schema.TypeList,
				Description: "Configuration for a Linode VLAN interface.",
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vlan_label": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringLenBetween(1, 64),
						},
						"ipam_address": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"vpc": {
				Type:        schema.TypeList,
				Description: "Configuration for a (regular) Linode VPC interface.",
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"ipv4": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"addresses": {
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
												"primary": {
													Type:     schema.TypeBool,
													Optional: true,
													ForceNew: true,
												},
												"nat_1_1_address": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
											},
										},
									},
									"ranges": {
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"range": {
													Type:     schema.TypeString,
													Required: true,
													ForceNew: true,
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
			"rdma_vpc": {
				Type: schema.TypeList,
				Description: "Configuration for a GPUDirect RDMA VPC interface. " +
					"RDMA VPC interfaces can only be attached at instance creation time. " +
					"NOTE: RDMA VPC interfaces may not currently be available to all users.",
				MaxItems: 1,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:        schema.TypeInt,
							Description: "The ID of the RDMA VPC subnet to attach this interface to.",
							Required:    true,
							ForceNew:    true,
						},
						"ipv4": {
							Type:        schema.TypeList,
							Description: "The IPv4 configuration for the RDMA VPC interface.",
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"addresses": {
										Type:        schema.TypeList,
										Description: "The list of IPv4 addresses for this RDMA VPC interface. Must contain exactly one element.",
										Required:    true,
										MaxItems:    1,
										ForceNew:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": {
													Type:        schema.TypeString,
													Description: "The IPv4 address (or 'auto' to allocate one from the subnet).",
													Optional:    true,
													ForceNew:    true,
													Default:     "auto",
												},
												"primary": {
													Type:        schema.TypeBool,
													Description: "Whether this is the primary IPv4 address for the interface.",
													Optional:    true,
													ForceNew:    true,
													Default:     true,
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
		},
	},
}

// expandLinodeInstanceInterfaces converts the TF-encoded `linode_interfaces`
// schema list into a slice of linodego.LinodeInstanceInterfaceCreateOptions
// suitable for InstanceCreateOptions.
func expandLinodeInstanceInterfaces(raw []any) ([]linodego.LinodeInstanceInterfaceCreateOptions, error) {
	result := make([]linodego.LinodeInstanceInterfaceCreateOptions, 0, len(raw))

	for i, entry := range raw {
		m, ok := entry.(map[string]any)
		if !ok {
			continue
		}

		var opt linodego.LinodeInstanceInterfaceCreateOptions
		typeCount := 0

		if v, ok := m["firewall_id"].(int); ok && v != 0 {
			opt.FirewallID = linodego.Pointer(v)
		}

		if dr, ok := m["default_route"].([]any); ok && len(dr) > 0 {
			if drMap, ok := dr[0].(map[string]any); ok {
				opt.DefaultRoute = expandLinodeDefaultRoute(drMap)
			}
		}

		if pub, ok := m["public"].([]any); ok && len(pub) > 0 {
			if pubMap, ok := pub[0].(map[string]any); ok {
				opt.Public = expandLinodePublicInterface(pubMap)
				typeCount++
			}
		}

		if vlan, ok := m["vlan"].([]any); ok && len(vlan) > 0 {
			if vlanMap, ok := vlan[0].(map[string]any); ok {
				opt.VLAN = expandLinodeVLANInterface(vlanMap)
				typeCount++
			}
		}

		if vpc, ok := m["vpc"].([]any); ok && len(vpc) > 0 {
			if vpcMap, ok := vpc[0].(map[string]any); ok {
				opt.VPC = expandLinodeVPCInterface(vpcMap)
				typeCount++
			}
		}

		if rdma, ok := m["rdma_vpc"].([]any); ok && len(rdma) > 0 {
			if rdmaMap, ok := rdma[0].(map[string]any); ok {
				opt.RDMAVPC = expandRDMAVPCInterface(rdmaMap)
				typeCount++
			}
		}

		if typeCount != 1 {
			return nil, fmt.Errorf(
				"linode_interfaces[%d]: exactly one of `public`, `vlan`, `vpc`, or `rdma_vpc` must be specified, got %d",
				i, typeCount,
			)
		}

		result = append(result, opt)
	}

	return result, nil
}

func expandLinodeDefaultRoute(m map[string]any) *linodego.InterfaceDefaultRoute {
	dr := &linodego.InterfaceDefaultRoute{}
	if v, ok := m["ipv4"].(bool); ok {
		dr.IPv4 = linodego.Pointer(v)
	}
	if v, ok := m["ipv6"].(bool); ok {
		dr.IPv6 = linodego.Pointer(v)
	}
	return dr
}

func expandLinodePublicInterface(m map[string]any) *linodego.PublicInterfaceCreateOptions {
	opts := &linodego.PublicInterfaceCreateOptions{}

	if ipv4, ok := m["ipv4"].([]any); ok && len(ipv4) > 0 {
		if ipv4Map, ok := ipv4[0].(map[string]any); ok {
			if addrs, ok := ipv4Map["addresses"].([]any); ok && len(addrs) > 0 {
				addresses := make([]linodego.PublicInterfaceIPv4AddressCreateOptions, 0, len(addrs))
				for _, a := range addrs {
					addrMap, ok := a.(map[string]any)
					if !ok {
						continue
					}
					addr := linodego.PublicInterfaceIPv4AddressCreateOptions{}
					if v, ok := addrMap["address"].(string); ok && v != "" {
						addr.Address = linodego.Pointer(v)
					}
					if v, ok := addrMap["primary"].(bool); ok {
						addr.Primary = linodego.Pointer(v)
					}
					addresses = append(addresses, addr)
				}
				opts.IPv4 = &linodego.PublicInterfaceIPv4CreateOptions{Addresses: &addresses}
			}
		}
	}

	if ipv6, ok := m["ipv6"].([]any); ok && len(ipv6) > 0 {
		if ipv6Map, ok := ipv6[0].(map[string]any); ok {
			if ranges, ok := ipv6Map["ranges"].([]any); ok && len(ranges) > 0 {
				rangeOpts := make([]linodego.PublicInterfaceIPv6RangeCreateOptions, 0, len(ranges))
				for _, r := range ranges {
					rMap, ok := r.(map[string]any)
					if !ok {
						continue
					}
					if v, ok := rMap["range"].(string); ok {
						rangeOpts = append(rangeOpts, linodego.PublicInterfaceIPv6RangeCreateOptions{Range: v})
					}
				}
				opts.IPv6 = &linodego.PublicInterfaceIPv6CreateOptions{Ranges: &rangeOpts}
			}
		}
	}

	return opts
}

func expandLinodeVLANInterface(m map[string]any) *linodego.VLANInterface {
	v := &linodego.VLANInterface{}
	if label, ok := m["vlan_label"].(string); ok {
		v.VLANLabel = label
	}
	if ipam, ok := m["ipam_address"].(string); ok && ipam != "" {
		v.IPAMAddress = linodego.Pointer(ipam)
	}
	return v
}

func expandLinodeVPCInterface(m map[string]any) *linodego.VPCInterfaceCreateOptions {
	opts := &linodego.VPCInterfaceCreateOptions{}
	if subnetID, ok := m["subnet_id"].(int); ok {
		opts.SubnetID = subnetID
	}

	if ipv4, ok := m["ipv4"].([]any); ok && len(ipv4) > 0 {
		if ipv4Map, ok := ipv4[0].(map[string]any); ok {
			ipv4Opts := &linodego.VPCInterfaceIPv4CreateOptions{}

			if addrs, ok := ipv4Map["addresses"].([]any); ok && len(addrs) > 0 {
				addrOpts := make([]linodego.VPCInterfaceIPv4AddressCreateOptions, 0, len(addrs))
				for _, a := range addrs {
					addrMap, ok := a.(map[string]any)
					if !ok {
						continue
					}
					ao := linodego.VPCInterfaceIPv4AddressCreateOptions{}
					if v, ok := addrMap["address"].(string); ok && v != "" {
						ao.Address = linodego.Pointer(v)
					}
					if v, ok := addrMap["primary"].(bool); ok {
						ao.Primary = linodego.Pointer(v)
					}
					if v, ok := addrMap["nat_1_1_address"].(string); ok && v != "" {
						ao.NAT1To1Address = linodego.Pointer(v)
					}
					addrOpts = append(addrOpts, ao)
				}
				ipv4Opts.Addresses = &addrOpts
			}

			if ranges, ok := ipv4Map["ranges"].([]any); ok && len(ranges) > 0 {
				rangeOpts := make([]linodego.VPCInterfaceIPv4RangeCreateOptions, 0, len(ranges))
				for _, r := range ranges {
					rMap, ok := r.(map[string]any)
					if !ok {
						continue
					}
					if v, ok := rMap["range"].(string); ok {
						rangeOpts = append(rangeOpts, linodego.VPCInterfaceIPv4RangeCreateOptions{Range: v})
					}
				}
				ipv4Opts.Ranges = &rangeOpts
			}

			opts.IPv4 = ipv4Opts
		}
	}

	return opts
}

func expandRDMAVPCInterface(m map[string]any) *linodego.RDMAVPCInterfaceCreateOptions {
	opts := &linodego.RDMAVPCInterfaceCreateOptions{}

	if subnetID, ok := m["subnet_id"].(int); ok {
		opts.SubnetID = subnetID
	}

	if ipv4, ok := m["ipv4"].([]any); ok && len(ipv4) > 0 {
		if ipv4Map, ok := ipv4[0].(map[string]any); ok {
			if addrs, ok := ipv4Map["addresses"].([]any); ok {
				addrOpts := make([]linodego.RDMAVPCInterfaceIPv4AddressOptions, 0, len(addrs))
				for _, a := range addrs {
					addrMap, ok := a.(map[string]any)
					if !ok {
						continue
					}
					ao := linodego.RDMAVPCInterfaceIPv4AddressOptions{}
					if v, ok := addrMap["address"].(string); ok {
						ao.Address = v
					}
					if v, ok := addrMap["primary"].(bool); ok {
						ao.Primary = v
					}
					addrOpts = append(addrOpts, ao)
				}
				opts.IPv4 = linodego.RDMAVPCInterfaceIPv4Options{Addresses: addrOpts}
			}
		}
	}

	return opts
}
