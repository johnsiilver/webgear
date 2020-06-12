package html

import (
	"context"
	"net/url"
	"strings"
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
		got := &strings.Builder{}
		pipe := NewPipeline(context.Background(), nil, got)
		test.base.Execute(pipe)
		if test.want != got.String() {
			t.Errorf("TestBase(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
