package formulate

import (
	"fmt"

	"github.com/go-humble/temple/temple"

	"honnef.co/go/js/dom"
)

type ListCol struct {
	Heading string
	Model   string
}

type ListForm struct {
	Title    string
	Icon     string
	ID       int
	Data     interface{}
	Cols     []*ListCol
	RowCB    func(string)
	CancelCB func(dom.Event)
	NewRowCB func(dom.Event)
}

// Init a new listform
func (f *ListForm) New(icon string, title string) *ListForm {
	f.Title = title
	f.Icon = icon
	return f
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

// Render the form
func (f *ListForm) Render(name string, selector string, data interface{}) {

	f.Data = data

	w := dom.GetWindow()
	doc := w.Document()

	renderTemplateT(f.generateTemplate(name), selector, f)

	if f.CancelCB == nil {
		print("Error - No cancel callback")
		return
	}
	if f.RowCB == nil {
		print("Error - No save callback")
		return
	}

	// If there is a focusfield, then focus on it
	if el := doc.QuerySelector("#focusme"); el != nil {
		el.(*dom.HTMLInputElement).Focus()
	}

	// plug in cancel callbacks
	if el := doc.QuerySelector("#legend"); el != nil {
		el.AddEventListener("click", false, f.CancelCB)
	}

	if el := doc.QuerySelector(".md-close"); el != nil {
		el.AddEventListener("click", false, f.CancelCB)
	}

	if el := doc.QuerySelector(".data-add-btn"); el != nil {
		if f.NewRowCB != nil {
			el.AddEventListener("click", false, f.NewRowCB)
		}
	}

	// Handlers on the table itself
	if el := doc.QuerySelector(".data-table"); el != nil {

		el.AddEventListener("click", false, func(evt dom.Event) {
			evt.PreventDefault()
			td := evt.Target()
			tr := td.ParentElement()
			key := tr.GetAttribute("key")
			f.RowCB(key)
		})

		el.AddEventListener("keyup", false, func(evt dom.Event) {
			if evt.(*dom.KeyboardEvent).KeyCode == 27 {
				evt.PreventDefault()
				el.AddEventListener("click", false, f.CancelCB)
			}
		})
	}

}

func (f *ListForm) generateTemplate(name string) *temple.Template {

	tmpl, err := generatedTemplates.GetTemplate(name)
	if err != nil {
		print("Generating template for", name)
		// Template doesnt exist, so create it
		src := `
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

<table class="data-table" id="list-form">
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
			src += fmt.Sprintf("<td>{{if .%s}}{{.%s}}{{end}}</td>\n", col.Model, col.Model)
		}

		src += `      
    </tr>
{{end}}  
  <tbody>
  </tbody>
</table>

</div>	
`
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
