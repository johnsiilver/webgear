package html

import (
	"html/template"
	"strings"
	"sync"
)

var pTmpl = template.Must(template.New("p").Parse(strings.TrimSpace(`
<p {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</p>
`)))

// P tag defines a paragraph.
type P struct {
	GlobalAttrs
	Events *Events

	Elements []Element

	pool sync.Pool
}

func (p *P) isElement() {}

func (p *P) Attr() template.HTMLAttr {
	output := structToString(p)
	return template.HTMLAttr(output)
}

func (p *P) Init() error {
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

	if err := pTmpl.Execute(buff, pipe); err != nil {
		panic(err)
	}

	return template.HTML(buff.String())
}
