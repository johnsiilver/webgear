// Package component provides a Gear, which represents an HTML shadow-dom component.
package component

import (
	"fmt"
	"html/template"
	"strings"
	"sync"

	"github.com/johnsiilver/webgear/html"

	"github.com/yosssi/gohtml"
)

var gearTmpl = strings.TrimSpace(`
<template id="{{.Self.Name}}Template">
	{{.Self.Doc.Render .}}
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

// Gear is a shadow-dom component.
type Gear struct {
	// Doc is public to allow its use in internal templating code. It should only be set by the call to New().
	Doc   *html.Doc
	gears []*Gear

	// Pretty says to make the HTML look pretty before outputting.
	Pretty bool

	pool sync.Pool

	name string

	tmpl *template.Template
}

// New creates a new Gear object called "name" using the HTML provided by the doc passed.
// This will call Compile() on the *html.Doc. "gears" provides other componenets that this Gear
// will execute before executing its own doc when rendering.  This allows a componenet to use
// other components.
func New(name string, doc *html.Doc, gears []*Gear) (*Gear, error) {
	if name == "" {
		return nil, fmt.Errorf("must provide a name for the Gear")
	}

	doc.Component = true
	doc.Pretty = false // If they want pretty, they need to set it in the Gear.

	if err := doc.Compile(); err != nil {
		return nil, err
	}

	g := Gear{
		Doc:  doc,
		name: name,
		pool: sync.Pool{
			New: func() interface{} {
				return &strings.Builder{}
			},
		},
	}

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

// Execute executes the internal templates and renders the html for output with the given pipeline.
func (g *Gear) Execute(pipe html.Pipeline) (template.HTML, error) {
	if g.tmpl == nil {
		return "", fmt.Errorf("Gear.Execute() called before Gear.Compile()")
	}

	w := g.pool.Get().(*strings.Builder)
	defer g.pool.Put(w)
	w.Reset()

	pipe.Self = g

	var err error
	for _, gear := range g.gears {
		h, err := gear.Execute(pipe)
		if err != nil {
			return "", err
		}
		w.WriteString(string(h))
	}

	if g.Pretty {
		err = g.tmpl.ExecuteTemplate(gohtml.NewWriter(w), "gear", pipe)
	} else {
		err = g.tmpl.ExecuteTemplate(w, "gear", pipe)
	}
	if err != nil {
		return "", err
	}
	return template.HTML(w.String()), nil
}
