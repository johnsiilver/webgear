package index

import (
	"github.com/johnsiilver/webgear/html/examples/components/banner"

	. "github.com/johnsiilver/webgear/html"
)

const (
	bannerGearName = "banner-component"
)

// New creates a new Page object that can have the .Doc called to render the index page.
func New() (*Doc, error) {
	bannerGear, err := banner.New(bannerGearName)
	if err != nil {
		return nil, err
	}

	doc := &Doc{
		Head: &Head{
			Elements: []Element{
				&Meta{Charset: "UTF-8"},
				&Title{TagValue: TextElement("Go Language Basics")},
				&Link{Rel: "stylesheet", Href: URLParse("/static/index/index.css")},
				&Link{Href: URLParse("https://fonts.googleapis.com/css2?family=Share+Tech+Mono&display=swap"), Rel: "stylesheet"},
				&Link{Href: URLParse("https://fonts.googleapis.com/css2?family=Nanum+Gothic&display=swap"), Rel: "stylesheet"},
			},
		},
		Body: &Body{
			Elements: []Element{
				bannerGear, // This causes the code to render.
				&Component{GlobalAttrs: GlobalAttrs{ID: "banner"}, Gear: bannerGear},
			},
		},
	}

	return doc, nil
}
