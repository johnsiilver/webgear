package html

import (
	"fmt"
	"html/template"
	"strings"
	"net/url"
)

var baseTmpl = strings.TrimSpace(`
<base {{.Self.Attr}} {{.Self.GlobalAttrs.Attr}}>
`)

// Base represents an HTML script tag.
type Base struct {
	GlobalAttrs

	// Href specifies the URL for all relative urls in the page.
	Href *url.URL

	// Target specifies the default target for all hyperlinks and forms in the page.
	Target string

	tmpl *template.Template

	str string
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

func (b *Base) compile() error {
	var err error
	b.tmpl, err = template.New("base").Parse(baseTmpl)
	if err != nil {
		return fmt.Errorf("Base object had error: %s", err)
	}

	return nil
}

func (b *Base) Execute(data interface{}) template.HTML {
	if b.str != "" {
		return template.HTML(b.str)
	}

	buff := strings.Builder{}

	if err := b.tmpl.Execute(&buff, pipeline{Self: b, Data: data}); err != nil {
		panic(err)
	}

	b.str = buff.String()
	return template.HTML(b.str)
}