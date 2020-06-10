package html

import (
	"html/template"
	"net/url"
	"strings"
	"sync"
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

var aTmpl = template.Must(template.New("a").Parse(strings.TrimSpace(`
<a {{.Self.Attr }} {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</a>
`)))

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

	Elements []Element

	Events *Events

	pool sync.Pool
}

func (a *A) Attr() template.HTMLAttr {
	output := structToString(a)
	return template.HTMLAttr(output)
}

func (a *A) Init() error {
	a.pool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}

	return nil
}

func (a *A) Execute(pipe Pipeline) template.HTML {
	buff := a.pool.Get().(*strings.Builder)
	defer a.pool.Put(buff)
	buff.Reset()

	pipe.Self = a

	if err := aTmpl.Execute(buff, pipe); err != nil {
		panic(err)
	}

	return template.HTML(buff.String())
}
