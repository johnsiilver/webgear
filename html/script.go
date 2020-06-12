package html

import (
	"html/template"
	"net/url"
	"strings"
)

var scriptTmpl = template.Must(template.New("script").Parse(strings.TrimSpace(`
<script {{.Self.Attr}} {{.Self.GlobalAttrs.Attr}}>
	{{.Self.TagValue}}
</script>
`)))

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
}

func (s *Script) Attr() template.HTMLAttr {
	output := structToString(s)
	return template.HTMLAttr(output)
}

func (s *Script) Execute(pipe Pipeline) string {
	pipe.Self = s

	if err := scriptTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}
