package html_test

import (
	"context"
	"fmt"
	"testing"

	. "github.com/johnsiilver/webgear/html"
	"github.com/johnsiilver/webgear/html/builder"

	"github.com/kylelemons/godebug/pretty"
)

func TestWalker(t *testing.T) {
	categories := []string{"film1", "film2"}

	b := builder.NewHTML(&Head{}, &Body{GlobalAttrs: GlobalAttrs{ID: "body"}})
	b.Add(&Link{GlobalAttrs: GlobalAttrs{ID: "stylesheet"}, Rel: "stylesheet", Href: URLParse("/static/components/categories/categories.css")})

	b.Into(&Div{GlobalAttrs: GlobalAttrs{ID: "div"}})
	b.Into(&Form{GlobalAttrs: GlobalAttrs{ID: "form"}})
	b.Into(&Select{GlobalAttrs: GlobalAttrs{ID: "select"}, Name: "categories"})
	for _, category := range categories {
		// Adds our option with an event attachment when somone selects something.
		b.Add(&Option{GlobalAttrs: GlobalAttrs{ID: fmt.Sprintf("%s-option", category)}, Value: category, TagValue: category})
	}

	doc := b.Doc()

	want := []string{
		"body",
		"stylesheet",
		"div",
		"form",
		"select",
		"film1-option",
		"film2-option",
	}

	got := []string{}
	for element := range Walker(context.Background(), doc.Body) {
		got = append(got, GetElementID(element))
	}

	if diff := pretty.Compare(want, got); diff != "" {
		t.Errorf("TestWalker: -want/+got:\n%s", diff)
		t.Errorf(pretty.Sprint(got))
	}
}
