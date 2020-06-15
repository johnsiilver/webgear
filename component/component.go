/*
Package component provides a Gear, which represents an HTML shadow-dom component. This allows creation of isolated
HTML components using the webgear/html package.

Simply attach an *html.Doc element with a .Body and no .Head. This will be automatically isolated from other css
styles in the top level document.  The component tag name will be the name you pass along in the constructor and you
can access this component in your main document via the html.Component{} object.
*/
//
// Prerequisites
//
/*
To properly understand this package, you will need the following:
	* Understanding of the webgears/html package
	* Understand the idea of the Shadow-DOM in Web Components

If you don't know about web components, a good introduction can be found here: https://developer.mozilla.org/en-US/docs/Web/Web_Components/Using_shadow_DOM
*/
//
// Example Component
//
// This creates a component that simply prints a name and attaches that to an html.Doc object for rendering.
/*
	type dynName struct {
		name string
	}

	func(d dynName) Name() []html.Element {
		return []html.Element{html.TextElement(d.name)}
	}

	// New is our custom Gear's constructor that will print out the name we pass in via printName.
	// The compName will be used to specify our custom HTML tag that allows exeuction of the component.
	// If compName is "my-component", the the <my-component></my-component> tag is used to execute this component.
	func New(compName, printName string) (*component.Gear, error) {
		// Create doc that will be used by the Gear.
		var doc := &html.Doc{
			Body: &html.Body{
				Elements: []html.Element{
					&html.Link{Rel: "stylesheet", Href: "/static/main/gear.css"},
					html.Dynamic(dynName{"John Doak"}.Name),
				},
			},
		}

		return component.New(name, doc)
	}
*/
//
// Attaching a Component
//
// To use a componenet in a page, it needs to be attached to an html.Doc to be rendered.
/*
	gear, err := printname.New("print-name-author", "John Doak")
	if err != nil {
		// Do something
	}

	// Use the Gear in your index page. This is usually not in the same package as the component.
	doc := &html.Doc{
		Head: &html.Head{
			&html.Meta{Charset: "UTF-8"},
			&html.Title{TagValue: html.TextElement("My site showing my name")},
			&html.Link{Rel: "stylesheet", Href: html.URLParse("/static/main/index.css")},
			&html.Link{Href: html.URLParse("https://fonts.googleapis.com/css2?family=Share+Tech+Mono&display=swap"), Rel: "stylesheet"},
		},
		Body: &html.Body{
			Elements: []html.Element{
				// This causes the gear's generated code to be written to output.
				gear,
				// This is the gear tag that causes the code to be executed, aka <print-name-author></print-name-author>.
				&html.Component{TagType: template.HTMLAttr(gear.Name())},
			},
		},
	},
*/
//
// Serving a page
//
// Now we need to serve the page and any external file required such as images or css files.
/*
	// Setup server handlers.
	h := handlers.New(handlers.DoNotCache())

	// Serve all files from the the binary working directory and below it (recursively) that have
	// the file extensions listed.
	h.ServeFilesWorkingDir([]string{".css", ".jpg", ".svg", ".png"})

	// Attach our page containing the gear to "/".
	h.MustHandle("/", doc)

	// Serve the content using the http.Server.
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", *port),
		Handler:        h.ServerMux(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("http server serving on :%d", *port)

	log.Fatal(server.ListenAndServe())
*/
package component

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/johnsiilver/webgear/html"
)

var gearTmpl = template.Must(template.New("gear").Parse(strings.TrimSpace(`
<template id="{{.Self.Name}}Template">
	{{.Self.Doc.ExecuteAsGear .}}
</template>

<script>
  window.customElements.define(
		'{{.Self.Name}}',
		class extends HTMLElement {
			constructor() {
				super();
				let template = document.getElementById('{{.Self.Name}}Template');
				let templateContent = template.content;

				const shadowRoot = this.attachShadow({mode: 'open'}).appendChild(templateContent.cloneNode(true));
			}
		}
	);
</script>
`)))

// DataFunc represents a function that provides data in the html.Pipeline.GearData. The DataFunc should
// return data that will be stored in the html.Pipeline.GearData field. The returned object must be thread-safe.
type DataFunc func(r *http.Request) (interface{}, error)

// Gear is a shadow-dom component.
type Gear struct {
	// Doc is public to allow its use in internal templating code. It should only be set by the call to New().
	Doc      *html.Doc
	gears    []*Gear
	dataFunc DataFunc

	name string
}

// Option is an optional argument to the New() constructor.
type Option func(g *Gear)

// ApplyDataFunc adds a function that populates the html.Pipeline.GearData attribute on Execute() calls.
func ApplyDataFunc(f DataFunc) Option {
	return func(g *Gear) {
		if g.dataFunc != nil {
			panic("cannot use ApplyDataFunc() more than once")
		}
		g.dataFunc = f
	}
}

// AddGear adds another Gear that will be called before this gear is called.  This allows a componenet to use
// other components.
func AddGear(newGear *Gear) Option {
	return func(g *Gear) {
		g.gears = append(g.gears, newGear)
	}
}

// New creates a new Gear object called "name" using the HTML provided by the doc passed.
func New(name string, doc *html.Doc, options ...Option) (*Gear, error) {
	if name == "" {
		return nil, fmt.Errorf("must provide a name for the Gear")
	}

	if !strings.Contains(name, "-") {
		return nil, fmt.Errorf("a componenent name must have a - in it, don't blame me, blame the spec")
	}

	doc.Component = true
	doc.Pretty = false // Inside a Gear, this should always be false.

	if err := doc.Init(); err != nil {
		return nil, err
	}

	g := &Gear{
		Doc:  doc,
		name: name,
	}

	for _, o := range options {
		o(g)
	}

	return g, nil
}

// Name returns the name of the Gear so that it may be referenced.
func (g *Gear) Name() string {
	return g.name
}

// Execute executes the internal templates and renders the html for output with the given pipeline.
func (g *Gear) Execute(pipe html.Pipeline) string {
	pipe.Self = g

	if g.dataFunc != nil {
		i, err := g.dataFunc(pipe.Req)
		if err != nil {
			panic(err)
		}
		pipe.GearData = i
	}

	var err error
	for _, gear := range g.gears {
		gear.Execute(pipe)
		if pipe.Ctx.Err() != nil {
			return html.EmptyString
		}
	}

	err = gearTmpl.Execute(pipe.W, pipe)
	if err != nil {
		panic(err)
	}

	return html.EmptyString
}
