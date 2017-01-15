package templates

import (
	"encoding/json"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
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
	"join":  strings.Join,
	"title": strings.Title,
	"lower": strings.ToLower,
	"upper": strings.ToUpper,
}

func init() {
	for k, v := range sprig.TxtFuncMap() {
		funcMap[k] = v
	}
}

// New creates a new template with funcMap and parses the given format.
func New(format string) (*template.Template, error) {
	return template.New("").Funcs(funcMap).Parse(format)
}
