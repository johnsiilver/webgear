package html

import (
	"fmt"
	"html/template"
	"net/url"
	"strings"
)

var imgTmpl = template.Must(template.New("img").Parse(strings.TrimSpace(`
<img {{.Self.Attr}} {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}/>
`)))

// Img details an image to be shown.
type Img struct {
	GlobalAttrs

	Events *Events

	// Src specifies the path to the image.
	Src *url.URL

	// SrcSet specifies the path to an image to use in different situations.
	SrcSet *url.URL

	// Alt specifies an alternate text for the image, if the image for some reason cannot be displayed.
	Alt string

	// UseMap specifies an image as a client-side image-map.
	UseMap string

	// CrossOrigin allows images from third-party sites that allow cross-origin access to be used with canvas
	CrossOrigin CrossOrigin

	// HeightPx is the height set in pixels.
	HeightPx uint `html:"height" suffix:"px"`

	// HeightEm is the height set in em.
	HeightEm uint `html:"height" suffix:"em"`

	// WidthPx is the width set in pixels.
	WidthPx uint `html:"width" suffix:"px"`

	// WidthEm is the width set in pixels.
	WidthEm uint `html:"width" suffix:"em"`

	// IsMap specifies an image as a server-side image-map.
	IsMap bool `html:"attr"`

	// LongDesc specifies a URL to a detailed description of an image.
	LongDesc *url.URL

	// ReferrerPolicy specifies which referrer to send.
	ReferrerPolicy ReferrerPolicy

	Sizes string
}

func (i *Img) validate() error {
	// TODO(johnsiilver): could be done more simply.
	if i.HeightPx > 0 && i.HeightEm > 0 {
		return fmt.Errorf("Img tag cannot have HeightPx and HeightEm both set")
	}

	if i.WidthPx > 0 && i.WidthEm > 0 {
		return fmt.Errorf("Img tag cannot have WidthPx and WidthEM both set")
	}

	var hpx, hem, wpx, wem bool
	switch {
	case i.HeightPx > 0:
		hpx = true
	case i.HeightEm > 0:
		hem = true
	case i.WidthPx > 0:
		wpx = true
	case i.WidthEm > 0:
		wem = true
	}

	switch {
	case hpx && wem:
		return fmt.Errorf("Img tag cannot have HeightPx and WidthEm both set")
	case hem && wpx:
		return fmt.Errorf("Img tag cannot have HeightEm and WidthPx both set")
	}
	return nil
}

func (i *Img) Attr() template.HTMLAttr {
	output := structToString(i)
	return template.HTMLAttr(output)
}

func (i *Img) Execute(pipe Pipeline) string {
	pipe.Self = i

	if err := imgTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}
