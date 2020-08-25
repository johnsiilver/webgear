package wasm

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"

	. "github.com/johnsiilver/webgear/html"
)

func TestGetElementID(t *testing.T) {
	tests := []struct{
		e Element
		want string
	}{
		{&Span{GlobalAttrs: GlobalAttrs{ID: "mySpan"}}, "mySpan"},
		{TextElement("hello"), ""},
	}

	for _, test := range tests {
		got := getElementID(test.e)

		if test.want != got {
			t.Errorf("TestGetElementID(%+v): got %q, want %q", test.e, got, test.want)
		}
	}
}

func TestMapIDs(t *testing.T) {
	body := &Body{
		GlobalAttrs: GlobalAttrs{ID: "myBody"},
		Elements: []Element{
			&Ul{
				GlobalAttrs: GlobalAttrs{ID: "myUl"},
				Elements: []Element{
					&Li{
						GlobalAttrs: GlobalAttrs{ID: "myLi"},
						Elements: []Element{
							TextElement("hello"),
						},
					},
				},
			},
			&Ul{
				Elements: []Element{
					&Li{
						GlobalAttrs: GlobalAttrs{ID: "myLi1"},
						Elements: []Element{
							TextElement("hello"),
						},
					},
				},
			},
		},
	}

	want := map[string]elemNode{
		"myUl": elemNode{
			parent: body,
			element: body.Elements[0],
		},
		"myLi":elemNode{
			parent: body.Elements[0],
			element: body.Elements[0].(*Ul).Elements[0],
		},
		"myLi1":elemNode{
			parent: body.Elements[1],
			element: body.Elements[1].(*Ul).Elements[0],
		},
	}

	got:= map[string]elemNode{}
	if err := mapIDs(body, body.Elements, got); err != nil {
		t.Fatalf("TestMapIDs: fatal error: %s", err)
	}

	if diff := pretty.Compare(want, got); diff != "" {
		t.Fatalf("TestMapIDs: -want/+got:\n%s", diff)
	}
}

func TestReplaceElementInNode(t *testing.T) {
	tests := []struct{
		desc string
		parent Element
		element Element
		want Element
		err bool
	}{
		{
			desc: "Parent does not contain element",
			parent: &Div{},
			element: &Span{GlobalAttrs: GlobalAttrs{ID: "mySpan"}},
			err: true,
		},
		{
			desc: "Element does not have an ID",
			parent: &Div{
				Elements: []Element{
					&Span{GlobalAttrs: GlobalAttrs{ID: "mySpan"}},
				},
			},
			element: &Span{},
			err: true,
		},
		{
			desc: "Replaced element",
			parent: &Div{
				GlobalAttrs: GlobalAttrs{ID: "myDiv"},
				Elements: []Element{
					&Span{GlobalAttrs: GlobalAttrs{ID: "myData"}},
				},
			},
			element: &Div{GlobalAttrs: GlobalAttrs{ID: "myData"}},
			want: &Div{
				GlobalAttrs: GlobalAttrs{ID: "myDiv"},
				Elements: []Element{
					&Div{GlobalAttrs: GlobalAttrs{ID: "myData"}},
				},
			},
		},
	}

	for _, test := range tests{
		err := replaceElementInNode(test.parent, test.element)
		switch {
		case err == nil && test.err:
			t.Errorf("TestReplaceElementInNode(%s): got err == nil, want != nil", test.desc)
			continue
		case err != nil && !test.err:
			t.Errorf("TestReplaceElementInNode(%s): got err == %s, want == nil", test.desc, err)
			continue
		case err != nil:
			continue
		}

		if diff := pretty.Compare(test.want, test.parent); diff != "" {
			t.Errorf("TestReplaceElementInNode(%s): -want/+got:\n%s", test.desc, diff)
		}
	}
}

func TestAddElementToNode(t *testing.T) {
	tests := []struct{
		desc string
		node Element
		element Element
		want Element
		err bool
	}{
		{
			desc: "Does does not contain Elements or Element",
			node: TextElement("hello"),
			element: &Span{GlobalAttrs: GlobalAttrs{ID: "mySpan"}},
			err: true,
		},
		{
			desc: "Add element",
			node: &Div{
				GlobalAttrs: GlobalAttrs{ID: "myDiv"},
				Elements: []Element{
					&Span{GlobalAttrs: GlobalAttrs{ID: "mySpan"}},
				},
			},
			element: &Div{GlobalAttrs: GlobalAttrs{ID: "myData"}},
			want: &Div{
				GlobalAttrs: GlobalAttrs{ID: "myDiv"},
				Elements: []Element{
					&Span{GlobalAttrs: GlobalAttrs{ID: "mySpan"}},
					&Div{GlobalAttrs: GlobalAttrs{ID: "myData"}},
				},
			},
		},
	}

	for _, test := range tests{
		err := addElementToNode(test.node, test.element)
		switch {
		case err == nil && test.err:
			t.Errorf("TestAddElementToNode(%s): got err == nil, want != nil", test.desc)
			continue
		case err != nil && !test.err:
			t.Errorf("TestAddElementToNode(%s): got err == %s, want == nil", test.desc, err)
			continue
		case err != nil:
			continue
		}

		if diff := pretty.Compare(test.want, test.node); diff != "" {
			t.Errorf("TestAddElementToNode(%s): -want/+got:\n%s", test.desc, diff)
		}
	}
}

func TestDeleteNode(t *testing.T) {
	tests := []struct{
		desc string
		id string
		m map[string]elemNode
		wantMap map[string]elemNode
		wantParent *Div
		err bool
	}{
		{
			desc: "Does does not contain Element",
			id: "myID",
			m: map[string]elemNode{},
			err: true,
		},
		{
			desc: "Parent does not have an ID",
			id: "mySpan",
			m: map[string]elemNode{
				"mySpan": {
					parent: &Div{ 
						Elements: []Element{
							&Span{GlobalAttrs: GlobalAttrs{ID: "mySpan"}},
						},
					},
					element: &Span{GlobalAttrs: GlobalAttrs{ID: "mySpan"}},
				},
			},
			err: true,
		},
		{
			desc: "Delete element",
			id: "mySpan",
			m: map[string]elemNode{
				"mySpan": {
					parent: &Div{ 
						GlobalAttrs: GlobalAttrs{ID: "myDiv"},
						Elements: []Element{
							&Span{GlobalAttrs: GlobalAttrs{ID: "mySpan"}},
						},
					},
					element: &Span{GlobalAttrs: GlobalAttrs{ID: "mySpan"}},
				},
			},
			wantMap: map[string]elemNode{},
			wantParent: &Div{ 
				GlobalAttrs: GlobalAttrs{ID: "myDiv"},
				Elements: []Element{},
			},
		},
	}

	for _, test := range tests{
		parent, err := deleteElement(test.id, test.m)
		switch {
		case err == nil && test.err:
			t.Errorf("TestDeleteNode(%s): got err == nil, want != nil", test.desc)
			continue
		case err != nil && !test.err:
			t.Errorf("TestDeleteNode(%s): got err == %s, want == nil", test.desc, err)
			continue
		case err != nil:
			continue
		}

		if diff := pretty.Compare(test.wantMap, test.m); diff != "" {
			t.Errorf("TestDeleteNode(%s): map:  -want/+got:\n%s", test.desc, diff)
		}
		if diff := pretty.Compare(test.wantParent, parent.(*Div)); diff != "" {
			t.Errorf("TestDeleteNode(%s): parent:  -want/+got:\n%s", test.desc, diff)
		}
	}
}