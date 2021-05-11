package templates

import (
	"errors"
	"html/template"
	"net/http"
	"strings"

	"github.com/fakovacic/aig/internal/meta"
	"github.com/fakovacic/aig/pkg/date"
	"github.com/fakovacic/aig/pkg/i18"
)

// TemplateFuncs return function on template
func TemplateFuncs() map[string]interface{} {
	return template.FuncMap{
		"i18":        i18.Get,
		"slice":      StringSlice,
		"prevPage":   meta.PrevPage,
		"nextPage":   meta.NextPage,
		"formatDate": date.Format,
	}
}

// New return new template struct
func New(name string, base *template.Template, files ...string) *Template {
	return &Template{
		Name:  name,
		Files: files,
		Base:  base,
	}
}

// Template struct
type Template struct {
	Name     string
	Base     *template.Template
	Template *template.Template
	Files    []string
}

// Parse parse template files
func (t *Template) Parse() error {
	var err error

	if t.Base != nil {
		t.Template, err = template.Must(t.Base.Clone()).
			Funcs(TemplateFuncs()).
			ParseFiles(t.Files...)
		if err != nil {
			return err
		}
		return nil
	}

	t.Template, err = template.New(t.Name).
		Funcs(TemplateFuncs()).
		ParseFiles(t.Files...)
	if err != nil {
		return err
	}
	return nil
}

// Execute execute template
func (t *Template) Execute(w http.ResponseWriter, data interface{}) error {
	if t.Template == nil {
		return errors.New("template empty")
	}

	err := t.Template.Execute(w, data)
	if err != nil {
		return err
	}
	return nil
}

// StringSlice join strings
func StringSlice(s []string) string {
	return strings.Join(s, ",")
}
