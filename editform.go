package formulate

import (
	"fmt"
	"html/template"
	"reflect"
	"strconv"
	"time"

	"honnef.co/go/js/dom"
)

type SelectOption struct {
	ID   int
	Name string
}

type SelectGroup struct {
	Title   string
	Options []SelectOption
}

type EditOption struct {
	Key      int
	Display  string
	Selected bool
}

type EditField struct {
	Span     int
	Label    string
	Type     string
	Model    string
	Value    string
	Focusme  bool
	Extras   template.CSS
	Class    string
	Step     string
	Options  []*EditOption
	Swapper  *Swapper
	Selected int
	Group    []SelectGroup
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

type Swapper struct {
	Name     string
	Selected int
	Panels   []*Panel
}

func (s *Swapper) AddPanel(panelName string) *Panel {
	p := Panel{
		Name: panelName,
	}
	s.Panels = append(s.Panels, &p)
	return &p
}

func (s *Swapper) Select(idx int) {
	w := dom.GetWindow()
	doc := w.Document()

	// Show or unshow all panels by name
	for i, p := range s.Panels {
		el := doc.QuerySelector(fmt.Sprintf("[name=%s-%s]", s.Name, p.Name)).(*dom.HTMLDivElement)
		if el != nil {
			cl := el.Class()
			if i == idx {
				s.Selected = i
				cl.Add("swapper-show")
			} else {
				cl.Remove("swapper-show")
			}
		}
	}
}

func (s *Swapper) SelectByName(name string) {
	w := dom.GetWindow()
	doc := w.Document()

	// Show or unshow all panels by name
	for i, p := range s.Panels {
		el := doc.QuerySelector(fmt.Sprintf("[name=%s-%s]", s.Name, p.Name)).(*dom.HTMLDivElement)
		if el != nil {
			cl := el.Class()
			if p.Name == name {
				s.Selected = i
				cl.Add("swapper-show")
			} else {
				cl.Remove("swapper-show")
			}
		}
	}
}

type Panel struct {
	Name string
	Div  *dom.HTMLDivElement
	Rows []*EditRow
}

func (p *Panel) AddRow(s int) *EditRow {
	r := EditRow{
		Span: s,
	}
	p.Rows = append(p.Rows, &r)
	return &r
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
		Span:    span,
		Label:   label,
		Type:    t,
		Model:   model,
		Focusme: false,
		Extras:  template.CSS(extras),
	}
	r.Fields = append(r.Fields, f)
	// print("add to row", span, label, t, model, extras)
	return r
}

// Add a Text input
func (r *EditRow) AddInput(span int, label string, model string) *EditRow {
	f := &EditField{
		Span:    span,
		Label:   label,
		Type:    "text",
		Focusme: false,
		Model:   model,
	}
	r.Fields = append(r.Fields, f)
	return r
}

// Add a Date input
func (r *EditRow) AddDate(span int, label string, model string) *EditRow {
	f := &EditField{
		Span:    span,
		Label:   label,
		Type:    "date",
		Focusme: false,
		Model:   model,
	}
	r.Fields = append(r.Fields, f)
	return r
}

// Add a Number input
func (r *EditRow) AddNumber(span int, label string, model string, step string) *EditRow {
	f := &EditField{
		Span:    span,
		Label:   label,
		Type:    "number",
		Focusme: false,
		Step:    step,
		Model:   model,
	}
	r.Fields = append(r.Fields, f)
	return r
}

// Add a Radio input
func (r *EditRow) AddRadio(span int, label string, model string,
	options interface{}, key string, value string, selectedKey int) *EditRow {
	fld := &EditField{
		Span:    span,
		Label:   label,
		Type:    "radio",
		Focusme: false,
		Model:   model,
	}

	// Now loop through the options and append to the options array
	ptrVal := reflect.ValueOf(options)

	olen := ptrVal.Len()

	for i := 0; i < olen; i++ {
		o := ptrVal.Index(i)
		okey := int(o.FieldByName(key).Int())
		oval := o.FieldByName(value).String()
		fld.Options = append(fld.Options, &EditOption{
			Key:      okey,
			Display:  oval,
			Selected: okey == selectedKey,
		})
	}

	r.Fields = append(r.Fields, fld)
	return r
}

