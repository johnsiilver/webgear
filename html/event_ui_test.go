// +build js,wasm

package html

import (
	"context"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestBody(t *testing.T) {
	a := &Body{
		GlobalAttrs: GlobalAttrs{ID: "body"},
		Elements: []Element{
			&Div{
				GlobalAttrs: GlobalAttrs{ID: "div"},
				Elements: []Element{
					&Span{
						GlobalAttrs: GlobalAttrs{ID: "span"},
					},
				},
			},
		},
	}

	want := &Body{
		GlobalAttrs: GlobalAttrs{ID: "myBody"},
		Elements: []Element{
			&Div{
				GlobalAttrs: GlobalAttrs{ID: "myDiv"},
				Elements: []Element{
					&Span{
						GlobalAttrs: GlobalAttrs{ID: "mySpan"},
					},
				},
			},
		},
	}

	ui := newUI(a, &Wasm{})

	b := ui.Body()

	b.GlobalAttrs.ID = "myBody"
	b.Elements[0].(*Div).GlobalAttrs.ID = "myDiv"
	b.Elements[0].(*Div).Elements[0].(*Span).GlobalAttrs.ID = "mySpan"

	if a == b {
		t.Fatalf("TestBody: our deep copy has the same top level pointer, that is a complete fail")
	}
	switch {
	case a.ID != "body":
		t.Errorf("TestBody: a.ID was %q, want %q", ui.body.ID, "body")
	case a.Elements[0].(*Div).ID != "div":
		t.Errorf("TestBody: a.Elements[0].ID was %q, want %q", a.Elements[0].(*Div).ID, "div")
	case a.Elements[0].(*Div).Elements[0].(*Span).ID != "span":
		t.Errorf("TestBody: a.Elements[0].(*Div)[0].(*Span).ID was %q, want %q", a.Elements[0].(*Div).Elements[0].(*Span).ID, "span")
	}

	if diff := pretty.Compare(want, ui.body); diff != "" {
		t.Errorf("TestBody: while the outside representation did not change(good), the internal representation did not change(bad): -want/+got:\n%s", diff)
	}
}

func TestUpdate(t *testing.T) {
	body := &Body{
		Elements: []Element{
			&Div{
				GlobalAttrs: GlobalAttrs{ID: "myDiv"},
				Elements: []Element{
					&Ul{
						GlobalAttrs: GlobalAttrs{ID: "list"},
						Elements: []Element{
							&Li{
								Elements: []Element{
									TextElement("Hello"),
								},
							},
						},
					},
				},
			},
		},
	}

	tests := []struct {
		desc    string
		id      string
		element Element
		err     bool
	}{
		{
			desc: "Successful replace of list",
			id:   "list",
			element: &Ul{
				GlobalAttrs: GlobalAttrs{ID: "list2"},
				Elements: []Element{
					&Li{
						Elements: []Element{
							TextElement("World"),
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		wasm := New(&Doc{Body: copyBody(body)})
		go wasm.Run(context.Background())

		wasm.Ready()

		ui := wasm.UI()

		err := ui.Update(test.id, test.element)
		switch {
		case err == nil && test.err:
			t.Errorf("TestUpdate(%s): got err == nil, want != nil", test.desc)
			continue
		case err != nil && !test.err:
			t.Errorf("TestUpdate(%s): got err == %s, want == nil", test.desc, err)
			continue
		case err != nil:
			continue
		}

		if err := ui.Close(); err != nil {
			t.Errorf("TestUpdate(%s): error on Close(): %s", test.desc, err)
			continue
		}
	}
}
