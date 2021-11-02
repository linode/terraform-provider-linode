package stackscript

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
)

func setStackScriptUserDefinedFields(d *schema.ResourceData, ss *linodego.Stackscript) {
	if ss.UserDefinedFields == nil {
		d.Set("user_defined_fields", nil)
		return
	}

	udfs := GetStackScriptUserDefinedFields(ss)
	d.Set("user_defined_fields", udfs)
}

func GetStackScriptUserDefinedFields(ss *linodego.Stackscript) []map[string]string {
	if ss.UserDefinedFields == nil {
		return nil
	}

	result := make([]map[string]string, len(*ss.UserDefinedFields))
	for i, udf := range *ss.UserDefinedFields {
		result[i] = map[string]string{
			"default": udf.Default,
			"example": udf.Example,
			"many_of": udf.ManyOf,
			"one_of":  udf.OneOf,
			"label":   udf.Label,
			"name":    udf.Name,
		}
	}

	return result
}
