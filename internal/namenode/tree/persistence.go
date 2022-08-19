package tree

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/cyb0225/gdfs/internal/namenode/config"
	"github.com/cyb0225/gdfs/pkg/log"
)

type Persistence struct {
	fd          *os.File // tree.log
	storagePath string
	mu          sync.Mutex
}

func NewPersistence(storagePath string) (*Persistence, error) {
	fd, err := os.OpenFile(storagePath+"tree.log", os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("create logging file failed: %w", err)
	}

	per := &Persistence{
		fd:          fd,
		storagePath: storagePath,
	}

	return per, nil
}

func (p *Persistence) ChangeFD(fd *os.File) {
	p.fd.Close()
	p.fd = fd
}

// write tree struct to file sysytem from memory
// read file system to memory
func (p *Persistence) BackupTree(node *Node) error {

	// start write memory into the bakcup file.
	_, err := os.OpenFile(p.storagePath+"tree.backup.new", os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("create tree.backup.new failed: %w", err)
	}

	// copy tree success, then remove old tree file and old logging file. and change file name.
	if err := os.Remove(config.Cfg.StoragePath + "tree.backup"); err != nil {
		return fmt.Errorf("remove tree.backup failed: %w", err)
	}
	if err := os.Remove(config.Cfg.StoragePath + "tree.log"); err != nil {
		return fmt.Errorf("remove tree.log failed: %w", err)
	}
	if err := os.Rename(config.Cfg.StoragePath+"tree.backup.new", config.Cfg.StoragePath+"tree.backup"); err != nil {
		return fmt.Errorf("rename tree.backup.new failed: %w", err)
	}
	if err := os.Rename(config.Cfg.StoragePath+"tree.log.new", config.Cfg.StoragePath+"tree.log"); err != nil {
		return fmt.Errorf("rename tree.log.new failed: %w", err)
	}

	return nil
}

// read bakcup file when restart namenode server.
func (p *Persistence) ReadBakcup() {
	// read tree file.

	// read log file.

	// read log.new file

	// put filekey into

}

// write to log.
func (p *Persistence) ReadLog() error {
	r := bufio.NewScanner(p.fd)

	for r.Scan() {
		line := r.Text()
		switch line {
		case "mkdir":

		case "delete":

		case "rename":

		case "put":

		default:
			log.Error("unknow method", log.String("method", line))
			return fmt.Errorf("unknow method")
		}
	}
	return nil
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
	if _, err := p.fd.WriteString("put\n"); err != nil {
		return err
	}
	if _, err := p.fd.WriteString(node.FilePath + "\n"); err != nil {
		return err
	}
	if _, err := p.fd.WriteString(node.FileName + "\n"); err != nil {
		return err
	}
	if _, err := p.fd.WriteString(strconv.FormatInt(node.FileSize, 10) + "\n"); err != nil {
		return err
	}
	keys := strings.Join(node.FileKeys, " ")
	if _, err := p.fd.WriteString(keys + "\n"); err != nil {
		return err
	}
	if _, err := p.fd.WriteString(node.UpdateTime.String() + "\n"); err != nil {
		return err
	}
	if _, err := p.fd.WriteString(node.CreateTime.String() + "\n"); err != nil {
		return err
	}
	if _, err := p.fd.WriteString("\n"); err != nil {
		return err
	}

	return nil
}

func putRead(r *bufio.Scanner) {
	r.Scan()
	// filepath := r.Text()
	// r.Scan()
	// filename := r.Text()
	// r.Scan()
	// filesize := r.Text()
	// r.Scan()
	// filekeys := r.Text()
	// r.Scan()
	// updatetime := r.Text()
	// r.Scan()
	// createtime := r.Text()
	// r.Scan()

	

}

func (p *Persistence) Delete(filepath string) error {
	p.mu.Lock()
	defer p.mu.Unlock() // write template:
	// delete
	// filepath
	//
	if _, err := p.fd.WriteString("delete\n"); err != nil {
		return err
	}
	if _, err := p.fd.WriteString(filepath + "\n"); err != nil {
		return err
	}
	if _, err := p.fd.WriteString("\n"); err != nil {
		return err
	}
	return nil
}

func deleteRead(r *bufio.Scanner, t *Tree) {

}

func (p *Persistence) Mkdir(node *Node) error {
	p.mu.Lock()
	defer p.mu.Unlock() // write template:
	// mkdir
	// filepath
	// filename
	// updatetime
	// createtime
	//
	if _, err := p.fd.WriteString("mkdir\n"); err != nil {
		return err
	}
	if _, err := p.fd.WriteString(node.FilePath + "\n"); err != nil {
		return err
	}
	if _, err := p.fd.WriteString(node.FileName + "\n"); err != nil {
		return err
	}
	if _, err := p.fd.WriteString(node.UpdateTime.String() + "\n"); err != nil {
		return err
	}
	if _, err := p.fd.WriteString(node.CreateTime.String() + "\n"); err != nil {
		return err
	}
	if _, err := p.fd.WriteString("\n"); err != nil {
		return err
	}
	return nil
}

func (p *Persistence) Rename(src string, dest string) error {
	p.mu.Lock()
	defer p.mu.Unlock() // write template:
	// rename
	// src
	// dest
	//

	if _, err := p.fd.WriteString("rename\n"); err != nil {
		return err
	}
	if _, err := p.fd.WriteString(src + "\n"); err != nil {
		return err
	}
	if _, err := p.fd.WriteString(dest + "\n"); err != nil {
		return err
	}
	if _, err := p.fd.WriteString("\n"); err != nil {
		return err
	}
	return nil
}
