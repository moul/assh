package templates

import (
	"encoding/json"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func funcMap() template.FuncMap {
	var m = template.FuncMap{
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
		"title": cases.Title(language.Und, cases.NoLower).String,
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
	}
	for k, v := range sprig.TxtFuncMap() {
		m[k] = v
	}
	return m
}

// New creates a new template with funcMap and parses the given format.
func New(format string) (*template.Template, error) {
	return template.New("").Funcs(funcMap()).Parse(format)
}
