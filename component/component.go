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
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"sync"

	"github.com/johnsiilver/webgear/html"
)

var htmlTemplateTxt = `
{{ define "template" }}
<template id="{{.Self.Name}}Template">
	{{.Self.Doc.ExecuteAsGear .}}
</template>
{{ end }}
`

var scriptTemplateTxt = `
{{ define "script" }}
<script>
	function {{.Self.LoaderName}}() {
		if (!window.customElements.get('{{.Self.Name}}')) {
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
		}
		let old = document.getElementById("{{.Self.Name}}");
		if (old !== null) {
			let newcomp = old.cloneNode(true);
			document.body.replaceChild(newcomp, old);
		}
	}
	console.log("hello world");
	{{.Self.LoaderName}}();
</script>
{{ end }}
`

var combinedTxt = `
{{ template "template" . }}

{{ template "script" . }}
`

var justTemplate = `
{{ template "template" . }}
`

var gearTmpl *template.Template

func init() {
	gearTmpl = template.Must(template.New("htmlTemplate").Parse(htmlTemplateTxt))
	gearTmpl = template.Must(gearTmpl.New("scriptTemplate").Parse(scriptTemplateTxt))
	gearTmpl = template.Must(gearTmpl.New("combinedTxt").Parse(combinedTxt))
	gearTmpl = template.Must(gearTmpl.New("justTemplate").Parse(justTemplate))
}

// DataFunc represents a function that provides data in the html.Pipeline.GearData. The DataFunc should
// return data that will be stored in the html.Pipeline.GearData field. The returned object must be thread-safe.
type DataFunc func(r *http.Request) (interface{}, error)

// Gear is a shadow-dom component.
type Gear struct {
	// Doc is public to allow its use in internal templating code. It should only be set by the call to New().
	Doc *html.Doc
	// Gears is public to allow external packages to access Gear data. This houls only be set by AddGear().
	Gears    []*Gear
	dataFunc DataFunc

	name       string
	loaderName string

	wasmUpdateMu sync.Mutex
	wasmUpdate   bool
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
// other components. You still must use html.Component{} to insert your custom tag where you want the componenet to be displayed.
func AddGear(newGear *Gear) Option {
	return func(g *Gear) {
		g.Gears = append(g.Gears, newGear)
	}
}

// New creates a new Gear object called "name" using the HTML provided by the doc passed. Name also is used for the tag's
// ID when using the html.Component{} type.
func New(name string, doc *html.Doc, options ...Option) (*Gear, error) {
	if err := validName(name); err != nil {
		return nil, err
	}

	doc.Component = true
	doc.Pretty = false // Inside a Gear, this should always be false.

	if err := doc.Init(); err != nil {
		return nil, err
	}

	walkCtx, cancel := context.WithCancel(context.Background())
	for walked := range html.Walker(walkCtx, doc.Body) {
		if g, ok := walked.Element.(*Gear); ok {
			cancel()
			return nil, fmt.Errorf("WebGear Component(%s) had another component(%s) added directly to the passed *html.Doc,"+
				"this can only be added using the component.AddGear() option to allow correct rendered ordering", name, g.name)
		}
	}

	g := &Gear{
		Doc:        doc,
		name:       name,
		loaderName: strings.ReplaceAll(name, "-", "") + "Loader",
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

// TagType is the same as Name() except the type works in the templates.
func (g *Gear) TagType() template.HTMLAttr {
	return template.HTMLAttr(g.name)
}

// GearID outputs the ID of the gear stored in the .Doc.  Gear's are special and are never output
// to HTML with this name (output with TemplateName() and LoaderName()), but we still need a way
// to reference them in our Doc tree. This is that ID.
func (g *Gear) GearID() string {
	return "gear-" + g.name
}

// LoaderName is the name of the JS function that defines the HTML custom element and causes the custom element to render.
func (g *Gear) LoaderName() template.JS {
	return template.JS(g.loaderName)
}

// TemplateName is the DOM id of the template that is used for this component.
func (g *Gear) TemplateName() string {
	return g.name + "Template"
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
	for _, gear := range g.Gears {
		gear.Execute(pipe)
		if pipe.Ctx.Err() != nil {
			return html.EmptyString
		}
	}

	err = gearTmpl.ExecuteTemplate(pipe.W, "combinedTxt", pipe)
	//err = gearTmpl.Execute(pipe.W, pipe)
	if err != nil {
		panic(err)
	}

	return html.EmptyString
}

// TemplateContent outputs the content of the HTML template object used for this component, but not the script.
func (g *Gear) TemplateContent() string {
	buff := bytes.Buffer{}
	pipe := html.NewPipeline(context.Background(), &http.Request{}, nil)
	pipe.Self = g
	pipe.W = &buff

	err := gearTmpl.ExecuteTemplate(pipe.W, "justTemplate", pipe)
	if err != nil {
		panic(err)
	}
	return buff.String()
}

func validName(s string) error {
	hasHyphen := false
	if len(s) == 0 {
		return fmt.Errorf("component name cannot have an empty name")
	}
	if s[0] < 97 || s[0] > 122 {
		return fmt.Errorf("component name cannot have first letter that is not a lowercase alpha character")
	}

	for i := 1; i < len(s); i++ {
		switch {
		case s[i] == 45:
			hasHyphen = true
		case s[i] < 48: // Non numeric
			return fmt.Errorf("component name cannot contain a non-numeric or non-ascii lower case letter, such as %q in %q", s[i], s)
		case s[i] > 57: // Possible non-alpha
			if s[i] < 97 || s[i] > 122 { // lower case ascii
				return fmt.Errorf("component name cannot contain this non-numeric non-lower case letter %q in %q", s[i], s)
			}
		}
	}

	if !hasHyphen {
		return fmt.Errorf("a componenent name must have a - in it, don't blame me, blame the spec")
	}
	return nil
}
