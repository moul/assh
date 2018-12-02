package templates // import "moul.io/assh/pkg/templates"

import (
	"encoding/json"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
)

var funcMap = template.FuncMap{
	"json": func(v interface{}) string {
		a, err := json.Marshal(v)
		if err != nil {
			return err.Error()
		}
		return string(a)
	},
	"prettyjson": func(v interface{}) string {
		a, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return err.Error()
		}
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
