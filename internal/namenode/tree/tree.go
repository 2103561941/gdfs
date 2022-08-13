package tree

import (
	"sync"
)

const (
	// File Type
	Direcotry  int = iota + 1 // 1
	NormalFile                // 2
)

// prefix tree
type Tree struct {
	rw   sync.RWMutex
	Root *Node // tree root '/'
}

// create a new Tree
func NewTree() *Tree {
	root := NewNode(
		SetFileName("/"),
		SetFilePath("/"),
		IsDirectory(true),
	)
	return &Tree{
		Root: root,
	}
}

// tree Node
// if file is directory, then the file keys is empty and filesize is zero.
// on the other hand, if file is normal file, than the children is empty.
type Node struct {
	FileName string
	FileType int    // directory or normal file
	FilePath string // filepath contains the filename and its parents' directory
	FileSize uint64
	FileKeys []string // file keys, it records the chunks' uuid of file

	Children []*Node
}

// in my use, i found that, if I use create the file by struct, I always forget to make a slice,
// so I chooose to use the function to init, and use options to fill this node.
func NewNode(opts ...Option) *Node {
	node := &Node{
		FileSize: 0,
		Children: make([]*Node, 0),
		FileKeys: make([]string, 0),
	}

	for _, opt := range opts {
		opt(node)
	}

	return node
}

type Option func(node *Node)

// if I pass true, then it will fill Directory to this node. on the other hand, it will fill NormalFile.
func IsDirectory(isDirectory bool) Option {
	return func(node *Node) {
		if isDirectory {
			node.FileType = Direcotry
			return
		}
		node.FileType = NormalFile
	}
}

func SetFileName(filename string) Option {
	return func(node *Node) {
		node.FileName = filename
	}
}

func SetFilePath(filepath string) Option {
	return func(node *Node) {
		node.FilePath = filepath
	}
}


func SetFileSize(filesize uint64) Option {
	return func(node *Node) {
		node.FileSize = filesize
	}
}