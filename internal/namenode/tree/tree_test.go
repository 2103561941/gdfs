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
		NewNode(SetFileName("abc"), IsDirectory(true)),
		NewNode(SetFileName("abc"), IsDirectory(true)),
	}
}

// test search
func TestSearch(t *testing.T) {
	createTree()
	patterns := split("/abc")
	t.Logf("%+v, len: %d\n", patterns, len(patterns))
	node := search(patterns, tree.Root, 0)
	t.Logf("search node: %+v\n", node)
}

func TestMkdir(t *testing.T) {
	patterns := split("/tmp/cyb/node")
	dir := mkdir(patterns, tree.Root, 0)
	node := search(patterns, tree.Root, 0)
	if node != dir {
		t.Logf("mkdir exec error, node: %#v, dir: %#v", node, dir)
	}
}

func TestPut(t *testing.T) {
	filepath := "/tmp/cyb/node/1.txt"

	node := NewNode()
	tree.Put(filepath, node)
	s := search(split(filepath), tree.Root, 0)
	if s != node {
		t.Fatalf(
			`tree Put function's logical error:
				real: %#v
				expe: %#v`, s, node)
	}

	if err := tree.Put(filepath, node); err == nil {
		t.Fatalf("put twice failed not catch, please check the Put function")
	}
}
