package html

import (
	"html/template"
	"testing"
)

func TestGlobalAttrs(t *testing.T) {
	tests := []struct {
		desc  string
		attrs GlobalAttrs
		want  template.HTMLAttr
	}{
		{
			desc: "Empty attributes",
		},
		{
			desc: "All attributes",
			attrs: GlobalAttrs{
				AccessKey:       "key",
				Class:           "class",
				ContentEditable: true,
				Dir:             RTLDir,
				Draggable:       true,
				Hidden:          true,
				ID:              "id",
				Lang:            "english",
				SpellCheck:      true,
				Style:           "style",
				TabIndex:        1,
				Title:           "title",
				Translate:       Yes,
			},
			want: `accesskey="key" class="class" contenteditable="true" dir="rtl" draggable="true" hidden id="id" ` +
				`lang="english" spellcheck="true" style="style" tabindex="1" title="title" translate="yes"`,
		},
	}

	for _, test := range tests {
		got := test.attrs.Attr()
		if test.want != got {
			t.Errorf("TestGlobalAttrs(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
