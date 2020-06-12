package html

import (
	"fmt"
	"html/template"
	"net/url"
	"strings"
)

var baseTmpl = template.Must(template.New("base").Parse(strings.TrimSpace(`
<base {{.Self.Attr}} {{.Self.GlobalAttrs.Attr}}>
`)))

// Base represents an HTML script tag.
type Base struct {
	GlobalAttrs

	// Href specifies the URL for all relative urls in the page.
	Href *url.URL

	// Target specifies the default target for all hyperlinks and forms in the page.
	Target string
}

func (b *Base) isElement() {}

func (b *Base) validate() error {
	if b.Href.String() == "" && b.Target == "" {
		return fmt.Errorf("Base tag must have either/both Href and Target attributes")
	}
	return nil
}

func (b *Base) Attr() template.HTMLAttr {
	output := structToString(b)
	return template.HTMLAttr(output)
}

func (b *Base) Execute(pipe Pipeline) string {

	if err := baseTmpl.Execute(pipe.W, Pipeline{Self: b}); err != nil {
		panic(err)
	}

	return EmptyString
}
