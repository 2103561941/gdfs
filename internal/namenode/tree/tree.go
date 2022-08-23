package tree

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cyb0225/gdfs/pkg/log"
)

// prefix tree
type Tree struct {
	rw   sync.RWMutex
	Root *Node // tree root '/'
	Per  *Persistence
}

// Create a new Tree
func NewTree(storagePath string) (*Tree, error) {
	tree, err := ReadLog(storagePath)
	if err != nil {
		log.Error("read tree.log file failed", log.Err(err))
		// If read tree.log File failed, then delete the old log and create a new file.
		if err = os.Remove(storagePath + "tree.log"); err != nil {
			log.Error("delete fault tree.log file failed", log.String("storagePath", storagePath + "tree.log"), log.Err(err))
		}
		
	}

	per, err := NewPersistence(storagePath)
	if err != nil {
		return nil, fmt.Errorf("new persistence failed: %w", err)
	}

	// Read logging file failed.
	if tree == nil {
		root := NewNode(
			SetFileName("/"),
			SetFilePath("/"),
			IsDirectory(true),
		)
		tree = &Tree{
			Root: root,
		}
	}

	tree.Per = per

	return tree, nil
}

// read tree.log file and Create a root node.
// If the logging file is damaged(can't load it), then return an error.
func ReadLog(storagePath string) (*Tree, error) {
	tree := &Tree{
		Root: NewNode(
			SetFileName("/"),
			SetFilePath("/"),
			IsDirectory(true),
		),
	}

	fd, err := os.Open(storagePath + "tree.log")
	if err != nil {
		return nil, fmt.Errorf("open tree.log faile: %w", err)
	}

	s := &scaner{
		s: bufio.NewScanner(fd),
		tree: tree,
	}

	for s.s.Scan() {
		method := s.s.Text()
		switch method {
		case "put":
			s.put()
		case "rename":
			s.rename()
		case "mkdir":
			s.mkdir()
		case "delete":
			s.delete()
		default:
			continue
		}
		if s.err != nil {
			return nil, fmt.Errorf("read tree.log failed: %w", s.err)
		}
	}

	return s.tree, nil
}

// This struct is used to simplified code.
type scaner struct {
	s *bufio.Scanner
	err error
	tree *Tree
}

// Used to deal with too many scan and text in read log functions.
func (s *scaner) scan() string {
	s.s.Scan()
	return s.s.Text()
}

func (s *scaner)put() {
	if s.err != nil {
		return
	}

	filepath := s.scan()
	filesizeStr := s.scan()
	filesize, err := strconv.ParseInt(filesizeStr, 10, 64)
	if err != nil {
		s.err = fmt.Errorf("transform the filesize of %s to int64 failed: %w", filepath, err)
		return
	}
	keys := s.scan()
	filekeys := strings.Split(keys, " ")
	updateTimeStr := s.scan()
	updateTime, err := time.ParseInLocation("2006-01-02 15:04:05.000", updateTimeStr, time.Local)
	if err != nil {
		s.err = fmt.Errorf("parse %s to time failed: %w",updateTimeStr, err)
		return
	}
	createTimeStr := s.scan()
	createTime, err := time.ParseInLocation("2006-01-02 15:04:05.000", createTimeStr, time.Local)
	if err != nil {
		s.err = fmt.Errorf("parse %s to time failed: %w",createTimeStr, err)
		return
	}
	node := &Node{
		FilePath: filepath,
		FileType: NormalFile,
		FileSize: filesize,
		FileKeys: filekeys,
		UpdateTime: updateTime,
		CreateTime: createTime,
	}

	s.err = s.tree.Put(filepath, node)
}

func (s *scaner)rename() {
	if s.err != nil {
		return
	}

	src := s.scan()
	dest := s.scan()
	s.err = s.tree.Rename(src, dest)
}

func (s *scaner)mkdir() {
	if s.err != nil {
		return
	}

	filepath := s.scan()
	updateTimeStr := s.scan()
	updateTime, err := time.ParseInLocation("2006-01-02 15:04:05.000", updateTimeStr, time.Local)
	if err != nil {
		s.err = fmt.Errorf("parse %s to time failed: %w",updateTimeStr, err)
		return
	}
	createTimeStr := s.scan()
	createTime, err := time.ParseInLocation("2006-01-02 15:04:05.000", createTimeStr, time.Local)
	if err != nil {
		s.err = fmt.Errorf("parse %s to time failed: %w",createTimeStr, err)
		return
	}

	node, err := s.tree.Mkdir(filepath)
	if err != nil {
		s.err = fmt.Errorf("make directory %s failed: %w", filepath, err)
		return
	}
	node.UpdateTime = updateTime
	node.CreateTime = createTime
	
}

func (s *scaner)delete() {
	if s.err != nil {
		return
	}

	filepath := s.scan()	
	if _, err := s.tree.Delete(filepath); err != nil {
		s.err = fmt.Errorf("delete %s failed :%w", filepath, err)
		return
	}
}
