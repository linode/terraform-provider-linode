package frameworkfilter

var testFilterConfig = Config{
	"foo":     {APIFilterable: false},
	"bar":     {APIFilterable: false},
	"api_foo": {APIFilterable: true},
	"api_bar": {APIFilterable: true},
}
