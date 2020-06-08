package html

import (
	"fmt"
	"html/template"
	"strings"
	"sync"
)

var pTmpl = strings.TrimSpace(`
<p {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{.Self.Element.Execute .Data}}
</p>
`)

// P tag defines a paragraph.
type P struct {
	GlobalAttrs
	Events *Events

	Element Element

	tmpl *template.Template

	pool sync.Pool
}

func (p *P) isElement() {}

func (p *P) Attr() template.HTMLAttr {
	output := structToString(p)
	return template.HTMLAttr(output)
}

func (p *P) compile() error {
	var err error
	p.tmpl, err = template.New("p").Parse(pTmpl)
	if err != nil {
		return fmt.Errorf("P object had error: %s", err)
	}

	p.pool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}
	return nil
}

func (p *P) Execute(data interface{}) template.HTML {
	buff := p.pool.Get().(*strings.Builder)
	defer p.pool.Put(buff)
	buff.Reset()

	if err := p.tmpl.Execute(buff, pipeline{Self: p, Data: data}); err != nil {
		panic(err)
	}

	return template.HTML(buff.String())
}
