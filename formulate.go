package formulate

import (
	"errors"

	"github.com/go-humble/temple/temple"
	"honnef.co/go/js/dom"
)

var gt func(name string) (*temple.Template, error)

func Templates(g func(string) (*temple.Template, error)) {
	gt = g
}

type EditField struct {
	Span  int
	ID    string
	Label string
	Type  string
	Name  string
	Value string
}

type EditRow struct {
	Span   int
	Fields []*EditField
}

type EditForm struct {
	Title    string
	Icon     string
	ID       int
	Rows     []*EditRow
	CancelCB func(dom.Event)
	SaveCB   func(dom.Event)
}

// Init a new editform
func (f *EditForm) New(icon string, title string) *EditForm {
	f.Title = title
	f.Icon = icon
	return f
}

// Associate a cancel event with the editform
func (f *EditForm) CancelEvent(c func(dom.Event)) *EditForm {
	f.CancelCB = c
	return f
}

// Associate a save event with the editform
func (f *EditForm) SaveEvent(c func(dom.Event)) *EditForm {
	f.SaveCB = c
	return f
}

// Add a row to an edit form
func (f *EditForm) Row(s int) *EditRow {
	r := EditRow{
		Span: s,
	}
	f.Rows = append(f.Rows, &r)
	return &r
}

// Add a field to a row on an edit form
func (r *EditRow) Add(f EditField) *EditRow {
	if f.Span == 0 {
		f.Span = 1
	}
	r.Fields = append(r.Fields, &f)
	// print("=", r.Fields)
	return r
}

// Render the form
func (f *EditForm) Render(template string, selector string) {

	w := dom.GetWindow()
	doc := w.Document()
	f.RenderTemplate(template, selector, f)

	if f.CancelCB == nil {
		print("Error - No cancel callback")
		return
	}
	if f.SaveCB == nil {
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

	if el := doc.QuerySelector(".grid-form"); el != nil {
		el.AddEventListener("keyup", false, func(evt dom.Event) {
			if evt.(*dom.KeyboardEvent).KeyCode == 27 {
				evt.PreventDefault()
				el.AddEventListener("click", false, f.CancelCB)
			}
		})
	}

	// plug in the save callback

	if el := doc.QuerySelector(".md-save"); el != nil {
		el.AddEventListener("click", false, f.SaveCB)
	}
}

// Add actions
func (f *EditForm) ActionGrid(template string, selector string, id int, cb func(string)) {

	print("add action grid")
	w := dom.GetWindow()
	doc := w.Document()

	f.RenderTemplate(template, selector, id)
	for _, ai := range doc.QuerySelectorAll(".action__item") {
		url := ai.(*dom.HTMLDivElement).GetAttribute("url")
		if url != "" {
			ai.AddEventListener("click", false, func(evt dom.Event) {
				url := evt.CurrentTarget().GetAttribute("url")
				cb(url)
			})
		}
	}
}

// Load a template and attach it to the specified element in the doc
func (f *EditForm) RenderTemplate(template string, selector string, data interface{}) error {
	w := dom.GetWindow()
	doc := w.Document()

	t, err := gt(template)
	if t == nil {
		print("Failed to load template", template)
		return errors.New("Invalid template")
	}
	if err != nil {
		print(err.Error())
		return err
	}

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

// Read the DOM values of each field back into the data
func (f *EditForm) Bind() {
	print("binding fields to data")
	w := dom.GetWindow()
	doc := w.Document()

	for _, row := range f.Rows {
		for _, field := range row.Fields {
			print("binding ", field)

			el := doc.QuerySelector(`[name="` + field.Name + `"]`)
			print("element = ", el)
			switch field.Type {
			case "text", "textarea":
				v := el.Value
				print("value =", v)
			}
		}
	}
}
