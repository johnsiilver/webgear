package html

import (
	"strings"
	"testing"
)

func TestTitle(t *testing.T) {
	tests := []struct {
		desc  string
		title *Title
		want  string
	}{
		{
			desc:  "Empty attributes",
			title: &Title{},
			want:  "<title ></title>",
		},
		{
			desc: "All attributes + 1 global",
			title: &Title{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				TagValue: TextElement("hello"),
			},

			want: strings.TrimSpace(`
<title accesskey="key">hello</title>
`),
		},
	}

	for _, test := range tests {
		if err := test.title.Init(); err != nil {
			panic(err)
		}
		got := test.title.Execute(Pipeline{})
		if test.want != string(got) {
			t.Errorf("TestTitle(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
