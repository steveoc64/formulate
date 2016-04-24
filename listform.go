package formulate

import (
	"fmt"

	"github.com/go-humble/temple/temple"

	"honnef.co/go/js/dom"
)

type ListCol struct {
	Heading string
	Model   string
	Format  string
	Width   string
}

type ListForm struct {
	Title       string
	Icon        string
	ID          int
	Data        interface{}
	Cols        []*ListCol
	RowCB       func(string)
	CancelCB    func(dom.Event)
	NewRowCB    func(dom.Event)
	HasSetWidth bool
}

// Init a new listform
func (f *ListForm) New(icon string, title string) *ListForm {
	f.Title = title
	f.Icon = icon
	return f
}

func (f *ListForm) SetWidths(w []string) {
	f.HasSetWidth = true
	for i, v := range w {
		if i < len(f.Cols) {
			f.Cols[i].Width = v
		}
	}
}

// Associate a cancel event with the listform
func (f *ListForm) CancelEvent(c func(dom.Event)) *ListForm {
	f.CancelCB = c
	return f
}

// Add new row
func (f *ListForm) NewRowEvent(c func(dom.Event)) *ListForm {
	f.NewRowCB = c
	return f
}

// Associate a Row Click event with the listform
func (f *ListForm) RowEvent(c func(string)) *ListForm {
	f.RowCB = c
	return f
}

// Add a colunm to the listform
func (f *ListForm) Column(heading string, model string) *ListForm {
	c := &ListCol{
		Heading: heading,
		Model:   model,
	}
	f.Cols = append(f.Cols, c)
	return f
}

// Add a colunm to the listform with format
func (f *ListForm) ColumnFormat(heading string, model string, format string) *ListForm {
	c := &ListCol{
		Heading: heading,
		Model:   model,
		Format:  format,
	}
	f.Cols = append(f.Cols, c)
	return f
}

// Add a colunm to the listform in Date Format
func (f *ListForm) DateColumn(heading string, model string) *ListForm {
	c := &ListCol{
		Heading: heading,
		Model:   model,
		Format:  "date",
	}
	f.Cols = append(f.Cols, c)
	return f
}

// Render the form using a template that we generate on the fly
func (f *ListForm) Render(name string, selector string, data interface{}) {

	f.Data = data
	renderTemplateT(f.generateTemplate(name), selector, f)
	f.decorate()
}

// Render the form using a custom template
func (f *ListForm) RenderCustom(name string, selector string, data interface{}) {

	f.Data = data
	renderTemplate(name, selector, data)
	f.decorate()
}

func (f *ListForm) decorate() {

	w := dom.GetWindow()
	doc := w.Document()

	// If there is a focusfield, then focus on it
	if el := doc.QuerySelector("#focusme"); el != nil {
		el.(*dom.HTMLInputElement).Focus()
	}

	// plug in cancel callbacks
	if f.CancelCB != nil {
		if el := doc.QuerySelector("#legend"); el != nil {
			el.AddEventListener("click", false, f.CancelCB)
		}

		if el := doc.QuerySelector(".md-close"); el != nil {
			el.AddEventListener("click", false, f.CancelCB)
		}
	}

	if f.NewRowCB != nil {
		if el := doc.QuerySelector(".data-add-btn"); el != nil {
			el.AddEventListener("click", false, f.NewRowCB)
		}
	}

	// Handlers on the table itself
	if el := doc.QuerySelector(".data-table"); el != nil {

		if f.RowCB != nil {
			el.AddEventListener("click", false, func(evt dom.Event) {
				evt.PreventDefault()
				td := evt.Target()
				tr := td.ParentElement()
				key := tr.GetAttribute("key")
				f.RowCB(key)
			})
		}

		if f.CancelCB != nil {
			el.AddEventListener("keyup", false, func(evt dom.Event) {
				if evt.(*dom.KeyboardEvent).KeyCode == 27 {
					evt.PreventDefault()
					el.AddEventListener("click", false, f.CancelCB)
				}
			})
		}
	}

}

func (f *ListForm) generateTemplate(name string) *temple.Template {

	tmpl, err := generatedTemplates.GetTemplate(name)
	if err != nil {
		print("Generating template for", name)
		// Template doesnt exist, so create it

		src := ""
		doTitle := false

		if f.Title != "" || f.Icon != "" {
			doTitle = true
		}

		if doTitle {

			src += `
<div class="data-container">
	<div class="row data-table-header">
    <h3 class="column column-90" id="legend">
      <i class="fa {{.Icon}} fa-lg" style="font-size: 3rem"></i> 
      {{.Title}}
    </h3>
`
			if f.NewRowCB != nil {
				src += `
    <div class="column col-center">
      <i class="data-add-btn fa fa-plus-circle fa-lg"></i>    
    </div>    
`
			}
			src += `    
  </div>
`
		}

		src += `<table class="data-table" id="list-form">
  <thead>
    <tr>
      {{range .Cols}}
      <th>{{.Heading}}</th>
      {{end}}
    </tr>
  </thead>
  <tbody>
{{$cols := .Cols}}
{{range .Data}}  
    <tr class="data-row" 
        key="{{.ID}}">
`
		// for each column, add a column renderer

		for _, col := range f.Cols {
			width := ""
			if f.HasSetWidth {
				width = fmt.Sprintf(` width="%s"`, col.Width)
			}

			src += fmt.Sprintf("<td %s %s>{{if .%s}}{{.%s}}{{end}}</td>\n", width, col.Format, col.Model, col.Model)
		}

		src += `      
    </tr>
{{end}}  
  <tbody>
  </tbody>
</table>`

		if doTitle {
			src += `
</div>
`
		}

		// print("list source = ", src)
		createErr := generatedTemplates.AddTemplate(name, src)
		if createErr != nil {
			print("failed to create template", name, createErr.Error())
		}
		tmpl, err = generatedTemplates.GetTemplate(name)
		if err != nil {
			print("could not get generated template !!")
		}

	}
	return tmpl

}
