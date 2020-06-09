package html

import (
	"fmt"
	"html/template"
	"net/url"
	"strings"
	"sync"
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

	pool sync.Pool
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

	b.pool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}

	return nil
}

func (b *Base) Execute(pipe Pipeline) template.HTML {
	buff := b.pool.Get().(*strings.Builder)
	defer b.pool.Put(buff)
	buff.Reset()

	if err := b.tmpl.Execute(buff, Pipeline{Self: b}); err != nil {
		panic(err)
	}

	return template.HTML(buff.String())
}
