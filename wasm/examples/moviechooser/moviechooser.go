package main

import (
	"context"
	"html/template"
	"log"

	"github.com/johnsiilver/webgear/wasm"
	"github.com/johnsiilver/webgear/wasm/examples/moviechooser/components/categories"
	"github.com/johnsiilver/webgear/wasm/examples/moviechooser/components/list"

	. "github.com/johnsiilver/webgear/html"
)

func main() {
	w := wasm.New()

	listGear, err := list.New("list-component", nil)
	if err != nil {
		log.Println(err)
		return
	}

	categoriesGear, err := categories.New("categories-component", listGear.Name(), w)
	if err != nil {
		panic(err)
	}

	doc := &Doc{
		Head: &Head{
			Elements: []Element{
				&Meta{Charset: "UTF-8"},
				&Title{TagValue: TextElement("Movie Chooser")},
				&Link{Rel: "stylesheet", Href: URLParse("/static/moviechooser.css")},
			},
		},
		Body: &Body{
			Elements: []Element{
				// Renders the component code.
				categoriesGear,
				listGear,
				// Renders the component tag that will be where the output is displayed.
				&Component{Gear: categoriesGear},
				&Component{Gear: listGear},
				&Script{TagValue: template.JS("console.log('js view:', document.documentElement.innerHTML)")},
			},
		},
	}
	w.SetDoc(doc)

	w.Run(context.Background())
}
