package html

import (
	"context"
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
		got := &strings.Builder{}
		pipe := NewPipeline(context.Background(), nil, got)
		test.p.Execute(pipe)
		if test.want != got.String() {
			t.Errorf("TestP(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
