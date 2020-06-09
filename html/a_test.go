package html

import (
	"net/url"
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
			want: "<a   ></a>",
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
				TagValue:       TextElement("text"),
				Events:         (&Events{}).OnError("handleError"),
			},
			want: `<a href="/subpage" download hreflang="english" media="query" ping="/reg" ` +
				`referrerpolicy="origin" rel="author" target="_blank" type="media" accesskey="key" ` +
				`onerror="handleError">text</a>`,
		},
	}

	for _, test := range tests {
		if err := test.a.compile(); err != nil {
			panic(err)
		}
		got := test.a.Execute(Pipeline{})
		if test.want != string(got) {
			t.Errorf("TestA(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
