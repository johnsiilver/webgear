package html

import (
	"context"
	"strings"
	"testing"
)

func TestIFrame(t *testing.T) {
	tests := []struct {
		desc   string
		iframe *IFrame
		want   string
	}{
		{
			desc: "Most attributes using Src + two sandboxing params + 1 global + 1 event",
			iframe: &IFrame{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				Events:              (&Events{}).AddScript(OnError, "handleError"),
				Name:                "name",
				Src:                 URLParse("https://vimeo.com"),
				Allow:               "autoplay; fullscreen",
				AllowFullscreen:     true,
				AllowPaymentRequest: true,
				Height:              110,
				Width:               110,
				ReferrerPolicy:      OriginWhenCrossOrigin,
				Sandboxing:          Sandboxing{AllowFormsSB, AllowPopupsSB},
				Loading:             LazyILoad,
			},

			want: strings.TrimSpace(`
<iframe name="name" src="https://vimeo.com" allow="autoplay; fullscreen" allowfullscreen="true" ` +
				`allowpaymentrequest="true" height="110" width="110" referrerpolicy="origin-when-cross-origin" ` +
				`sandbox="allow-forms allow-popups" loading="lazy" accesskey="key" onerror="handleError"></iframe>
`),
		},
	}

	for _, test := range tests {
		got := &strings.Builder{}
		pipe := NewPipeline(context.Background(), nil, got)
		test.iframe.Execute(pipe)
		if test.want != got.String() {
			t.Errorf("TestIFrame(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
