package html

import (
	"fmt"
	"strings"
	"text/template"
)

// Attribute provides a custom attribute for a user to provide for custom componenets.
type Attribute interface {
	fmt.Stringer
	IsAttr()
}

var componenetTmpl = template.Must(template.New("component").Parse(strings.TrimSpace(`
<{{.Self.Gear.TagType}} {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{with .Self.TagValue}}{{.}}{{end}}
</{{.Self.Gear.TagType}}>
`)))

// Component is for providing custom componenets registered through the javascript window.customElements type.
// Be aware that the .Gear.Name() will override a provided GlobalAttrs value.
type Component struct {
	GlobalAttrs

	// Attributes are custom attributes to apply to the component.
	Attributes []Attribute

	// Gear is the *component.Gear that implements the componenent. The name of that Gear will be both the tag type and
	// the id of the Gear.
	Gear GearType

	// TagValue provides the value inside a reference. This can be another element such as a div, but most commonly
	// it is not defined.
	TagValue Element

	Events *Events
}

func (c *Component) validate() error {
	if c.Gear == nil {
		return fmt.Errorf("Component.Gear cannot be empty")
	}

	return nil
}

func (c *Component) isElement() {}

func (c *Component) Execute(pipe Pipeline) string {
	pipe.Self = c

	ga := c.GlobalAttrs
	ga.ID = string(c.Gear.TagType())
	c.GlobalAttrs = ga

	if err := componenetTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}
