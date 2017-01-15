package formulate

import (
	"fmt"

	"github.com/go-humble/temple/temple"

	"honnef.co/go/js/dom"
)

type ListCol struct {
	Heading   string
	Model     string
	Format    string
	Width     string
	IsImg     bool
	IsArray   bool
	Fieldname string
	IsBool    bool
	MaxChars  int
	IsIcon    bool
	CanEdit   bool
}

type ListForm struct {
	Title string
	Icon  string
	ID    int
	// KeyField    string
	Data        interface{}
	Cols        []*ListCol
	RowCB       func(string)
	CancelCB    func(dom.Event)
	NewRowCB    func(dom.Event)
	PrintCB     func(dom.Event)
	HasSetWidth bool
	Draggable   bool
	HasImages   bool
	MaxChars    int
}

// Init a new listform
func (f *ListForm) New(icon string, title string) *ListForm {
	f.Title = title
	f.Icon = icon
	f.PrintCB = nil
	f.MaxChars = 120
	// f.KeyField = "ID"
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

// Add print button
func (f *ListForm) PrintEvent(c func(dom.Event)) *ListForm {
	f.PrintCB = c
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
		Heading:  heading,
		Model:    model,
		MaxChars: f.MaxChars,
	}
	f.Cols = append(f.Cols, c)
	return f
}

// Add a colunm to the listform with format
func (f *ListForm) ColumnFormat(heading string, model string, format string) *ListForm {
	c := &ListCol{
		Heading:  heading,
		Model:    model,
		Format:   format,
		MaxChars: f.MaxChars,
	}
	f.Cols = append(f.Cols, c)
	return f
}

// Add a colunm to the listform in Date Format
func (f *ListForm) DateColumn(heading string, model string) *ListForm {
	c := &ListCol{
		Heading:  heading,
		Model:    model,
		Format:   "date",
		MaxChars: f.MaxChars,
	}
	f.Cols = append(f.Cols, c)
	return f
}

// Add a colunm to the listform in Email Format
func (f *ListForm) AvatarColumn(heading string, model string) *ListForm {
	c := &ListCol{
		Heading:  heading,
		Model:    model,
		Format:   "avatar",
		MaxChars: f.MaxChars,
	}
	f.Cols = append(f.Cols, c)
	return f
}

// Add a colunm to the listform in Email Format
func (f *ListForm) EmailAvatarColumn(heading string, model string) *ListForm {
	c := &ListCol{
		Heading:  heading,
		Model:    model,
		Format:   "email-avatar",
		MaxChars: f.MaxChars,
	}
	f.Cols = append(f.Cols, c)
	return f
}

// Add a colunm to the listform in Img Format
func (f *ListForm) ImgColumn(heading string, model string) *ListForm {
	c := &ListCol{
		Heading: heading,
		Model:   model,
		IsImg:   true,
	}
	f.Cols = append(f.Cols, c)
	f.HasImages = true
	return f
}

// Add a colunm that is edittable
func (f *ListForm) EditColumn(heading string, model string) *ListForm {
	c := &ListCol{
		Heading: heading,
		Model:   model,
		CanEdit: true,
	}
	f.Cols = append(f.Cols, c)
	f.HasImages = true
	return f
}

// Add a colunm to the listform in Img Format
func (f *ListForm) MultiImgColumn(heading string, model string, field string) *ListForm {
	c := &ListCol{
		Heading:   heading,
		Model:     model,
		IsImg:     true,
		IsArray:   true,
		Fieldname: field,
	}
	f.Cols = append(f.Cols, c)
	f.HasImages = true
	return f
}

// Add a colunm to the listform with a checkbox
func (f *ListForm) BoolColumn(heading string, model string) *ListForm {
	c := &ListCol{
		Heading: heading,
		Model:   model,
		IsBool:  true,
	}
	f.Cols = append(f.Cols, c)
	f.HasImages = true
	return f
}

// Add a colunm with a custom icon
func (f *ListForm) IconColumn(heading string, model string) *ListForm {
	c := &ListCol{
		Heading: heading,
		Model:   model,
		IsIcon:  true,
	}
	f.Cols = append(f.Cols, c)
	f.HasImages = true
	return f
}

func (f *ListForm) Render(name string, selector string, data interface{}) {
	f.Data = data
	if dom.GetWindow().Document().QuerySelector(selector) == nil {
		return
	}
	print("loading into selector", selector)
	renderTemplateT(f.generateTemplate(name, true), selector, f)
	f.decorate(selector)
}

func (f *ListForm) RenderNoContainer(name string, selector string, data interface{}) {
	f.Data = data
	if dom.GetWindow().Document().QuerySelector(selector) == nil {
		return
	}
	renderTemplateT(f.generateTemplate(name, false), selector, f)
	f.decorate(selector)
}

// Render the form using a custom template
func (f *ListForm) RenderCustom(name string, selector string, data interface{}) {

	f.Data = data
	renderTemplate(name, selector, data)
	f.decorate(selector)
}

