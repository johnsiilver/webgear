// Package component provides a Gear, which represents an HTML shadow-dom component.
package component

import (
	"fmt"
	"html/template"
	"strings"
	"sync"

	"github.com/johnsiilver/webgear/html"

	"github.com/google/uuid"
	"github.com/yosssi/gohtml"
)

var gearTmpl = strings.TrimSpace(`
<template id="{{.Self.Name}}Template">
	{{.Self.Doc.Render .Data}}
</template>

<script>
	window.customElements.define('{{.Self.Name}}',
		class extends HTMLElement {
			constructor() {
				super();
				let template = document.getElementById('{{.Self.Name}}');
				let templateContent = template.content;

				const shadowRoot = this.attachShadow({mode: 'open'}).appendChild(templateContent.cloneNode(true));
			}
		}
	);
</script>
`)

type pipeline struct {
	Self interface{}
	Data interface{}
}

// Gear is a shadow-dom component. It has a randomly created unique ID that is found by using .Name().
type Gear struct {
	Doc *html.Doc

	// Pretty says to make the HTML look pretty before outputting.
	Pretty bool

	pool sync.Pool

	namePrefix string
	name       string

	tmpl *template.Template
}

// NewGear creates a new Gear object called "name" using the HTML provided by the doc passed.
// This will call Compile() on the *html.Doc.
func NewGear(name string, doc *html.Doc) (*Gear, error) {
	if name == "" {
		return nil, fmt.Errorf("must provide a name for the Gear")
	}

	doc.Component = true
	doc.Pretty = false // If they want pretty, they need to set it in the Gear.

	if err := doc.Compile(); err != nil {
		return nil, err
	}

	g := Gear{
		Doc:        doc,
		namePrefix: name,
		pool: sync.Pool{
			New: func() interface{} {
				return &strings.Builder{}
			},
		},
	}

	id := uuid.New()

	g.name = g.namePrefix + "-" + id.String()

	if err := g.compile(); err != nil {
		return nil, err
	}

	return &g, nil
}

// Name returns the name of the Gear so that it may be referenced.
func (g *Gear) Name() string {
	return g.name
}

// compile compiles the internal templates before execution. This should only be done once.
func (g *Gear) compile() error {
	gear, err := template.New("gear").Parse(gearTmpl)
	if err != nil {
		return err
	}
	g.tmpl = gear
	return nil
}

// Execute executes the internal templates and renders the html for output with the given "data" pipeline.
// If you are unfamiliar with data pipelines, see https://golang.org/pkg/text/template/.
func (g *Gear) Execute(data interface{}) (template.HTML, error) {
	if g.tmpl == nil {
		return "", fmt.Errorf("Gear.Execute() called before Gear.Compile()")
	}

	w := g.pool.Get().(*strings.Builder)
	defer g.pool.Put(w)
	w.Reset()

	var err error
	if g.Pretty {
		err = g.tmpl.ExecuteTemplate(gohtml.NewWriter(w), "gear", pipeline{Self: g, Data: data})
	} else {
		err = g.tmpl.ExecuteTemplate(w, "gear", pipeline{Self: g, Data: data})
	}
	if err != nil {
		return "", err
	}
	return template.HTML(w.String()), nil
}
