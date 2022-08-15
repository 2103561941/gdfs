package tree

import (
	"fmt"
	"sync"

	"github.com/cyb0225/gdfs/internal/pkg/util"
	"github.com/spf13/viper"
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
	FileSize int64
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

func SetFileSize(filesize int64) Option {
	return func(node *Node) {
		node.FileSize = filesize
	}
}

func (n *Node) IsDirectory() bool {
	return n.FileType == Direcotry
}

// do some checks before append child file.
// here, I not check the filepath, it will check by tree Put funciton.
// because, the file's relation belongs to the tree, not a node.
func (n *Node) AppendChild(node *Node) error {
	if n.FileType != Direcotry {
		return fmt.Errorf("file: %s is not a directory", n.FilePath)
	}

	for _, child := range n.Children {
		if child.FileName == node.FileName {
			return fmt.Errorf("file: %s is already exist in %s", node.FileName, n.FilePath)
		}
	}

	n.Children = append(n.Children, node)
	return nil
}

// depends on the size of file stored in node.
func (n *Node) CreateFileKeys() error {
	if n.FileType != NormalFile {
		return fmt.Errorf("file: %s is a directory", n.FilePath)
	}

	if n.FileSize == 0 {
		return fmt.Errorf("file: %s size is 0", n.FilePath)
	}

	//Rounded up
	size := viper.GetInt64("chunckSize")
	// fmt.Println("size: ", size)
	// fmt.Println("filesize: ", n.FileSize)
	num := int64(n.FileSize / size)
	if n.FileSize%size != 0 {
		num += 1
	}
	// fmt.Println("num: ", num)

	n.FileKeys = make([]string, num)
	for i := 0; i < len(n.FileKeys); i++ {
		uuid, err := util.GetUUID()
		if err != nil {
			n.FileKeys = make([]string, 0)
			return fmt.Errorf("creat file: %s keys: %w", n.FilePath, err)
		}
		n.FileKeys[i] = uuid
	}
	return nil
}
