package html

import (
	"testing"
	"strings"
)

func TestComponent(t *testing.T) {
	tests := []struct {
		desc string
		component    *Component
		want string
	}{
		{
			desc: "All attributes + 1 global + 1 event ",
			component: &Component{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				TagType: "myComponent",
				TagValue: TextElement("value"),
				Events:         (&Events{}).OnError("handleError"),
			},

			want: strings.TrimSpace(`
<myComponent accesskey="key" onerror="handleError">
	value
</myComponent>
`),
		},
	}

	for _, test := range tests {
		if err := test.component.compile(); err != nil {
			panic(err)
		}
		got := test.component.Execute(struct{}{})
		if test.want != string(got) {
			t.Errorf("TestComponent(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
