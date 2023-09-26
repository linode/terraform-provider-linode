package helper

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func EmptySetDefault(elemType attr.Type) defaults.Set {
	return setdefault.StaticValue(
		types.SetValueMust(
			elemType,
			[]attr.Value{},
		),
	)
}
