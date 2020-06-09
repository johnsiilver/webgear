package html

import (
	"fmt"
	"html/template"
	"net/url"
	"strings"
	"sync"
)

var scriptTmpl = strings.TrimSpace(`
<script {{.Self.Attr}} {{.Self.GlobalAttrs.Attr}}>
	{{.Self.TagValue}}
</script>
`)

// Script represents an HTML script tag.
type Script struct {
	GlobalAttrs

	// Src specifies the URL of an external script file.
	Src *url.URL

	// Type specifies the media type of the script.
	Type string

	// Async specifies that the script is executed asynchronously (only for external scripts).
	Async bool `html:"attr"`

	// Defer specifies that the script is executed when the page has finished parsing (only for external scripts).
	Defer bool `html:"attr"`

	// TagValue holds the value that is between the begin and ending tag. This should be a script of some type.
	TagValue template.JS

	tmpl *template.Template

	pool sync.Pool
}

func (s *Script) isElement() {}

func (s *Script) Attr() template.HTMLAttr {
	output := structToString(s)
	return template.HTMLAttr(output)
}

func (s *Script) compile() error {
	var err error
	s.tmpl, err = template.New("script").Parse(scriptTmpl)
	if err != nil {
		return fmt.Errorf("Script object had error: %s", err)
	}

	s.pool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}

	return nil
}

func (s *Script) Execute(pipe Pipeline) template.HTML {
	buff := s.pool.Get().(*strings.Builder)
	defer s.pool.Put(buff)
	buff.Reset()

	pipe.Self = s

	if err := s.tmpl.Execute(buff, pipe); err != nil {
		panic(err)
	}

	return template.HTML(buff.String())
}
