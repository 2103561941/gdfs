package tree

import (
	"fmt"
	"os"
	"sync"

	"github.com/cyb0225/gdfs/pkg/log"
	"github.com/huandu/go-clone"
)


// prefix tree
type Tree struct {
	rw   sync.RWMutex
	Root *Node // tree root '/'
	Per  *Persistence
}

// create a new Tree
func NewTree(storagePath string) (*Tree, error) {
	root := NewNode(
		SetFileName("/"),
		SetFilePath("/"),
		IsDirectory(true),
	)

	per, err := NewPersistence(storagePath)
	if err != nil {
		return nil, fmt.Errorf("new persistence failed: %w", err)
	}

	tree := &Tree{
		Root:      root,
		Per:       per,
	}

	return tree, nil
}

// file tree persistence
func (t *Tree) Persistence() {
	// get copy tree.
	t.rw.RLock()
	tree := clone.Clone(t.Root).(*Node)

	// start record new logging file and close old file's fd
	fd, err := os.OpenFile(t.Per.storagePath+"tree.log.new", os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil { // if create file failed, then skip this backup.
		log.Error("create tree.log file failed", log.Err(err))
		return
	}
	t.Per.ChangeFD(fd)
	t.rw.RUnlock()

	t.Per.BackupTree(tree)
}
