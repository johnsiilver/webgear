package html

import (
	"html/template"
	"strings"
	"sync"
)

var navTmpl = template.Must(template.New("nav").Parse(strings.TrimSpace(`
<nav {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</nav>
`)))

// Nav defines an HTML nav tag.
type Nav struct {
	GlobalAttrs
	Events *Events

	Elements []Element

	pool sync.Pool
}

func (n *Nav) isElement() {}

func (n *Nav) Init() error {
	n.pool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}
	return nil
}

func (n *Nav) Execute(pipe Pipeline) template.HTML {
	buff := n.pool.Get().(*strings.Builder)
	defer n.pool.Put(buff)
	buff.Reset()

	pipe.Self = n

	if err := navTmpl.Execute(buff, pipe); err != nil {
		panic(err)
	}

	return template.HTML(buff.String())
}
