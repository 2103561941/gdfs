package tree

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
)

var (
	tree, _ = NewTree("./stroage/")
	chunckSize int64 = 1024 *1024
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

func TestPut(t *testing.T) {
	filepath := "/tmp/cyb/node/1.txt"

	node := NewNode(
		SetFilePath(filepath),
		IsDirectory(false),
	)
	if err := tree.Put(filepath, node); err != nil {
		t.Fatalf("check the file tree, file has not been push into the file tree, but it has been found: %s", err.Error())
	}

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

// test basic funciton.
// test if the file is existed.
// test if the file is a normalfile
func TestAppendChild(t *testing.T) {
	n := NewNode(
		IsDirectory(true),
		SetFileName("hello"),
		SetFilePath("/hello"),
	)

	f := NewNode(
		IsDirectory(false),
		SetFileName("cyb"),
	)

	if err := n.AppendChild(f); err != nil {
		t.Fatalf("append child error: %s", err.Error())
	}

	if err := n.AppendChild(f); err == nil {
		t.Fatalf("AppendChild have append the same file twice")
	}

	if err := f.AppendChild(n); err == nil {
		t.Fatalf("a normal file can append file")
	}
}

// test normal
// test directory file or empty file.
func TestCreateFileKeys(t *testing.T) {
	viper.Set("chunckSize", 1024)
	n := NewNode(
		IsDirectory(true),
		SetFileName("hello"),
		SetFilePath("/hello"),
	)

	if err := n.CreateFileKeys(chunckSize); err == nil {
		t.Fatalf("directory file can CreateFileKeys")
	}

	e := NewNode(
		IsDirectory(false),
		SetFileName("cyb"),
		SetFileSize(0),
	)

	if err := e.CreateFileKeys(chunckSize); err == nil {
		t.Fatalf("empty file can CreateFileKeys")
	}

	f := NewNode(
		IsDirectory(false),
		SetFileName("cyb"),
		SetFileSize(1024+1),
	)
	if err := f.CreateFileKeys(chunckSize); err != nil {
		t.Fatalf("normal file cannot CreateFileKeys: %s", err.Error())
	}
	fmt.Println(len(f.FileKeys))
	fmt.Println(f.FileKeys)
}

func TestDelete(t *testing.T) {

}
