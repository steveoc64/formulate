package formulate

import (
	"fmt"
	"html/template"
	"reflect"
	"strconv"
	"time"
	"unsafe"

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
	Span      int
	Label     string
	Type      string
	Model     string
	Value     string
	Checked   bool
	Focusme   bool
	Readonly  bool
	Extras    template.CSS
	Class     string
	Step      string
	IsFloat   bool
	Decimals  int
	Options   []*EditOption
	Swapper   *Swapper
	Selected  int
	Group     []SelectGroup
	CodeBlock bool
	BigText   bool
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
	PrintCB  func(dom.Event)
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

func (p *Panel) Row(s int) *EditRow {
	r := EditRow{
		Span: s,
	}
	p.Rows = append(p.Rows, &r)
	return &r
}

func (p *Panel) AddRow(s int) *EditRow {
	return p.Row(s)
}

func (p *Panel) Paint(data interface{}) {

	w := dom.GetWindow()
	doc := w.Document()

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
			for _, r := range p.Rows {
				for _, f := range r.Fields {
					print("render swapper field", f.Model)
					dataField := reflect.Indirect(ptrVal).FieldByName(f.Model)
					print("Field Type", f.Type)
					switch dataField.Kind() {
					case reflect.Float64:
						f.Value = fmt.Sprintf("%.2f", dataField.Float())
					case reflect.Int:
						f.Value = fmt.Sprintf("%d", dataField.Int())
					case reflect.Ptr:
						// print("odd case of a swapper field being a ptr", f.Model, f.Type)
						switch f.Type {
						case "date":
							f.Value = ""
							ptr := unsafe.Pointer(dataField.Pointer())
							if ptr != nil {
								t := *(*time.Time)(ptr)
								f.Value = t.Format(rfc3339DateLayout)
							}
						case "number":
							f.Value = ""
							ptr := unsafe.Pointer(dataField.Pointer())
							if ptr != nil {
								if f.IsFloat {
									v := *(*float64)(ptr)
									f.Value = fmt.Sprintf("%f", v)
								} else {
									v := *(*int)(ptr)
									f.Value = fmt.Sprintf("%d", v)

								}
							}
						}
					case reflect.String:
						f.Value = dataField.String()
					default:
						f.Value = dataField.String()
					}
					print("Field", f.Type, f.Model, f.Value)
					switch f.Type {
					case "text", "number":
						el := doc.QuerySelector(fmt.Sprintf("[name=%s]", f.Model)).(*dom.HTMLInputElement)
						el.Value = f.Value
					case "textarea":
						el := doc.QuerySelector(fmt.Sprintf("[name=%s]", f.Model)).(*dom.HTMLTextAreaElement)
						el.Value = f.Value
					}
				}
			}
		}
	}

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

// Associate a print event with the editform
func (f *EditForm) PrintEvent(c func(dom.Event)) *EditForm {
	f.PrintCB = c
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
		Span:     span,
		Label:    label,
		Type:     "text",
		Focusme:  false,
		Model:    model,
		Readonly: false,
	}
	r.Fields = append(r.Fields, f)
	return r
}

