package banner

import (
	"github.com/johnsiilver/webgear/component"
	"github.com/johnsiilver/webgear/html"
)

// New constructs a new component that shows a banner.
func New(name string, options ...component.Option) (*component.Gear, error) {
	doc := &html.Doc{
		Body: &html.Body{
			Elements: []html.Element{
				&html.Div{
					GlobalAttrs: html.GlobalAttrs{ID: "banner"},
					Elements: []html.Element{
						&html.Link{Rel: "stylesheet", Href: html.URLParse("/static/apps/snippets/components/banner/banner.css")},
						&html.A{
							Href: "/",
							Elements: []html.Element{
								&html.Img{
									GlobalAttrs: html.GlobalAttrs{ID: "gopher"},
									Src:         html.URLParse("/static/apps/snippets/components/banner/surfing-js.svg"),
								},
							},
						},
						&html.A{
							Href: "/",
							Elements: []html.Element{
								&html.Span{
									GlobalAttrs: html.GlobalAttrs{ID: "title"},
									Elements: []html.Element{
										html.TextElement("Webgear Snippets"),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	options = append(options)

	gear, err := component.New(name, doc, options...)
	if err != nil {
		return nil, err
	}

	return gear, nil
}