func (f *ListForm) decorate(selector string) {

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

	if f.PrintCB != nil {
		if el := doc.QuerySelector(".data-print-btn"); el != nil {
			el.AddEventListener("click", false, f.PrintCB)
		}
	}

	// Handlers on the table itself
	sel := doc.QuerySelector(selector)
	if el := sel.QuerySelector(".data-table"); el != nil {

		if f.RowCB != nil {
			el.AddEventListener("click", false, func(evt dom.Event) {
				evt.PreventDefault()
				// Fix the issue where the user clicked on some clickable element inside the row
				// which adds an extra level which we dont usually need
				el := evt.Target()
				switch el.TagName() {
				case "INPUT":
					break
				case "TD":
					f.RowCB(el.ParentElement().GetAttribute("key"))
				case "TR":
					f.RowCB(el.GetAttribute("key"))
				}
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

func (f *ListForm) generateTemplate(name string, container bool) *temple.Template {

	isMobile := dom.GetWindow().InnerWidth() < 740
	//print("looking for template", name)
	tmpl, err := generatedTemplates.GetTemplate(name)
	if err != nil {
		// print("Generating template for", name)
		// Template doesnt exist, so create it

		src := ""
		doTitle := false

		if f.Title != "" || f.Icon != "" {
			doTitle = true
		}

		if doTitle {

			if container {

				src += `
<div class="data-container">
	<div class="row data-table-header">
    <h3 class="column column-90" id="legend">
      <i class="fa {{.Icon}} fa-lg" style="font-size: 3rem"></i> 
      {{.Title}}
    </h3>
`
			} else {
				src += `
	<div class="row data-table-header">
    <h3 class="column column-90" id="legend">
      <i class="fa {{.Icon}} fa-lg" style="font-size: 3rem"></i> 
      {{.Title}}
    </h3>
`
			}
			if f.NewRowCB != nil {
				src += `
    <div class="column col-center">
      <i class="data-add-btn fa fa-plus-circle fa-lg no-print"></i>    
    </div>    
`
			}

			if f.PrintCB != nil {
				src += `
    <div class="column col-center">
      <i class="data-print-btn fa fa-print fa-lg no-print"></i>    
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
    <tr class="data-row`

		if f.Draggable {
			src += ` draggable" draggable="true"`
		} else {
			src += `"`
		}

		src += ` key="{{.ID}}">`

		// for each column, add a column renderer

		for _, col := range f.Cols {
			width := ""
			if f.HasSetWidth {
				width = fmt.Sprintf(` width="%s"`, col.Width)
			}

			if col.IsImg {
				if col.IsArray { // MultiImgColunn
					src += fmt.Sprintf("<td %s %s>{{range $k,$v := .%s}}", width, col.Format, col.Model)

					src += fmt.Sprintf("{{if $v.%s}}<img name=%s-{{$k}}-{{.ID}} src={{$v.%s | safeURL}}>{{end}}",
						col.Fieldname, col.Model, col.Fieldname)

					src += "{{end}}</td>\n"
				} else { // ImgColumn
					src += fmt.Sprintf("<td %s>{{if .%s}}<img name=%s-{{.ID}} src={{.%s | safeURL}}>{{end}}</td>\n",
						width, col.Model, col.Model, col.Model)
				}
			} else if col.IsBool {
				src += fmt.Sprintf("<td %s %s>{{if .%s}}<i class=\"fa fa-check fa-lg\">{{end}}</td>\n",
					width, col.Format, col.Model)
			} else if col.IsIcon {
				src += fmt.Sprintf("<td %s %s><i class=\"{{.%s}}\"></td>\n",
					width, col.Format, col.Model)
			} else if col.Format == "date" {
				if isMobile {
					print("is a date col on mobile layout")
					src += fmt.Sprintf("<td %s>{{if .%s}}{{.%s.Format \"2 Jan\"}}{{end}}</td>\n",
						width, col.Model, col.Model)
				} else {
					src += fmt.Sprintf("<td %s>{{if .%s}}{{.%s.Format \"Mon, Jan 2 2006\"}}{{end}}</td>\n",
						width, col.Model, col.Model)
				}
			} else if col.Format == "avatar" {
				src += fmt.Sprintf("<td %s>{{if .%s}}<img src=\"{{.GetAvatar 64}}\">{{end}}</td>\n",
					width, col.Model)
			} else if col.Format == "email-avatar" {
				src += fmt.Sprintf("<td %s>{{if .%s}}<img src=\"{{.GetAvatar 40}}\"> {{.%s}}{{end}}</td>\n",
					width, col.Model, col.Model)
			} else if col.CanEdit {
				src += fmt.Sprintf("<td %s>{{if .%s}}<input type=\"text\" value=\"{{.%s}}\">{{end}}</td>",
					width, col.Model, col.Model)
			} else {
				if col.Format != "" {
					src += fmt.Sprintf("<td class=\"{{.%s}}\">{{if .%s}}{{.%s}}{{end}}</td>\n",
						col.Format, col.Model, col.Model)
				} else {
					src += fmt.Sprintf("<td %s>{{if .%s}}{{.%s}}{{end}}</td>\n",
						width, col.Model, col.Model)
				}
			}
		}

		if container {

			src += `      
    </tr>
{{end}}  
  <tbody>
  </tbody>
</table>
<div id="action-grid" class="no-print hidden"></div>
`
		} else {
			src += `      
    </tr>
{{end}}  
  <tbody>
  </tbody>
</table>
`
		}
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

// Add actions
func (f *ListForm) ActionGrid(template string, selector string, id interface{}, cb func(string)) {
	ActionGrid(template, selector, id, cb)
}

func (f *ListForm) OldActionGrid(template string, selector string, id interface{}, cb func(string)) {

	// print("add action grid")
	w := dom.GetWindow()
	doc := w.Document()

	el := doc.QuerySelector(selector)
	if el == nil {
		print("Could not find selector", selector)
		return
	}
	el.Class().Remove("hidden")

	renderTemplate(template, selector, id)
	for _, ai := range doc.QuerySelectorAll(".action__item") {
		url := ai.(*dom.HTMLDivElement).GetAttribute("url")
		if url != "" {
			ai.AddEventListener("click", false, func(evt dom.Event) {
				if evt.Target().TagName() == "INPUT" {
					evt.PreventDefault()
					print("clicked on input field")
				} else {
					print("cliked not on input field")
					url := evt.CurrentTarget().GetAttribute("url")
					cb(url)
				}
			})
		}
	}
}
