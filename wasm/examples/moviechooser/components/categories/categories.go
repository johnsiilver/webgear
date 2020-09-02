package categories

import (
	"bytes"
	"fmt"
	"log"
	"syscall/js"

	"github.com/johnsiilver/webgear/component"
	"github.com/johnsiilver/webgear/html/builder"
	"github.com/johnsiilver/webgear/wasm"
	"github.com/johnsiilver/webgear/wasm/examples/moviechooser/components/list"

	. "github.com/johnsiilver/webgear/html"
)

var categories = []string{
	"Action",
	"Drama",
	"Romance",
	"Sci-Fi/Fantasy",
}

// Filter provides a way of filtering our list by updating our list compoonent.
type Filter struct {
	listName string
	wasm     *wasm.Wasm
	buff     *bytes.Buffer
	options  []component.Option
}

// List will change the list in our "list" component to filter out movies not in categories we are interested in.
func (f Filter) List(this js.Value, root js.Value) {
	if this.IsUndefined() {
		log.Println("Filter.List() received 'this' argument that was undefined")
		return
	}

	filterOpt := this.Get("options").Index(this.Get("selectedIndex").Int()).Get("value").String()
	if filterOpt == "" {
		return
	}
	log.Println("the selected option's value: ", filterOpt)
	listGear, err := list.New(f.listName, []string{filterOpt}, f.options...)
	if err != nil {
		log.Println(err)
		return
	}

	up := wasm.GetDocUpdater()
	defer func() {
		if err == nil {
			if err := up.UpdateDOM(); err != nil {
				log.Printf("UpdateDOM() failed: %s", err)
				return
			}
		}
	}()
	if err = up.UpdateElementByID(listGear.GearID(), listGear); err != nil {
		log.Printf("Filter.List: had error trying to update the list gear element: %s", err)
		return
	}

	/*
		log.Println("listGear template name: ", listGear.TemplateName())
		js.Global().Get("document").Call("getElementById", listGear.TemplateName()).Set("outerHTML", listGear.TemplateContent())
		log.Println("wrote: \n", listGear.TemplateContent())
		js.Global().Call(string(listGear.LoaderName()))
		listGear.Doc.ExecuteDomCalls()
	*/
}

// New constructs a new component that shows a list of categories that are selectable.
func New(name string, listName string, w *wasm.Wasm, options ...component.Option) (*component.Gear, error) {
	filter := Filter{
		listName: listName,
		wasm:     w,
		buff:     bytes.NewBuffer([]byte{}),
		options:  options,
	}

	b := builder.NewHTML(&Head{}, &Body{})
	b.Into(&Div{GlobalAttrs: GlobalAttrs{ID: "leftPane"}})
	b.Into(&Form{})
	b.Into(
		wasm.AttachHandler(
			OnChange,
			false,
			filter.List,
			&Select{GlobalAttrs: GlobalAttrs{ID: "categories"}, Name: "categories"},
		),
		//&Select{GlobalAttrs: GlobalAttrs{ID: "categories"}, Name: "categories"},
	)
	for _, category := range categories {
		// Adds our option with an event attachment when somone selects something.
		b.Add(&Option{GlobalAttrs: GlobalAttrs{ID: fmt.Sprintf("%s-option", category)}, Value: category, TagValue: category})
	}
	b.Up().Up().Up()

	gear, err := component.New(name, b.Doc())
	if err != nil {
		return nil, err
	}

	return gear, nil
}
