package html

import (
	"fmt"
	"html/template"
	"net/url"
	"strings"
)

var iframeTmpl = template.Must(template.New("iframe").Parse(strings.TrimSpace(`
<iframe {{.Self.Attr}} {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}></iframe>
`)))

// Sandboxing reprsents all the Sandbox attributes to apply.
type Sandboxing []Sandbox

func (s Sandboxing) isRaw() {}
func (s Sandboxing) String() string {
	if len(s) == 0 {
		return ""
	}

	buff := &strings.Builder{}
	for _, item := range s {
		if item == "all" {
			return "sandbox"
		}
		buff.WriteString(string(item + " "))
	}

	return fmt.Sprintf("sandbox=%q", strings.TrimSpace(buff.String()))
}

// Sandbox reprsents the HTML sandbox attribute commonly used on iframes.
type Sandbox string

const (
	AllSB              Sandbox = "all"
	AllowFormsSB       Sandbox = "allow-forms"
	AllowPointerLockSB Sandbox = "allow-pointer-lock"
	AllowPopupsSB      Sandbox = "allow-popups"
	AllowSameOriginSB  Sandbox = "allow-same-origin"
	AllowScriptsSB     Sandbox = "allow-scripts"
	AllowTopNavigation Sandbox = "allow-top-navigation"
)

// IFrameLoad indicates when to load an IFrame's content.
type IFrameLoad string

const (
	// EagerILoad loads the IFrame even if it isn't visible yet on the screen. This is the default.
	EagerILoad IFrameLoad = "eager"
	// LazyILoad loads the IFrame only when it becomes visible on the screen.
	LazyILoad IFrameLoad = "lazy"
)

// IFrame represents a division tag.
type IFrame struct {
	GlobalAttrs

	// Name specifies the name of the <iframe>.
	Name string
	// Src specifies the address of the document to embed in the <iframe>.
	Src *url.URL
	// SrcDoc specifies the HTML content of the page to show in the <iframe>.
	SrcDoc template.HTMLAttr
	// Allow specifies a feature policy for the <iframe>.
	Allow string
	// Allow if set to true allows activating fullscreen mode by calling the requestFullscreen() method.
	AllowFullscreen bool
	// AllowPaymentRequest if set to true if will allow a cross-origin <iframe> to invoke the Payment Request API.
	AllowPaymentRequest bool
	// Height specifies the height of an <iframe>. Default height is 150 pixels.
	Height uint
	// Width specifies the width of an <iframe>. Default width is 300 pixels.
	Width uint
	// ReferrerPolicy specifies how much/which referrer information that will be sent when processing the iframe attributes.
	ReferrerPolicy ReferrerPolicy
	// Sandboxing enables an extra set of restrictions for the content in an <iframe>.
	Sandboxing Sandboxing
	// Loading indicates the way the browser loads the iframe (immediately or when on the visible screen).
	Loading IFrameLoad

	Events *Events
}

func (i *IFrame) validate() error {
	if i.Src != nil && i.SrcDoc != "" {
		return fmt.Errorf("cannot have Src and SrcDoc both set")
	}
	return nil
}

func (i *IFrame) Attr() template.HTMLAttr {
	output := structToString(i)
	return template.HTMLAttr(output)
}

func (i *IFrame) Execute(pipe Pipeline) string {
	pipe.Self = i

	if err := iframeTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}
