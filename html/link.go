

package html

import (
	"fmt"
	"html/template"
	"strings"
	"net/url"
)

type CrossOrigin string

const (
	AnnonymousCO CrossOrigin = "anonymous"
	UseCredentialsCO CrossOrigin = "use-credentials"
)

type RelationshipLink string

const (
	AlternateRL RelationshipLink = "alternate"
	AuthorRL RelationshipLink = "author"
	DNSPrefetchRL RelationshipLink = "dns-prefetch"
	HelpRL  RelationshipLink = "help"
	IconRL  RelationshipLink = "icon"
	LicenseRL RelationshipLink = "license"
	NextRL RelationshipLink = "next"
	PingBackRL RelationshipLink = "pingback"
	PreConnectRL  RelationshipLink = "preconnect"
	PreFetchRL RelationshipLink = "prefetch"
	PreLoadRL RelationshipLink = "preload"
	PreRenderRL RelationshipLink = "prerender"
	PrevRL RelationshipLink = "prev"
	SearchRL RelationshipLink = "search"
	StylesheetRL RelationshipLink = "stylesheet"
)

type Sizes struct {
	Height int
	Width int
}

func (s Sizes) outputAble() {}

func (s Sizes) String() string {
	if s.Width == 0 && s.Height == 0 {
		return ""
	}

	return fmt.Sprintf("%dx%d", s.Height, s.Width)
}

func (s Sizes) isZero() bool {
	if s.Height == 0 && s.Width == 0 {
		return true
	}
	return false
}

var linkTmpl = strings.TrimSpace(`
<link {{.Self.Attr}} {{.Self.GlobalAttrs.Attr}}>
`)


// Link defines an HTML link tag.
type Link struct {
	GlobalAttrs
	
	// Href specifies the location of the linked document.
	Href *url.URL

	CrossOrigin CrossOrigin

	// HrefLang specifies the language of the text in the linked document.
	HrefLang string

	// Media specifies on what device the linked document will be displayed.
	Media string

	ReferrerPolicy	ReferrerPolicy

	// Rel (required) specifies the relationship between the current document and the linked document.
	Rel RelationshipLink 

	// Sizes specifies the size of the linked resource. Only for rel="icon".
	Sizes	Sizes

	// Type specifies the media type of the linked document.
	Type string

	tmpl *template.Template

	str string
}

var zeroLink = Link{}

func (l *Link) validate() error {
	if l == nil {
		return nil
	}

	if *l == zeroLink {
		return fmt.Errorf("Link tag defined with no attributes set")
	}

	if l.Rel == "" {
		return fmt.Errorf("Link tag must include .Rel")
	}

	if !l.Sizes.isZero() && l.Rel != IconRL {
		return fmt.Errorf("Link tag with Sizes set must have Rel=IconRL")
	}

	return nil
}

func (l *Link) Attr() template.HTMLAttr {
	output := structToString(l)
	return template.HTMLAttr(output)
}

func (l *Link) isElement() {}

func (l *Link) compile() error {
	var err error
	l.tmpl, err = template.New("l").Parse(linkTmpl)
	if err != nil {
		return fmt.Errorf("Link object had error: %s", err)
	}

	return nil
}

func (l *Link) Execute(data interface{}) template.HTML {
	if l.str != "" {
		return template.HTML(l.str)
	}

	buff := strings.Builder{}

	if err := l.tmpl.Execute(&buff, pipeline{Self: l, Data: data}); err != nil {
		panic(err)
	}

	l.str = buff.String()
	return template.HTML(l.str)
}
