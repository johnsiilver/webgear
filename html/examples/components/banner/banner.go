package banner

import (
	"github.com/johnsiilver/webgear/component"
	"github.com/johnsiilver/webgear/html/builder"

	. "github.com/johnsiilver/webgear/html"
)

// New constructs a new component that shows a banner.
func New(name string, options ...component.Option) (*component.Gear, error) {
	build := builder.NewHTML(&Head{}, &Body{GlobalAttrs: GlobalAttrs{ID: "banner"}})
	build.Into(&Div{GlobalAttrs: GlobalAttrs{ID: "banner"}})
	build.Add(&Link{Rel: "stylesheet", Href: URLParse("/static/banner/banner.css")})

	build.Into(&A{Href: URLParse("/")})
	build.Add(&Img{GlobalAttrs: GlobalAttrs{ID: "gopher"}, Src: URLParse("/static/banner/scientist.svg")})
	build.Up()

	build.Into(&A{Href: URLParse("/")})
	build.Into(&Span{GlobalAttrs: GlobalAttrs{ID: "title"}})
	build.Add(TextElement("Example Banner"))

	gear, err := component.New(name, build.Doc(), options...)
	if err != nil {
		return nil, err
	}

	return gear, nil
}
