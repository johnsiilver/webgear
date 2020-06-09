package html

import (
	"fmt"
	"html/template"
	"strings"
	"sync"
)

var pTmpl = strings.TrimSpace(`
<p {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</p>
`)

// P tag defines a paragraph.
type P struct {
	GlobalAttrs
	Events *Events

	Elements []Element

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

func (p *P) Execute(pipe Pipeline) template.HTML {
	buff := p.pool.Get().(*strings.Builder)
	defer p.pool.Put(buff)
	buff.Reset()

	pipe.Self = p

	if err := p.tmpl.Execute(buff, pipe); err != nil {
		panic(err)
	}

	return template.HTML(buff.String())
}
