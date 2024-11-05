package tmpl

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type ResourceTemplateData struct {
	Prefix   string
	ID       string
	PubKey   string
	Region   string
	RootPass string
}

type TemplateData struct {
	Instances     []ResourceTemplateData
	NodeBalancers []ResourceTemplateData

	Label string
}

func Basic(t testing.TB, label, devicePrefix, region string) string {
	resources := []ResourceTemplateData{
		{
			Prefix:   devicePrefix,
			ID:       "one",
			PubKey:   acceptance.PublicKeyMaterial,
			Region:   region,
			RootPass: acctest.RandString(64),
		},
	}

	return acceptance.ExecuteTemplate(t,
		"firewall_basic", TemplateData{
			Label:         label,
			Instances:     resources,
			NodeBalancers: resources,
		})
}

func Updates(t testing.TB, label, devicePrefix, region string) string {
	resources := []ResourceTemplateData{
		{
			Prefix:   devicePrefix,
			ID:       "one",
			PubKey:   acceptance.PublicKeyMaterial,
			Region:   region,
			RootPass: acctest.RandString(64),
		},
		{
			Prefix:   devicePrefix,
			ID:       "two",
			PubKey:   acceptance.PublicKeyMaterial,
			Region:   region,
			RootPass: acctest.RandString(64),
		},
	}

	return acceptance.ExecuteTemplate(t,
		"firewall_updates", TemplateData{
			Label:         label,
			Instances:     resources,
			NodeBalancers: resources,
		})
}

func Minimum(t testing.TB, label string) string {
	return acceptance.ExecuteTemplate(t,
		"firewall_minimum", TemplateData{
			Label: label,
		})
}

func MultipleRules(t testing.TB, label, devicePrefix, region string) string {
	return acceptance.ExecuteTemplate(t,
		"firewall_multiple_rules", TemplateData{
			Label: label,
			Instances: []ResourceTemplateData{
				{
					Prefix:   devicePrefix,
					ID:       "one",
					PubKey:   acceptance.PublicKeyMaterial,
					Region:   region,
					RootPass: acctest.RandString(64),
				},
			},
		})
}

func NoDevice(t testing.TB, label string) string {
	return acceptance.ExecuteTemplate(t,
		"firewall_no_device", TemplateData{
			Label: label,
		})
}

func NoIPv6(t testing.TB, label string) string {
	return acceptance.ExecuteTemplate(t,
		"firewall_no_ipv6", TemplateData{
			Label: label,
		})
}

func NoRules(t testing.TB, label string) string {
	return acceptance.ExecuteTemplate(t,
		"firewall_no_rules", TemplateData{
			Label: label,
		})
}

func DataBasic(t testing.TB, label, devicePrefix, region string) string {
	resources := []ResourceTemplateData{
		{
			Prefix:   devicePrefix,
			ID:       "one",
			PubKey:   acceptance.PublicKeyMaterial,
			Region:   region,
			RootPass: acctest.RandString(64),
		},
	}

	return acceptance.ExecuteTemplate(t,
		"firewall_data_basic", TemplateData{
			Label:         label,
			Instances:     resources,
			NodeBalancers: resources,
		})
}
