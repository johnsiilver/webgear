package html

import (
	"context"
	"net/url"
	"strings"
	"testing"
)

func TestA(t *testing.T) {
	u, _ := url.Parse("/reg")

	tests := []struct {
		desc string
		a    *A
		want string
	}{
		{
			desc: "Empty attributes",
			a:    &A{},
			want: `<a   >
</a>`,
		},
		{
			desc: "All attributes + 1 global + 1 event",
			a: &A{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				Href:           "/subpage",
				Download:       true,
				HrefLang:       "english",
				Media:          "query",
				Ping:           u,
				ReferrerPolicy: Origin,
				Rel:            AuthorRel,
				Target:         BlankTarget,
				Type:           "media",
				Elements:       []Element{TextElement("text")},
				Events:         (&Events{}).OnError("handleError"),
			},
			want: `<a href="/subpage" download hreflang="english" media="query" ping="/reg" ` +
				`referrerpolicy="origin" rel="author" target="_blank" type="media" accesskey="key" ` +
				`onerror="handleError">
	text
</a>`,
		},
	}

	for _, test := range tests {
		got := &strings.Builder{}
		pipe := NewPipeline(context.Background(), nil, got)
		test.a.Execute(pipe)
		if test.want != got.String() {
			t.Errorf("TestA(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