// Add a panel swapper
func (r *EditRow) AddSwapper(span int, label string, swapper *Swapper) *EditRow {

	fld := &EditField{
		Span:    span,
		Label:   label,
		Type:    "swapper",
		Swapper: swapper,
	}

	r.Fields = append(r.Fields, fld)
	return r
}

// Add a Select element
func (r *EditRow) AddSelect(span int, label string, model string,
	options interface{}, key string, value string,
	min int, selectedKey int) *EditRow {

	fld := &EditField{
		Span:     span,
		Label:    label,
		Type:     "select",
		Focusme:  false,
		Model:    model,
		Selected: selectedKey,
	}

	// If min = 0, then we start with a blank option for "nothing selected"
	if min == 0 {
		fld.Options = append(fld.Options, &EditOption{
			Key:     0,
			Display: "",
		})
	}

	// Now loop through the options and append to the options array
	ptrVal := reflect.ValueOf(options)

	olen := ptrVal.Len()

	for i := 0; i < olen; i++ {
		o := ptrVal.Index(i)
		okey := int(o.FieldByName(key).Int())
		oval := o.FieldByName(value).String()
		fld.Options = append(fld.Options, &EditOption{
			Key:      okey,
			Display:  oval,
			Selected: okey == selectedKey,
		})
	}

	r.Fields = append(r.Fields, fld)
	return r
}

// Add a GroupedSelect element
func (r *EditRow) AddGroupedSelect(span int, label string, model string,
	group []SelectGroup, selectedKey int) *EditRow {

	fld := &EditField{
		Span:     span,
		Label:    label,
		Type:     "groupselect",
		Focusme:  false,
		Model:    model,
		Group:    group,
		Selected: selectedKey,
	}
	r.Fields = append(r.Fields, fld)
	return r
}

// Add a Textarea
func (r *EditRow) AddTextarea(span int, label string, model string) *EditRow {
	f := &EditField{
		Span:    span,
		Label:   label,
		Type:    "textarea",
		Focusme: false,
		Model:   model,
	}
	r.Fields = append(r.Fields, f)
	return r
}

