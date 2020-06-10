package html

import (
	"strings"
	"testing"
)

func TestP(t *testing.T) {
	tests := []struct {
		desc string
		p    *P
		want string
	}{
		{
			desc: "All attributes + 1 global + 1 event ",
			p: &P{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				Elements: []Element{TextElement("text")},
				Events:   (&Events{}).OnError("handleError"),
			},

			want: strings.TrimSpace(`
<p accesskey="key" onerror="handleError">
	text
</p>
`),
		},
	}

	for _, test := range tests {
		if err := test.p.Init(); err != nil {
			panic(err)
		}
		got := test.p.Execute(Pipeline{})
		if test.want != string(got) {
			t.Errorf("TestP(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
