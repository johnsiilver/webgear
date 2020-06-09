package html

import (
	"fmt"
	"strings"
	"sync"
	"text/template"

	html "html/template"
)

// Attribute provides a custom attribute for a user to provide for custom componenets.
type Attribute interface {
	fmt.Stringer
	IsAttr()
}

var componenetTmpl = strings.TrimSpace(`
<{{.Self.TagType}} {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{.Self.TagValue}}
</{{.Self.TagType}}>
`)

// Component is for providing custom componenets registered through the javascript window.customElements type.
type Component struct {
	GlobalAttrs

	// Attributes are custom attributes to apply to the component.
	Attributes []Attribute

	// TageType (required) is the name of the custom componenet tag, like "myComponent".
	TagType html.HTMLAttr

	// TagValue provides the value inside a reference. This can be another element such as a div, but most commonly
	// it is not defined.
	TagValue Element

	Events *Events

	tmpl *template.Template

	pool sync.Pool
}

func (c *Component) validate() error {
	if c.TagType == "" {
		return fmt.Errorf("Component.TagValue cannot be empty")
	}

	if !strings.Contains(string(c.TagType), "-") {
		return fmt.Errorf("Components.TagValue does not have a '-' in the name, as required by the spec")
	}
	return nil
}

func (c *Component) isElement() {}

func (c *Component) compile() error {
	var err error
	c.tmpl, err = template.New("c").Parse(componenetTmpl)
	if err != nil {
		return fmt.Errorf("Component object had error: %s", err)
	}

	c.pool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}
	return nil
}

func (c *Component) Execute(pipe Pipeline) html.HTML {
	buff := c.pool.Get().(*strings.Builder)
	defer c.pool.Put(buff)
	buff.Reset()

	pipe.Self = c

	if err := c.tmpl.Execute(buff, pipe); err != nil {
		panic(err)
	}

	return html.HTML(buff.String())
}