// Add a Div
func (r *EditRow) AddDiv(span int, label string, model string, class string) *EditRow {
	f := &EditField{
		Span:    span,
		Label:   label,
		Type:    "div",
		Focusme: false,
		Model:   model,
		Class:   class,
	}
	r.Fields = append(r.Fields, f)
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
					if field.Model != "" {
						switch field.Type {
						case "div":
							// is just a placeholder div field, so dont bind it
						default:
							dataField := reflect.Indirect(ptrVal).FieldByName(field.Model)
							switch dataField.Kind() {
							case reflect.Float64:
								// print(field.Model + " of type " + dataField.Kind().String())
								field.Value = fmt.Sprintf("%.2f", dataField.Float())
							case reflect.Int:
								// print(field.Model + " of type " + dataField.Kind().String())
								field.Value = fmt.Sprintf("%d", dataField.Int())
							case reflect.Ptr:
								// print(field.Model + " of type " + dataField.Kind().String())
								field.Value = dataField.String()
							case reflect.String:
								field.Value = dataField.String()
							default:
								// print(field.Model + " of type " + dataField.Kind().String())
								field.Value = dataField.String()
							}
						}
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

			name := `[name="` + field.Model + `"]`
			el := doc.QuerySelector(name)
			dataField := reflect.Indirect(ptrVal).FieldByName(field.Model)

			switch field.Type {
			case "text":
				setFromString(dataField, el.(*dom.HTMLInputElement).Value)
			case "textarea":
				setFromString(dataField, el.(*dom.HTMLTextAreaElement).Value)
			case "select":
				idx := el.(*dom.HTMLSelectElement).SelectedIndex
				setFromInt(dataField, field.Options[idx].Key)
			case "groupselect":
				idx := el.(*dom.HTMLSelectElement).SelectedIndex
				setFromInt(dataField, idx)
			case "radio":
				els := doc.QuerySelectorAll(name)
				for _, rel := range els {
					ie := rel.(*dom.HTMLInputElement)
					if ie.Checked {
						v, _ := strconv.Atoi(ie.Value)
						setFromInt(dataField, v)
						break
					}
				}
			case "number":
				ie := el.(*dom.HTMLInputElement)
				v, _ := strconv.Atoi(ie.Value)
				setFromInt(dataField, v)
			case "date":
				ie := el.(*dom.HTMLInputElement)
				setFromDate(dataField, ie.Value)
				print("TODO - bind from date field", ie.Value)
			case "div":
				// is just a placeholder, dont bind it
			case "swapper":
				// Swapper has a slice of panels
				for _, p := range field.Swapper.Panels {
					// Panel has a slice of rows
					for _, r := range p.Rows {
						// Row has a slice of fields
						for _, f := range r.Fields {
							name := `[name="` + f.Model + `"]`
							el := doc.QuerySelector(name)
							dataField := reflect.Indirect(ptrVal).FieldByName(f.Model)
							switch f.Type {
							case "text":
								setFromString(dataField, el.(*dom.HTMLInputElement).Value)
							case "textarea":
								setFromString(dataField, el.(*dom.HTMLTextAreaElement).Value)
							case "select":
								idx := el.(*dom.HTMLSelectElement).SelectedIndex
								setFromInt(dataField, f.Options[idx].Key)
							case "radio":
								els := doc.QuerySelectorAll(name)
								for _, rel := range els {
									ie := rel.(*dom.HTMLInputElement)
									if ie.Checked {
										// print("swapper radio", name, "value =", ie.Value)
										v, _ := strconv.Atoi(ie.Value)
										setFromInt(dataField, v)
										break
									}
								}
							case "number":
								ie := el.(*dom.HTMLInputElement)
								v, _ := strconv.Atoi(ie.Value)
								setFromInt(dataField, v)
							case "date":
								ie := el.(*dom.HTMLInputElement)
								setFromDate(dataField, ie.Value)
								// print("TODO - bind swapper from date field", ie.Value)
							}

						}
					}

				}

			default:
				print("TODO - bind from ", field.Type)
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
		// print("conversion of string to int")
		i, _ := strconv.ParseInt(str, 0, 64)
		target.SetInt(i)
	case reflect.Float64:
		// print("conversion of string to float")
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

const (
	rfc3339DateLayout          = "2006-01-02"
	rfc3339DatetimeLocalLayout = "2006-01-02T15:04:05.999999999"
)

func setFromDate(target reflect.Value, str string) {

	thedate, _ := time.Parse(rfc3339DateLayout, str)
	// print("Parse", str, "as", thedate.String())

	k := target.Kind()
	switch k {
	case reflect.Ptr:
		// print("target should be a *time.Time")
		target.Set(reflect.ValueOf(&thedate))
	case reflect.Struct:
		// print("target should be a time.Time")
		target.Set(reflect.ValueOf(thedate))
	default:
		print("conversion of date to unknown type", k.String())
		// target.SetString(str)
	}
}

func setFromInt(target reflect.Value, v int) {

	k := target.Kind()
	switch k {
	case reflect.Bool:
		// print("conversion of int to bool")
		target.SetBool(v != 0)
	case reflect.Int:
		// print("conversion of int to int")
		target.SetInt(int64(v))
	case reflect.Float64:
		// print("conversion of int to float")
		target.SetFloat(float64(int64(v)))
	case reflect.Ptr:
		print("conversion of int to ptr")
		target.Set(reflect.ValueOf(&v))
		// print("lets just try setting the *int directly from the value ", v)
		// target.SetInt(int64(v))
		// print("that gets us", target.String())
	case reflect.String:
		// print("conversion of int to string")
		target.SetString(fmt.Sprintf("%d", v))
	default:
		print("conversion of int to unknown type", k.String())
		target.SetInt(int64(v))
	}
}
