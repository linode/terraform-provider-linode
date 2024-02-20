package customtypes

//
//import (
//	"context"
//	"fmt"
//	"github.com/hashicorp/terraform-plugin-log/tflog"
//
//	"github.com/hashicorp/terraform-plugin-framework/attr"
//	"github.com/hashicorp/terraform-plugin-framework/diag"
//	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
//	"github.com/hashicorp/terraform-plugin-go/tftypes"
//)
//
//// Ensure the implementation satisfies the expected interfaces
//var (
//	_ basetypes.SetTypable                    = CaseInsensitiveSetType{}
//	_ basetypes.SetValuableWithSemanticEquals = CaseInsensitiveSetValue{}
//)
//
//type CaseInsensitiveSetType struct {
//	basetypes.SetType
//}
//
//func (t CaseInsensitiveSetType) Equal(o attr.Type) bool {
//	other, ok := o.(CaseInsensitiveSetType)
//
//	if !ok {
//		return false
//	}
//
//	return t.SetType.Equal(other.SetType)
//}
//
//func (t CaseInsensitiveSetType) String() string {
//	return "CaseInsensitiveSetType"
//}
//a
//
//func (t CaseInsensitiveSetType) ValueFromString(
//	ctx context.Context,
//	in basetypes.SetValue,
//) (basetypes.SetValuable, diag.Diagnostics) {
//	// CaseInsensitiveSetValue defined in the value type section
//	value := CaseInsensitiveSetValue{
//		SetValue: in,
//	}
//
//	return value, nil
//}
//
//func (t CaseInsensitiveSetType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
//	attrValue, err := t.SetType.ValueFromTerraform(ctx, in)
//	if err != nil {
//		return nil, err
//	}
//
//	setValue, ok := attrValue.(basetypes.SetValue)
//
//	if !ok {
//		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
//	}
//
//	stringValuable, diags := t.ValueFromSet(ctx, setValue)
//
//	if diags.HasError() {
//		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
//	}
//
//	return stringValuable, nil
//}
//
//func (t CaseInsensitiveSetType) ValueType(ctx context.Context) attr.Value {
//	return CaseInsensitiveSetValue{}
//}
//
//var _ basetypes.SetValuable = CaseInsensitiveSetValue{}
//
//type CaseInsensitiveSetValue struct {
//	basetypes.SetValue
//}
//
//func (v CaseInsensitiveSetValue) Equal(o attr.Value) bool {
//	other, ok := o.(CaseInsensitiveSetValue)
//
//	if !ok {
//		return false
//	}
//
//	return v.SetValue.Equal(other.SetValue)
//}
//
//func (v CaseInsensitiveSetValue) Type(ctx context.Context) attr.Type {
//	return CaseInsensitiveSetType{}
//}
//
//// CaseInsensitiveSetValue defined in the value type section
//// Ensure the implementation satisfies the expected interfaces
//
//func (v CaseInsensitiveSetValue) SetSemanticEquals(
//	ctx context.Context,
//	newValuable basetypes.SetValuable,
//) (bool, diag.Diagnostics) {
//	//var diags diag.Diagnostics
//
//	//newValue, ok := newValuable.(CaseInsensitiveSetValue)
//	//if !ok {
//	//	diags.AddError(
//	//		"Semantic Equality Check Error",
//	//		"An unexpected value type was received while performing semantic equality checks. "+
//	//			"Please report this to the provider developers.\n\n"+
//	//			"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
//	//			"Got Value Type: "+fmt.Sprintf("%T", newValuable),
//	//	)
//	//
//	//	return false, diags
//	//}
//
//	tflog.Info(ctx, "RGIREIONJREIOHNRE")
//	return false, nil
//}
