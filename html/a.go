package html

import (
	"fmt"
	"html/template"
	"net/url"
	"strings"
)

// LanguageCode is use to specify the language in use.
type LanguageCode string

// MediaQuery provides a query to what media is being used.
type MediaQuery string

// MediaType is the type of media to use.
type MediaType string

// ReferrerPolicy specifies which referrer to send.
type ReferrerPolicy string

const (
	NoReferrer              ReferrerPolicy = "no-referrer"
	NoReferrerWhenDowngrade ReferrerPolicy = "no-referrer-when-downgrade"
	Origin                  ReferrerPolicy = "origin"
	OriginWhenCrossOrigin   ReferrerPolicy = "origin-when-cross-origin"
	UnsafeUrl               ReferrerPolicy = "unsafe-url"
)

type Relationship string

const (
	AlternateRel  Relationship = "alternate"
	AuthorRel     Relationship = "author"
	BookmarkRel   Relationship = "bookmark"
	ExternalRel   Relationship = "external"
	HelpRel       Relationship = "help"
	LicenseRel    Relationship = "license"
	NextRel       Relationship = "next"
	NoFollowRel   Relationship = "nofollow"
	NoReferrerRel Relationship = "noreferrer"
	NoOpenerRel   Relationship = "noopener"
	PrevRel       Relationship = "prev"
	SearchRel     Relationship = "search"
	TagRel        Relationship = "tag"
)

// Target specifies where to open the linked document.
type Target string

const (
	BlankTarget  = "_blank"
	ParentTarget = "_parent"
	SelfTarget   = "_self"
	TopTarget    = "_top"
)

var aTmpl = strings.TrimSpace(`
<a {{.Self.Attr }} {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>{{.Self.TagValue}}</a>
`)

// A defines a hyperlink, which is used to link from one page to another.
type A struct {
	// Href specifies the URL of the page the link goes to.
	Href string
	// Download specifies that the target will be downloaded when a user clicks on the hyperlink.
	Download bool `html:"attr"`
	// HrefLang specifies the language of the linked document.
	HrefLang LanguageCode
	// Media specifies what media/device the linked document is optimized for.
	Media MediaQuery
	// Ping specifies a space-separated list of URLs to which, when the link is followed, post requests with the body ping will be sent by the browser (in the background). Typically used for tracking.
	Ping *url.URL
	// ReferrerPolicy specifies which referrer to send.
	ReferrerPolicy ReferrerPolicy
	// Rel specifies the relationship between the current document and the linked document.
	Rel Relationship
	// Target specifies where to open the linked document.
	Target Target
	// Type specifies the media type of the linked document.
	Type MediaType

	GlobalAttrs

	// TagValue provides the value inside a reference. This can be another element such as a div, but most commonly
	// it is TextElement.
	TagValue Element

	Events *Events

	tmpl *template.Template

	str string
}

func (a *A) Attr() template.HTMLAttr {
	output := structToString(a)
	return template.HTMLAttr(output)
}

func (a *A) isElement() {}

func (a *A) compile() error {
	var err error
	a.tmpl, err = template.New("a").Parse(aTmpl)
	if err != nil {
		return fmt.Errorf("A object had error: %s", err)
	}

	return nil
}

func (a *A) Execute(data interface{}) template.HTML {
	if a.str != "" {
		return template.HTML(a.str)
	}

	buff := strings.Builder{}

	if err := a.tmpl.Execute(&buff, pipeline{Self: a, Data: data}); err != nil {
		panic(err)
	}

	a.str = buff.String()
	return template.HTML(a.str)
}
