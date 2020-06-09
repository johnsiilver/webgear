package html

import (
	"net/url"
	"testing"
)

func TestBase(t *testing.T) {
	u, _ := url.Parse("/subpage")

	tests := []struct {
		desc string
		base *Base
		want string
	}{
		{
			desc: "Everything + 1 global",
			base: &Base{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				Href:   u,
				Target: "_blank",
			},
			want: `<base href="/subpage" target="_blank" accesskey="key">`,
		},
	}

	for _, test := range tests {
		if err := test.base.compile(); err != nil {
			panic(err)
		}
		got := test.base.Execute(Pipeline{})
		if test.want != string(got) {
			t.Errorf("TestBase(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
