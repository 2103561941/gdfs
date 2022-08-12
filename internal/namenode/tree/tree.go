package tree

import "sync"

const (
	// File Type
	Direcotry  int = iota + 1 // 1
	NormalFile                // 2
)

// prefix tree
type Tree struct {
	rw sync.RWMutex
	Root *Node // tree root '/'
}

// create a new Tree
func NewTree() *Tree {
	root := &Node{
		Partten: "/",
		Children: make([]*Node, 0),
		FileMeta: &Meta{
			FileType: Direcotry,
		},
	}
	return &Tree{
		Root: root,
	}
}

// tree Node
type Node struct {
	Partten  string
	FileMeta *Meta
	Children []*Node
}


// file Meta data
type Meta struct {
	Filename string
	FileType int // directory or normal file
	Path     string
	FileSize uint64
	FileKey  string // string key
}