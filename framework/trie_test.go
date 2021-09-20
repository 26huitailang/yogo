package framework

import "testing"

func Test_filterChildNodes(t *testing.T) {
	root := &node{
		isLast:  false,
		segment: "",
		handler: func(c *Context) error { return nil },
		children: []*node{
			{
				isLast:   true,
				segment:  "FOO",
				handler:  func(c *Context) error { return nil },
				children: nil,
			},
			{
				isLast:   false,
				segment:  ":id",
				handler:  nil,
				children: nil,
			},
		},
	}

	{
		nodes := root.filterChildNodes("FOO")
		if len(nodes) != 2 {
			t.Error("foo error")
		}
	}

	{
		nodes := root.filterChildNodes(":foo")
		if len(nodes) != 2 {
			t.Error(":foo error")
		}
	}
}

func Test_matchNode(t *testing.T) {
	root := &node{
		isLast:  false,
		segment: "",
		handler: func(*Context) error { return nil },
		children: []*node{
			{
				isLast:  true,
				segment: "FOO",
				handler: nil,
				children: []*node{
					&node{
						isLast:   true,
						segment:  "BAR",
						handler:  func(*Context) error { panic("not implemented") },
						children: []*node{},
					},
				},
			},
			{
				isLast:   true,
				segment:  ":id",
				handler:  nil,
				children: nil,
			},
		},
	}

	{
		node := root.matchNode("foo/bar")
		if node == nil {
			t.Error("match normal node error")
		}
	}

	{
		node := root.matchNode("test")
		if node == nil {
			t.Error("match test")
		}
	}
}
