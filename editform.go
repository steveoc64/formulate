package formulate

import (
	"fmt"
	"reflect"
	"strconv"

	"honnef.co/go/js/dom"
)

type EditOption struct {
	Key      int
	Display  string
	Selected bool
}

type EditField struct {
	Span    int
	Label   string
	Type    string
	Model   string
	Value   string
	Extras  string
	Options []*EditOption
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
	DeleteCB func(dom.Event)
	SaveCB   func(dom.Event)
}

// Get the editfield of the given name
func (f *EditForm) GetField(name string) *EditField {

	for _, row := range f.Rows {
		for _, field := range row.Fields {
			if field.Model == name {
				return field
			}
		}
	}
	return nil
}

// Apply options to a select field
func (f *EditForm) SetSelectOptions(name string,
	options interface{},
	key string,
	value string,
	min int,
	selectedKey int) {

	fld := f.GetField(name)
	if fld == nil {
		print("Cannot find field by name", name)
		return
	}

	// If min = 0, then we start with a blank option for "nothing selected"
	if min == 0 {
		fld.Options = append(fld.Options, &EditOption{
			Key:     0,
			Display: "",
		})
	}

	// Now loop through the options and append to the options array
	ptrType := reflect.TypeOf(options)
	// print("options kind =", ptrType.Kind().String())
	if ptrType.Kind() != reflect.Slice {
		return
	}
	typ := ptrType.Elem()
	// print("element kind =", typ.Kind().String())
	if typ.Kind() != reflect.Struct {
		return
	}
	ptrVal := reflect.ValueOf(options)
	if ptrVal.IsNil() {
		// print("contents of options is null")
		return
	}

	olen := ptrVal.Len()

	for i := 0; i < olen; i++ {
		o := ptrVal.Index(i)
		okey := int(o.FieldByName(key).Int())
		oval := o.FieldByName(value).String()
		// print("key/val", i, okey, oval)
		fld.Options = append(fld.Options, &EditOption{
			Key:      okey,
			Display:  oval,
			Selected: okey == selectedKey,
		})
	}

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

// Associate a delete event with the editform
func (f *EditForm) DeleteEvent(c func(dom.Event)) *EditForm {
	f.DeleteCB = c
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
func (r *EditRow) AddField(f EditField) *EditRow {
	if f.Span == 0 {
		f.Span = 1
	}
	r.Fields = append(r.Fields, &f)
	// print("=", r.Fields)
	return r
}

// Add a field with params
func (r *EditRow) Add(span int, label string, t string, model string, extras string) *EditRow {
	f := &EditField{
		Span:   span,
		Label:  label,
		Type:   t,
		Model:  model,
		Extras: extras,
	}
	r.Fields = append(r.Fields, f)
	// print("add to row", span, label, t, model, extras)
	return r
}

// Render the form
func (f *EditForm) Render(template string, selector string, data interface{}) {

	w := dom.GetWindow()
	doc := w.Document()

	// Tricky part here - if data is passed in, then
	// load the field values from the data

	if data != nil {
		// Make sure the type of v is a pointer to a struct.
		doit := true
		ptrType := reflect.TypeOf(data)
		if ptrType.Kind() != reflect.Ptr {
			doit = false
		}
		typ := ptrType.Elem()
		if typ.Kind() != reflect.Struct {
			doit = false
		}
		ptrVal := reflect.ValueOf(data)
		if ptrVal.IsNil() {
			doit = false
		}

		if doit {
			for _, row := range f.Rows {
				for _, field := range row.Fields {
					switch field.Type {
					case "div":
						// is just a placeholder div field, so dont bind it
					default:
						dataField := reflect.Indirect(ptrVal).FieldByName(field.Model)
						field.Value = dataField.String()
					}
				}
			}
		}

	}

	renderTemplate(template, selector, f)

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
	if f.CancelCB != nil {

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
	}

	if f.DeleteCB != nil {
		if el := doc.QuerySelector(".md-confirm-del"); el != nil {
			el.AddEventListener("click", false, f.DeleteCB)
		}

		if el := doc.QuerySelector(".data-del-btn"); el != nil {
			el.AddEventListener("click", false, func(evt dom.Event) {
				doc.QuerySelector("#confirm-delete").Class().Add("md-show")
			})
		}

		if el := doc.QuerySelector(".md-close-del"); el != nil {
			el.AddEventListener("click", false, func(evt dom.Event) {
				doc.QuerySelector("#confirm-delete").Class().Remove("md-show")
			})
		}

		if el := doc.QuerySelector("#confirm-delete"); el != nil {
			el.AddEventListener("keyup", false, func(evt dom.Event) {
				if evt.(*dom.KeyboardEvent).KeyCode == 27 {
					evt.PreventDefault()
					doc.QuerySelector("#confirm-delete").Class().Remove("md-show")
				}
			})
		}

	}

	// plug in the save callback
	if el := doc.QuerySelector(".md-save"); el != nil {
		el.AddEventListener("click", false, f.SaveCB)
	}
}

// Add actions
func (f *EditForm) ActionGrid(template string, selector string, id int, cb func(string)) {

	// print("add action grid")
	w := dom.GetWindow()
	doc := w.Document()

	renderTemplate(template, selector, id)
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

// Read the DOM values of each field back into the data
func (f *EditForm) Bind(data interface{}) {
	// print("binding fields to data")
	w := dom.GetWindow()
	doc := w.Document()

	// Make sure the type of v is a pointer to a struct.
	ptrType := reflect.TypeOf(data)
	if ptrType.Kind() != reflect.Ptr {
		print("form: Bind expects a pointer to a struct, but got: %T", data)
		return
	}
	typ := ptrType.Elem()
	if typ.Kind() != reflect.Struct {
		print("form: Bind expects a pointer to a struct, but got: %T", data)
		return
	}
	ptrVal := reflect.ValueOf(data)
	if ptrVal.IsNil() {
		print("form: Argument to Bind was nil")
		return
	}

	for _, row := range f.Rows {
		for _, field := range row.Fields {

			el := doc.QuerySelector(`[name="` + field.Model + `"]`)
			dataField := reflect.Indirect(ptrVal).FieldByName(field.Model)

			switch field.Type {
			case "text":
				setFromString(dataField, el.(*dom.HTMLInputElement).Value)
			case "textarea":
				setFromString(dataField, el.(*dom.HTMLTextAreaElement).Value)
			case "select":
				idx := el.(*dom.HTMLSelectElement).SelectedIndex
				setFromInt(dataField, field.Options[idx].Key)
			case "div":
				// is just a placeholder, dont bind it
			}
		}
	}

}

func setFromString(target reflect.Value, str string) {

	k := target.Kind()
	switch k {
	case reflect.Bool:
		print("conversion of string to bool")
		switch str {
		case "false", "False", "no", "No", "":
			target.SetBool(false)
		default:
			target.SetBool(true)
		}
	case reflect.Int:
		print("conversion of string to int")
		i, _ := strconv.ParseInt(str, 0, 64)
		target.SetInt(i)
	case reflect.Float64:
		print("conversion of string to float")
		i, _ := strconv.ParseFloat(str, 64)
		target.SetFloat(i)
	case reflect.Ptr:
		print("conversion of string to ptr")
		target.SetString(str)
	case reflect.String:
		// print("conversion of string to string")
		target.SetString(str)
	default:
		print("conversion of string to unknown type", k.String())
		target.SetString(str)
	}
}

func setFromInt(target reflect.Value, v int) {

	k := target.Kind()
	i := int64(v)
	switch k {
	case reflect.Bool:
		print("conversion of int to bool")
		target.SetBool(v != 0)
	case reflect.Int:
		// print("conversion of int to int")
		target.SetInt(i)
	case reflect.Float64:
		print("conversion of int to float")
		target.SetFloat(float64(i))
	case reflect.Ptr:
		print("conversion of int to ptr")
		target.SetInt(i)
	case reflect.String:
		print("conversion of int to string")
		target.SetString(fmt.Sprintf("%d", v))
	default:
		print("conversion of int to unknown type", k.String())
		target.SetInt(i)
	}
}