// Add a Text input
func (r *EditRow) AddDisplay(span int, label string, model string) *EditRow {
	f := &EditField{
		Span:     span,
		Label:    label,
		Type:     "text",
		Focusme:  false,
		Model:    model,
		Readonly: true,
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

// Add a Floating Point Number input
func (r *EditRow) AddDecimal(span int, label string, model string, decimals int, step string) *EditRow {
	f := &EditField{
		Span:     span,
		Label:    label,
		Type:     "number",
		Focusme:  false,
		Step:     step,
		Model:    model,
		IsFloat:  true,
		Decimals: decimals,
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

// Add a Checkbok input
func (r *EditRow) AddCheck(span int, label string, model string) *EditRow {

	fld := &EditField{
		Span:    span,
		Label:   label,
		Type:    "checkbox",
		Focusme: false,
		Model:   model,
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

// Add a Textarea in bigtext mode
func (r *EditRow) AddBigTextarea(span int, label string, model string) *EditRow {
	f := &EditField{
		Span:    span,
		Label:   label,
		Type:    "textarea",
		Focusme: false,
		Model:   model,
		BigText: true,
	}
	r.Fields = append(r.Fields, f)
	return r
}

// Add a Textarea in readonly mode
func (r *EditRow) AddDisplayArea(span int, label string, model string) *EditRow {
	f := &EditField{
		Span:     span,
		Label:    label,
		Type:     "textarea",
		Focusme:  false,
		Model:    model,
		Readonly: true,
	}
	r.Fields = append(r.Fields, f)
	return r
}

// Add a Textarea in codeblock mode
func (r *EditRow) AddCodeBlock(span int, label string, model string) *EditRow {
	f := &EditField{
		Span:      span,
		Label:     label,
		Type:      "textarea",
		Focusme:   false,
		Model:     model,
		Readonly:  true,
		CodeBlock: true,
	}
	r.Fields = append(r.Fields, f)
	return r
}

// Add a Button
func (r *EditRow) AddButton(span int, label string, model string) *EditRow {
	f := &EditField{
		Span:      span,
		Label:     label,
		Type:      "button",
		Focusme:   false,
		Model:     model,
		Readonly:  true,
		CodeBlock: true,
	}
	r.Fields = append(r.Fields, f)
	return r
}

// Add a Div
func (r *EditRow) AddCustom(span int, label string, model string, class string) *EditRow {
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
							case reflect.Bool:
								field.Checked = dataField.Bool()
							default:
								// print(field.Model + " of type " + dataField.Kind().String())
								field.Value = dataField.String()
							}
						}
					} else { // field has no model - it could be a swapper
						switch field.Type {
						case "swapper":
							// Swapper has a slice of panels, which has a slice of rows, with a slice of fields
							for _, p := range field.Swapper.Panels {
								for _, r := range p.Rows {
									for _, f := range r.Fields {
										// print("render swapper field", f.Model)
										dataField := reflect.Indirect(ptrVal).FieldByName(f.Model)
										switch dataField.Kind() {
										case reflect.Float64:
											f.Value = fmt.Sprintf("%.2f", dataField.Float())
										case reflect.Int:
											f.Value = fmt.Sprintf("%d", dataField.Int())
										case reflect.Ptr:
											// print("odd case of a swapper field being a ptr", f.Model, f.Type)
											switch f.Type {
											case "date":
												f.Value = ""
												ptr := unsafe.Pointer(dataField.Pointer())
												if ptr != nil {
													t := *(*time.Time)(ptr)
													f.Value = t.Format(rfc3339DateLayout)
												}
											case "number":
												f.Value = ""
												ptr := unsafe.Pointer(dataField.Pointer())
												if ptr != nil {
													if f.IsFloat {
														v := *(*float64)(ptr)
														f.Value = fmt.Sprintf("%f", v)
													} else {
														v := *(*int)(ptr)
														f.Value = fmt.Sprintf("%d", v)

													}
												}
											}
										case reflect.String:
											f.Value = dataField.String()
										default:
											f.Value = dataField.String()
										}
									}
								}
							}
						}
					}
				}
			}
		}

	}

	renderTemplate(template, selector, f)

	// if f.CancelCB == nil {
	// 	print("Error - No cancel callback")
	// 	return
	// }
	// if f.SaveCB == nil {
	// 	print("Error - No save callback")
	// 	return
	// }

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
	if f.SaveCB != nil {
		if el := doc.QuerySelector(".md-save"); el != nil {
			el.AddEventListener("click", false, f.SaveCB)
		}
	}

	// plug in the print callback
	if f.PrintCB != nil {
		if el := doc.QuerySelector(".data-print-btn"); el != nil {
			el.AddEventListener("click", false, f.PrintCB)
		}
	}
}

