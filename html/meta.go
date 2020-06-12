package html

import (
	"fmt"
	"html/template"
	"strings"
)

var metaTmpl = template.Must(template.New("meta").Parse(strings.TrimSpace(`
<meta {{.Self.Attr}} {{.Self.GlobalAttrs.Attr}}>
`)))

type HTTPEquiv string

const (
	ContentSecurityPolicyHE HTTPEquiv = "content-security-policy"
	ContentTypeHE           HTTPEquiv = "content-type"
	DefaultStyleHE          HTTPEquiv = "default-style"
	RefreshHE               HTTPEquiv = "refresh"
)

type MetaName string

const (
	ApplicationNameMN MetaName = "application-name"
	DescriptionMN     MetaName = "description"
	GeneratorMN       MetaName = "generator"
	KeywordsMN        MetaName = "keywords"
	ViewportMN        MetaName = "viewport"
)

// Meta defines an HTML meta tag.
type Meta struct {
	GlobalAttrs

	// Charset holds the character encoding of the html document. We only support the value UTF-8.
	Charset string

	HTTPEquiv HTTPEquiv `html:"http-equiv"`

	MetaName MetaName

	// Content specifies the value associated with the http-equiv or name attribute.
	Content string
}

func (m *Meta) validate() error {
	if m == nil {
		return nil
	}

	attrs := []string{
		m.Charset,
		string(m.MetaName),
		string(m.HTTPEquiv),
	}

	i := 0
	for _, a := range attrs {
		if a != "" {
			i++
		}
	}
	if i != 1 {
		return fmt.Errorf("one and only one value of Meta.Charset/MetaName/HTTPEquiv can be set per Meta")
	}

	switch {
	case m.Charset != "":
		if m.Charset != "UTF-8" {
			return fmt.Errorf("Meta.Charset can only be UTF-8")
		}
	case m.HTTPEquiv != "" || m.MetaName != "":
		if m.Content == "" {
			return fmt.Errorf("HTTPEquiv cannot be set if Content is not set")
		}
	}

	return nil
}

func (m *Meta) Attr() template.HTMLAttr {
	output := structToString(m)
	return template.HTMLAttr(output)
}

func (m *Meta) Execute(pipe Pipeline) string {
	pipe.Self = m

	if err := metaTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}
