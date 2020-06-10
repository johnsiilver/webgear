package html

import (
	"html/template"
	"strings"
	"sync"
)

var liTmpl = template.Must(template.New("li").Parse(strings.TrimSpace(`
<li {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</li>
`)))

// Li defines an HTML li tag.
type Li struct {
	GlobalAttrs
	Events *Events

	Elements []Element

	pool sync.Pool
}

func (l *Li) isElement() {}

func (l *Li) Init() error {
	l.pool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}
	return nil
}

func (l *Li) Execute(pipe Pipeline) template.HTML {
	buff := l.pool.Get().(*strings.Builder)
	defer l.pool.Put(buff)
	buff.Reset()

	pipe.Self = l

	if err := liTmpl.Execute(buff, pipe); err != nil {
		panic(err)
	}

	return template.HTML(buff.String())
}
