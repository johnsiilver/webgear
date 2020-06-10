package html

import (
	"html/template"
	"strings"
	"sync"
)

var ulTmpl = template.Must(template.New("ul").Parse(strings.TrimSpace(`
<ul {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</ul>
`)))

// Ul defines an HTML ul tag.
type Ul struct {
	GlobalAttrs
	Events *Events

	Elements []Element

	pool sync.Pool
}

func (u *Ul) isElement() {}

func (u *Ul) Init() error {
	u.pool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}
	return nil
}

func (u *Ul) Execute(pipe Pipeline) template.HTML {
	buff := u.pool.Get().(*strings.Builder)
	defer u.pool.Put(buff)
	buff.Reset()

	pipe.Self = u

	if err := ulTmpl.Execute(buff, pipe); err != nil {
		panic(err)
	}

	return template.HTML(buff.String())
}
