package formulate

import (
	"github.com/go-humble/temple/temple"

	"honnef.co/go/js/dom"
)

type TreeForm struct {
	Title       string
	Icon        string
	ID          int
	Data        interface{}
	Cols        []*ListCol
	RowCB       func(string)
	CancelCB    func(dom.Event)
	NewRowCB    func(dom.Event)
	PrintCB     func(dom.Event)
	HasSetWidth bool
}

// Init a new listform
func (f *TreeForm) New(icon string, title string) *TreeForm {
	f.Title = title
	f.Icon = icon
	f.PrintCB = nil
	return f
}

func (f *TreeForm) SetWidths(w []string) {
	f.HasSetWidth = true
	for i, v := range w {
		if i < len(f.Cols) {
			f.Cols[i].Width = v
		}
	}
}

// Associate a cancel event with the listform
func (f *TreeForm) CancelEvent(c func(dom.Event)) *TreeForm {
	f.CancelCB = c
	return f
}

// Add new row
func (f *TreeForm) NewRowEvent(c func(dom.Event)) *TreeForm {
	f.NewRowCB = c
	return f
}

// Add print button
func (f *TreeForm) PrintEvent(c func(dom.Event)) *TreeForm {
	f.PrintCB = c
	return f
}

// Associate a Row Click event with the listform
func (f *TreeForm) RowEvent(c func(string)) *TreeForm {
	f.RowCB = c
	return f
}

// Render the form using a template that we generate on the fly
func (f *TreeForm) Render(name string, selector string, data interface{}) {

	f.Data = data
	renderTemplateT(f.generateTemplate(name), selector, f)
	f.decorate(selector)
}

// Render the form using a custom template
func (f *TreeForm) RenderCustom(name string, selector string, data interface{}) {

	f.Data = data
	renderTemplate(name, selector, data)
	f.decorate(selector)
}

func (f *TreeForm) decorate(selector string) {

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

func (f *TreeForm) generateTemplate(name string) *temple.Template {

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
<div class="data-container css-treeview">
	<div class="row data-table-header">
    <h3 class="column column-90" id="legend">
      <i class="fa {{.Icon}} fa-lg" style="font-size: 3rem"></i> 
      {{.Title}}
    </h3>
`
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

		src += `
	    <ul>
        <li><input type="checkbox" id="item-0" /><label for="item-0">This Folder is Closed By Default</label>
            <ul>
                <li><input type="checkbox" id="item-0-0" /><label for="item-0-0">Ooops! A Nested Folder</label>
                    <ul>
                        <li><input type="checkbox" id="item-0-0-0" /><label for="item-0-0-0">Look Ma - No Hands!</label>
                            <ul>
                                <li><a href="./">First Nested Item</a></li>
                                <li><a href="./">Second Nested Item</a></li>
                                <li><a href="./">Third Nested Item</a></li>
                                <li><a href="./">Fourth Nested Item</a></li>
                            </ul>
                        </li>
                        <li><a href="./">Item 1</a></li>
                        <li><a href="./">Item 2</a></li>
                        <li><a href="./">Item 3</a></li>
                    </ul>
                </li>
                <li><input type="checkbox" id="item-0-1" /><label for="item-0-1">Yet Another One</label>
                    <ul>
                        <li><a href="./">item</a></li>
                        <li><a href="./">item</a></li>
                        <li><a href="./">item</a></li>
                        <li><a href="./">item</a></li>
                        <li><a href="./">item</a></li>
                        <li><a href="./">item</a></li>
                        <li><a href="./">item</a></li>
                        <li><a href="./">item</a></li>
                    </ul>
                </li>
                <li><input type="checkbox" id="item-0-2" disabled="disabled" /><label for="item-0-2">Disabled Nested Items</label>
                    <ul>
                        <li><a href="./">item</a></li>
                        <li><a href="./">item</a></li>
                        <li><a href="./">item</a></li>
                        <li><a href="./">item</a></li>
                        <li><a href="./">item</a></li>
                        <li><a href="./">item</a></li>
                        <li><a href="./">item</a></li>
                        <li><a href="./">item</a></li>
                    </ul>
                </li>
                <li><a href="./">item</a></li>
                <li><a href="./">item</a></li>
                <li><a href="./">item</a></li>
                <li><a href="./">item</a></li>
        </ul>
</li>
<li><input type="checkbox" id="item-1" checked="checked" /><label for="item-1">This One is Open by Default...</label>
        <ul>
            <li><input type="checkbox" id="item-1-0" /><label for="item-1-0">And Contains More Nested Items...</label>
                <ul>
                    <li><a href="./">Look Ma - No Hands</a></li>
                    <li><a href="./">Another Item</a></li>
                    <li><a href="./">And Yet Another</a></li>
                </ul>
            </li>
            <li><a href="./">Lorem</a></li>
            <li><a href="./">Ipsum</a></li>
            <li><a href="./">Dolor</a></li>
            <li><a href="./">Sit Amet</a></li>
        </ul>
</li>
<li><input type="checkbox" id="item-2" /><label for="item-2">Can You Believe...</label>
        <ul>
                <li><input type="checkbox" id="item-2-0" /><label for="item-2-0">That This Treeview...</label>
                    <ul>
                        <li><input type="checkbox" id="item-2-2-0" /><label for="item-2-2-0">Does Not Use Any JavaScript...</label>
                            <ul>
                                <li><a href="./">But Relies Only</a></li>
                                <li><a href="./">On the Power</a></li>
                                <li><a href="./">Of CSS3</a></li>
                            </ul>
                        </li>
                        <li><a href="./">Item 1</a></li>
                        <li><a href="./">Item 2</a></li>
                        <li><a href="./">Item 3</a></li>
                    </ul>
                </li>
                <li><input type="checkbox" id="item-2-1" /><label for="item-2-1">This is a Folder With...</label>
                    <ul>
                        <li><a href="./">Some Nested Items...</a></li>
                        <li><a href="./">Some Nested Items...</a></li>
                        <li><a href="./">Some Nested Items...</a></li>
                        <li><a href="./">Some Nested Items...</a></li>
                        <li><a href="./">Some Nested Items...</a></li>
                    </ul>
                </li>
                <li><input type="checkbox" id="item-2-2" disabled="disabled" /><label for="item-2-2">Disabled Nested Items</label>
                    <ul>
                        <li><a href="./">item</a></li>
                        <li><a href="./">item</a></li>
                        <li><a href="./">item</a></li>
                        <li><a href="./">item</a></li>
                        <li><a href="./">item</a></li>
                        <li><a href="./">item</a></li>
                        <li><a href="./">item</a></li>
                        <li><a href="./">item</a></li>
                    </ul>
                </li>
            </ul>
        </li>
    </ul>	
`

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
