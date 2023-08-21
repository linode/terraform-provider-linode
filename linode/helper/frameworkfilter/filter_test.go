//go:build unit

package frameworkfilter

var testFilterConfig = Config{
	"foo":          {APIFilterable: false, TypeFunc: FilterTypeString},
	"bar":          {APIFilterable: false, TypeFunc: FilterTypeString},
	"api_foo":      {APIFilterable: true, TypeFunc: FilterTypeString},
	"api_bar":      {APIFilterable: true, TypeFunc: FilterTypeString},
	"api_foo_int":  {APIFilterable: true, TypeFunc: FilterTypeInt},
	"api_foo_bool": {APIFilterable: true, TypeFunc: FilterTypeBool},
}
