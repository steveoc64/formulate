package formulate

import (
	"github.com/go-humble/temple/temple"

	"honnef.co/go/js/dom"
)

type TreeCategories interface {
	Len() int
	Get(int) TreeData
}

type TreeElements interface {
	Len() int
	Render(int)
}

type TreeData interface {
	String() string
	Categories() TreeCategories
	Elements() TreeElements
}

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
func (f *TreeForm) Render(name string, selector string, data ...TreeData) {

	f.Data = data
	print("passed in treedata", data)
	renderTemplateT(f.generateTemplate(name), selector, f)
	f.decorate(selector)
}

// Render the form using a custom template
func (f *TreeForm) RenderCustom(name string, selector string, data ...TreeData) {

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
	    <ul class="css-treeview">
        <li><input type="checkbox" id="item-0" />
        		<label for="item-0">Bracket B4/5</label>
            <ul>
                <li><input type="checkbox" id="item-0-0" />
                	<label for="item-0-0">Lower Crop Blade</label>
                	<ul>
		                <li><input type="checkbox" id="item-1-0-0" />
		                	<label for="item-1-0-0">Consumables</label>
		                	<ul>
				                <li>Sharpening Stone</li>
				                <li>Oil</li>
				                <li>Razor Blades</li>
		                	</ul>
		                </li>
		                <li>Part 1</li>
		                <li>Part 2</li>
		                <li>Part 3</li>
		                <li>Part 4</li>
		                <li>Part 5</li>
                	</ul>
                </li>

                </li>
                <li><input type="checkbox" id="item-0-1" />
                	<label for="item-0-1">Upper Crop Blade</label>
                	<ul>
		                <li><input type="checkbox" id="item-1-0-0" />
		                	<label for="item-1-0-0">Part 1</label>
		                </li>
		                <li><input type="checkbox" id="item-1-0-1" />
		                	<label for="item-1-0-1">Part 2</label>
		                </li>
		                <li><input type="checkbox" id="item-1-0-2" />
		                	<label for="item-1-0-2">Part 3</label>
		                </li>
		                <li><input type="checkbox" id="item-1-0-3" />
		                	<label for="item-1-0-3">Part 4</label>
		                </li>
		                <li><input type="checkbox" id="item-1-0-4" />
		                	<label for="item-1-0-4">Part 5</label>
		                </li>
                	</ul>
                </li>

                </li>
            </ul>
        </li>
        <li><input type="checkbox" id="item-1" />
        		<label for="item-1">Chord</label>
            <ul>
                <li><input type="checkbox" id="item-1-0" />
                	<label for="item-1-0">Down Dimple</label>
                	<ul>
		                <li><input type="checkbox" id="item-1-0-0" />
		                	<label for="item-1-0-0">Part 1</label>
		                </li>
		                <li><input type="checkbox" id="item-1-0-1" />
		                	<label for="item-1-0-1">Part 2</label>
		                </li>
		                <li><input type="checkbox" id="item-1-0-2" />
		                	<label for="item-1-0-2">Part 3</label>
		                </li>
		                <li><input type="checkbox" id="item-1-0-3" />
		                	<label for="item-1-0-3">Part 4</label>
		                </li>
		                <li><input type="checkbox" id="item-1-0-4" />
		                	<label for="item-1-0-4">Part 5</label>
		                </li>
                	</ul>
                </li>
                <li><input type="checkbox" id="item-1-1" />
                	<label for="item-1-1">Tie Down Slot</label>
                	<ul>
		                <li><input type="checkbox" id="item-1-0-0" />
		                	<label for="item-1-0-0">Part 1</label>
		                </li>
		                <li><input type="checkbox" id="item-1-0-1" />
		                	<label for="item-1-0-1">Part 2</label>
		                </li>
		                <li><input type="checkbox" id="item-1-0-2" />
		                	<label for="item-1-0-2">Part 3</label>
		                </li>
		                <li><input type="checkbox" id="item-1-0-3" />
		                	<label for="item-1-0-3">Part 4</label>
		                </li>
		                <li><input type="checkbox" id="item-1-0-4" />
		                	<label for="item-1-0-4">Part 5</label>
		                </li>
                	</ul>
                </li>
                <li><input type="checkbox" id="item-1-2" />
                	<label for="item-1-2">Up Dimple</label>
                	<ul>
		                <li><input type="checkbox" id="item-1-0-0" />
		                	<label for="item-1-0-0">Part 1</label>
		                </li>
		                <li><input type="checkbox" id="item-1-0-1" />
		                	<label for="item-1-0-1">Part 2</label>
		                </li>
		                <li><input type="checkbox" id="item-1-0-2" />
		                	<label for="item-1-0-2">Part 3</label>
		                </li>
		                <li><input type="checkbox" id="item-1-0-3" />
		                	<label for="item-1-0-3">Part 4</label>
		                </li>
		                <li><input type="checkbox" id="item-1-0-4" />
		                	<label for="item-1-0-4">Part 5</label>
		                </li>
                	</ul>
                </li>
                <li><input type="checkbox" id="item-1-3" />
                	<label for="item-1-3">Half Notch</label>
                </li>
                <li><input type="checkbox" id="item-1-4" />
                	<label for="item-1-4">Full Notch</label>
                </li>
                <li><input type="checkbox" id="item-1-5" />
                	<label for="item-1-5">Right Angle Guillo</label>
                </li>
                <li><input type="checkbox" id="item-1-6" />
                	<label for="item-1-6">Straight Guillo</label>
                </li>
                <li><input type="checkbox" id="item-1-7" />
                	<label for="item-1-7">Left Angle Guillo</label>
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
