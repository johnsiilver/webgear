package component

import (
	"testing"
	"strings"

	"github.com/johnsiilver/webgear/html"

	"github.com/yosssi/gohtml"
	"github.com/kylelemons/godebug/pretty"
)

func TestComponent(t *testing.T) {
	tests := []struct{
		desc string
		name string
		doc *html.Doc
		want string
		err bool
	}{
		{
			desc: "Success",
			name: "myComponent",
			doc: &html.Doc{
				Body: &html.Body{
					Elements: []html.Element{
						&html.Div{
							Elements: []html.Element{
								&html.A{Href:"/self", TagValue: html.TextElement("link")},
							},
						},
					},
				},
			},
			want: strings.TrimSpace(`
<template id=".Self.NameTemplate">
	<div  >
		<a href="/self"  >link</a>
	</div>
</template>

<script>
	window.customElements.define('.Self.Name',
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
		g, err := NewGear(test.name, test.doc)
		switch {
		case err == nil && test.err:
			t.Errorf("TestComponent(%s): got err == nil, want err != nil", test.desc)
		case err != nil && !test.err:
			t.Errorf("TestComponent(%s): got err == %s, want err == nil", test.desc, err)
		case err != nil: 
			continue
		}

		got, err := g.Execute(struct{}{})
		if err != nil {
			t.Errorf("TestComponent(%s).Execute(): got err == %s, want err == nil", test.desc, err)
			continue
		}
		gotStr := gohtml.Format(string(got))

		want := strings.ReplaceAll(test.want, ".Self.Name", g.Name())
		want = gohtml.Format(want)

		if diff := pretty.Compare(want, gotStr); diff != "" {
			t.Errorf("TestComponent(%s): -want/+got:\n%s", test.desc, diff)
		}
	}
}