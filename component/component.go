// Package component provides a Gear, which represents an HTML shadow-dom component.
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
