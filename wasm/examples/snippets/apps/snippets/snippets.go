package main

import (
	"context"
	"log"

	"github.com/johnsiilver/webgear/wasm"
	"github.com/johnsiilver/webgear/wasm/examples/snippets/apps/snippets/components/banner"
	"github.com/johnsiilver/webgear/wasm/examples/snippets/apps/snippets/components/calendar"
	"github.com/johnsiilver/webgear/wasm/examples/snippets/apps/snippets/components/content"

	. "github.com/johnsiilver/webgear/html"
)

func main() {
	w := wasm.New()

	bannerGear, err := banner.New("banner-component")
	if err != nil {
		panic(err)
	}

	contentGear, err := content.New("content-component", content.Args{Mode: content.View, RestEndpoint: "127.0.0.1:8081"}, w)
	if err != nil {
		panic(err)
	}

	calendarGear, err := calendar.New(
		"calendar-component",
		"content-component",
		calendar.Args{
			CSSPath: "/static/apps/snippets/components/calendar/calendar.css",
		},
		w,
	)
	if err != nil {
		log.Println(err)
		return
	}

	doc := &Doc{
		Head: &Head{
			Elements: []Element{
				&Meta{Charset: "UTF-8"},
				&Title{TagValue: TextElement("Snippets App")},
				&Link{Rel: "stylesheet", Href: URLParse("/static/apps/snippets/snippets.css")},
			},
		},
		Body: &Body{
			Elements: []Element{
				// Renders the component code.
				bannerGear,
				calendarGear,
				contentGear,

				// Below we render the component tags that will be where the output is displayed.
				// The content is dynamically generated using the template code that the
				// above *Gear Elements inserted into the page along with javascript that registered
				// those templates to the component tags below (also called customElements).
				// Each component below gets an id=gear.Name(), aka the Component tag for
				// bannerGear looks like <banner-component id="banner-component"></banner-component>.
				&Component{Gear: bannerGear},
				&Component{Gear: calendarGear},
				&Component{Gear: contentGear},
			},
		},
	}
	log.Println("got here")
	w.SetDoc(doc)

	w.Run(context.Background())
}
