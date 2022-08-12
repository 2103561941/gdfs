package tree

import (
	"testing"
)

var (
	tree = NewTree()
)

// create a mock tree
func createTree() {
	tree.Root.Children = []*Node{
		&Node{Partten: "abc", FileMeta: &Meta{
			FileType: Direcotry,
		}},
		&Node{Partten: "abc", FileMeta: &Meta{
			FileType: NormalFile,
		}},
	}
}

// test search
func TestSearch(t *testing.T) {
	createTree()
	patterns := split("/abc")
	t.Logf("%+v, len: %d\n", patterns, len(patterns))
	node := search(patterns, tree.Root)
	t.Logf("search node: %+v\n", node)
}

func TestMkdir(t *testing.T) {
	patterns := split("/tmp/cyb/node")
	dir := mkdir(patterns, tree.Root)
	node := search(patterns, tree.Root)
	if node != dir {
		t.Logf("mkdir exec error, node: %#v, dir: %#v", node, dir)
	}
}

func TestPut(t *testing.T) {
	filepath := "/tmp/cyb/node/1.txt"
	node := &Node{
		Children: make([]*Node, 0),
		FileMeta: &Meta{
			FileType: NormalFile,
		},
	}
	tree.Put(filepath, node)
	s := search(split(filepath), tree.Root)
	if s != node {
		t.Fatalf(
			`tree Put function's logical error:
				real: %#v
				expe: %#v`, s, node)
	}
}
