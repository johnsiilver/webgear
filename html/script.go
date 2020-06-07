package html

import (
	"fmt"
	"html/template"
	"strings"
	"net/url"
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
	Defer 	bool `html:"attr"` 	

	// TagValue holds the value that is between the begin and ending tag. This should be a script of some type.
	TagValue template.JS

	tmpl *template.Template

	str string
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

	return nil
}

func (s *Script) Execute(data interface{}) template.HTML {
	if s.str != "" {
		return template.HTML(s.str)
	}

	buff := strings.Builder{}

	if err := s.tmpl.Execute(&buff, pipeline{Self: s, Data: data}); err != nil {
		panic(err)
	}

	s.str = buff.String()
	return template.HTML(s.str)
}