// Add actions
func (f *EditForm) ActionGrid(template string, selector string, id interface{}, cb func(string)) {

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

// Programmatically reset the Form title
func (f *EditForm) SetTitle(title string) {
	w := dom.GetWindow()
	doc := w.Document()
	el := doc.QuerySelector("#titletext")
	// print("setting element", el, " was =", el.InnerHTML())
	el.SetInnerHTML(title)
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

			// If its a display only field, or a custom div
			// then dont bother binding it == much speed ++ safety
			if field.Readonly {
				continue
			}
			if field.Type == "div" {
				continue
			}

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
				// print("here with field", field)
				// print("datafield", dataField)
				// print("idx", idx)
				// print("opts key", field.Options[idx])
				setFromInt(dataField, field.Options[idx].Key)
			case "groupselect":
				idx := el.(*dom.HTMLSelectElement).SelectedIndex
				setFromInt(dataField, idx)
			case "checkbox":
				//print("binding into", dataField)
				print("with checked", el.(*dom.HTMLInputElement).Checked)
				print("with value", el.(*dom.HTMLInputElement).Value)
				setFromBool(dataField, el.(*dom.HTMLInputElement).Checked)
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
				if field.IsFloat {
					v, ferr := strconv.ParseFloat(ie.Value, 64)
					if ferr != nil {
						print("strconv.ParseFloat err ", ferr.Error())
					}
					setFromFloat(dataField, v)
				} else {
					v, ferr := strconv.Atoi(ie.Value)
					if ferr != nil {
						print("strconv.Atoi err ", ferr.Error())
					}
					setFromInt(dataField, v)
				}
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
							case "checkbox":
								setFromString(dataField, el.(*dom.HTMLInputElement).Value)
							case "radio":
								els := doc.QuerySelectorAll(name)
								for _, rel := range els {
									print("having a look at rel", rel)
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
								if field.IsFloat {
									if ie.Value == "" {
										setFromFloat(dataField, 0.0)
									} else {
										v, ferr := strconv.ParseFloat(ie.Value, 64)
										if ferr != nil {
											print("strconv.ParseFloat err ", ferr.Error(), "field=", f, "val=", ie.Value)
										}
										setFromFloat(dataField, v)
									}
								} else {
									if ie.Value == "" {
										setFromInt(dataField, 0)
									} else {
										v, ferr := strconv.Atoi(ie.Value)
										if ferr != nil {
											print("strconv.Atoi err ", ferr.Error(), "val =", ie.Value, "model =", f.Model, "field =", f)
										}
										setFromInt(dataField, v)
									}
								}
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

func setFromBool(target reflect.Value, v bool) {

	k := target.Kind()
	switch k {
	case reflect.Bool:
		target.SetBool(v)
	case reflect.Int:
		i := int64(0)
		if v {
			i = int64(1)
		}
		target.SetInt(i)
	case reflect.Float64:
		// print("conversion of string to float")
		i := 0.0
		if v {
			i = 1.0
		}
		target.SetFloat(i)
	case reflect.String:
		str := "true"
		if !v {
			str = "false"
		}
		target.SetString(str)
	default:
		print("conversion of bool to unknown type", k.String())
		str := "true"
		if !v {
			str = "false"
		}
		target.SetString(str)
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
		// print("conversion of int to ptr")
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

func setFromFloat(target reflect.Value, v float64) {

	k := target.Kind()
	switch k {
	case reflect.Bool:
		// print("conversion of int to bool")
		target.SetBool(v != 0.0)
	case reflect.Int:
		// print("conversion of int to int")
		target.SetInt(int64(v))
	case reflect.Float64:
		// print("conversion of int to float")
		target.SetFloat(v)
	case reflect.Ptr:
		// print("conversion of float to ptr")
		target.Set(reflect.ValueOf(&v))
	case reflect.String:
		// print("conversion of int to string")
		target.SetString(fmt.Sprintf("%f", v))
	default:
		print("conversion of float to unknown type", k.String())
		target.SetFloat(v)
	}
}
