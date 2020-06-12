package html

import (
	"context"
	"net/url"
	"strings"
	"testing"
)

func TestLink(t *testing.T) {

	u, _ := url.Parse("/subpage")

	tests := []struct {
		desc string
		link *Link
		want string
	}{

		{
			desc: "Rel + 1 global",
			link: &Link{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				Rel: "UTF-8",
			},
			want: `<link rel="UTF-8" accesskey="key">`,
		},
		{
			desc: "Everything + 1 global",
			link: &Link{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				Href:           u,
				CrossOrigin:    AnnonymousCO,
				HrefLang:       "english",
				Media:          "media",
				ReferrerPolicy: Origin,
				Rel:            AlternateRL,
				Sizes:          Sizes{Width: 5, Height: 10},
				Type:           "type",
			},
			want: `<link href="/subpage" crossorigin="anonymous" hreflang="english" media="media" ` +
				`referrerpolicy="origin" rel="alternate" sizes="10x5" type="type" accesskey="key">`,
		},
	}

	for _, test := range tests {
		got := &strings.Builder{}
		pipe := NewPipeline(context.Background(), nil, got)
		test.link.Execute(pipe)
		if test.want != got.String() {
			t.Errorf("TestLink(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
