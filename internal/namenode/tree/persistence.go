package tree

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Directory tree persistence
type Persistence struct {
	fd          *os.File // tree.log
	storagePath string
	mu          sync.Mutex
}

func NewPersistence(storagePath string) (*Persistence, error) {
	fd, err := os.OpenFile(storagePath+"tree.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, fmt.Errorf("create logging file failed: %w", err)
	}

	per := &Persistence{
		fd:          fd,
		storagePath: storagePath,
	}
	return per, nil
}



func (p *Persistence) Put(node *Node) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	// write template:
	// put
	// filepath
	// filename
	// filesize
	// file keys (k1 k2 k3)
	// updatetime
	// createtime
	//
	w := bufio.NewWriter(p.fd)

	w.WriteString("put\n")
	w.WriteString(node.FilePath + "\n")
	w.WriteString(strconv.FormatInt(node.FileSize, 10) + "\n")

	keys := strings.Join(node.FileKeys, " ")
	w.WriteString(keys + "\n")
	w.WriteString(node.UpdateTime.Format("2006-01-02 15:04:05.000")+ "\n")
	w.WriteString(node.CreateTime.Format("2006-01-02 15:04:05.000") + "\n")
	w.WriteString("\n")
	return w.Flush()
}

func (p *Persistence) Delete(filepath string) error {
	p.mu.Lock()
	defer p.mu.Unlock() // write template:
	// delete
	// filepath
	//
	w := bufio.NewWriter(p.fd)
	w.WriteString("delete\n")
	w.WriteString(filepath + "\n")
	w.WriteString("\n")

	return w.Flush()
}

func (p *Persistence) Mkdir(node *Node) error {
	p.mu.Lock()
	defer p.mu.Unlock() // write template:
	// mkdir
	// filepath
	// updatetime
	// createtime
	//
	w := bufio.NewWriter(p.fd)
	w.WriteString("mkdir\n")
	w.WriteString(node.FilePath + "\n")
	w.WriteString(node.UpdateTime.Format("2006-01-02 15:04:05.000") + "\n")
	w.WriteString(node.UpdateTime.Format("2006-01-02 15:04:05.000") + "\n")
	w.WriteString("\n")
	return w.Flush()
}

func (p *Persistence) Rename(src string, dest string) error {
	p.mu.Lock()
	defer p.mu.Unlock() // write template:
	// rename
	// src
	// dest
	//
	w := bufio.NewWriter(p.fd)
	w.WriteString("rename\n")
	w.WriteString(src + "\n")
	w.WriteString(dest + "\n")
	w.WriteString("\n")
	return w.Flush()
}
