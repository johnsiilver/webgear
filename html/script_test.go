package html

import (
	"testing"
	"net/url"
	"strings"
)

func TestScript(t *testing.T) {
	u, _ := url.Parse("/subpage")

	tests := []struct {
		desc string
		script    *Script
		want string
	}{
		{
			desc: "Empty",
			script: &Script{},
			want: "<script  >\n\t\n</script>",
		},
		{
			desc: "Everything + 1 global",
			script: &Script{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				Src: u,
				Type: "media",
				Async: true,
				Defer: true,
				TagValue: "javascript",
			},
			want: strings.TrimSpace(
`<script src="/subpage" type="media" async defer accesskey="key">
	javascript
</script>`),
		},
	}

	for _, test := range tests {
		if err := test.script.compile(); err != nil {
			panic(err)
		}
		got := test.script.Execute(struct{}{})
		if test.want != string(got) {
			t.Errorf("TestScript(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
