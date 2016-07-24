package templates

import (
	"encoding/json"
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"json": func(v interface{}) string {
		a, _ := json.Marshal(v)
		return string(a)
	},
	"prettyjson": func(v interface{}) string {
		a, _ := json.MarshalIndent(v, "", "  ")
		return string(a)
	},
	// yaml
	// xml
	// toml
	"split": strings.Split,
	"join":  strings.Join,
	"title": strings.Title,
	"lower": strings.ToLower,
	"upper": strings.ToUpper,
}

// New creates a new template with funcMap and parses the given format.
func New(format string) (*template.Template, error) {
	return template.New("").Funcs(funcMap).Parse(format)
}
