package tree

import (
	"fmt"
	"time"

	"github.com/cyb0225/gdfs/internal/pkg/util"
)

const (
	// File Type
	Direcotry  int = iota + 1 // 1
	NormalFile                // 2
)

// tree Node
// If file is directory, then the file keys is empty and filesize is zero.
// On the other hand, if file is normal file, than the children is empty.
type Node struct {
	FileName   string
	FileType   int    // directory or normal file
	FilePath   string // filepath contains the filename and its parents' directory
	FileSize   int64
	FileKeys   []string // file keys, it records the chunks' uuid of file
	UpdateTime time.Time
	CreateTime time.Time

	Children []*Node
}

func NewNode(opts ...Option) *Node {
	node := &Node{
		FileSize:   0,
		Children:   make([]*Node, 0),
		FileKeys:   make([]string, 0),
		UpdateTime: time.Now(),
		CreateTime: time.Now(),
	}

	for _, opt := range opts {
		opt(node)
	}

	return node
}

type Option func(node *Node)

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

// Do some checks before append child file.
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
	n.UpdateTime = time.Now()
	return nil
}

// File chunk.
func (n *Node) CreateFileKeys(chunkSize int64) error {
	if n.FileType != NormalFile {
		return fmt.Errorf("file: %s is a directory", n.FilePath)
	}

	if n.FileSize == 0 {
		return fmt.Errorf("file: %s size is 0", n.FilePath)
	}

	//Rounded up
	num := int64(n.FileSize / chunkSize)
	if n.FileSize%chunkSize != 0 {
		num += 1
	}

	n.FileKeys = make([]string, num)
	for i := 0; i < len(n.FileKeys); i++ {
		uuid, err := util.GetUUID()
		if err != nil {
			n.FileKeys = make([]string, 0)
			return fmt.Errorf("creat file: %s keys: %w", n.FilePath, err)
		}
		n.FileKeys[i] = uuid
	}

	n.UpdateTime = time.Now()
	return nil
}