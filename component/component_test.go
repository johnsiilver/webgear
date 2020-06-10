package component

import (
	"regexp"
	"strings"
	"testing"

	"github.com/johnsiilver/webgear/html"

	"github.com/kylelemons/godebug/pretty"
)

func TestComponent(t *testing.T) {
	tests := []struct {
		desc string
		name string
		doc  *html.Doc
		want string
		err  bool
	}{
		{
			desc: "Success",
			name: "my-component",
			doc: &html.Doc{
				Body: &html.Body{
					Elements: []html.Element{
						&html.Div{
							Elements: []html.Element{
								&html.A{Href: "/self", Elements: []html.Element{html.TextElement("link")}},
							},
						},
					},
				},
			},
			want: strings.TrimSpace(`
<template id=".Self.NameTemplate">
	<div  >
		<a href="/self"  >
	link
</a>
	</div>
</template>

<script>
	window.customElements.define(
		'.Self.Name',
		class extends HTMLElement {
			constructor() {
				super();
				let template = document.getElementById('.Self.Name');
				let templateContent = template.content;

				const shadowRoot = this.attachShadow({mode: 'open'}).appendChild(templateContent.cloneNode(true));
			}
		}
	);
</script>
			`),
		},
	}

	for _, test := range tests {
		g, err := New(test.name, test.doc, nil)
		switch {
		case err == nil && test.err:
			t.Errorf("TestComponent(%s): got err == nil, want err != nil", test.desc)
		case err != nil && !test.err:
			t.Errorf("TestComponent(%s): got err == %s, want err == nil", test.desc, err)
		case err != nil:
			continue
		}

		h, err := g.Execute(html.Pipeline{})
		if err != nil {
			t.Errorf("TestComponent(%s).Execute(): got err == %s, want err == nil", test.desc, err)
			continue
		}

		space := regexp.MustCompile(`\s+`)
		got := strings.TrimSpace(space.ReplaceAllString(string(h), " "))
		want := strings.TrimSpace(space.ReplaceAllString(string(test.want), " "))

		want = strings.ReplaceAll(want, ".Self.Name", g.Name())

		if diff := pretty.Compare(want, got); diff != "" {
			t.Errorf("TestComponent(%s): -want/+got:\n%s", test.desc, diff)
		}
	}
}
