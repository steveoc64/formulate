package formulate

import (
	"errors"
	"fmt"
	"html/template"

	"github.com/go-humble/temple/temple"
	"honnef.co/go/js/dom"
)

var gt func(name string) (*temple.Template, error)
var generatedTemplates *temple.Group

func Templates(g func(string) (*temple.Template, error)) {
	gt = g
	generatedTemplates = temple.NewGroup()
}

// Load a template and attach it to the specified element in the doc
func renderTemplate(name string, selector string, data interface{}) error {

	t, err := gt(name)
	if t == nil {
		print("Failed to load template", name)
		return errors.New("Invalid template")
	}
	if err != nil {
		print(err.Error())
		return err
	}

	return renderTemplateT(t, selector, data)
}

// Load a template and attach it to the specified element in the doc
func renderTemplateT(t *temple.Template, selector string, data interface{}) error {
	w := dom.GetWindow()
	doc := w.Document()

	funcMap := template.FuncMap{
		"Truncate": func(s string) string {
			if len(s) > 120 {
				return fmt.Sprintf("%s ...", s[:120])
			}
			return s
		},
	}

	print(funcMap)

	el := doc.QuerySelector(selector)
	if el == nil {
		print("Could not find selector", selector)
		return errors.New("Invalid selector")
	}

	if err := t.ExecuteEl(el, data); err != nil {
		print(err.Error())
		return err
	}
	return nil
}
