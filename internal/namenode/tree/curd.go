package tree

import (
	"fmt"
	"strings"
	"time"

	"github.com/cyb0225/gdfs/pkg/log"
)

// get file's information.
func (t *Tree) Get(filepath string) (*Node, error) {
	patterns := split(filepath)
	if err := checkFilename(patterns); err != nil {
		return nil, err
	}

	t.rw.RLock()
	node := search(patterns, t.Root, 0)
	t.rw.RUnlock()

	if node == nil {
		return nil, fmt.Errorf("cannot find the file")
	}

	if node.IsDirectory() {
		return nil, fmt.Errorf("file is a directory")
	}

	log.Info("get file success", log.String("file", filepath))
	return node, nil
}

// put a node to the perfix tree.
// it can not only put a normal file but a directory, it depends on the file Type of node.
// user should pass the whole filepath(filepath + filename) when insert node.
func (t *Tree) Put(filepath string, node *Node) error {
	patterns := split(filepath)
	if err := checkFilename(patterns); err != nil {
		return err
	}

	if err := checkHaveDir(patterns); err != nil {
		return err
	}

	lastIndex := len(patterns) - 1
	putFile := patterns[lastIndex]
	node.FileName = putFile

	t.rw.RLock()
	parentDir := search(patterns[:lastIndex], t.Root, 0)
	t.rw.RUnlock()

	if parentDir == nil {
		return fmt.Errorf("parent dir %s not exist", filepath)
	}

	t.rw.Lock()
	defer t.rw.Unlock()
	if err := parentDir.AppendChild(node); err != nil {
		return err
	}

	log.Info("put file success", log.String("file", filepath))
	return nil
}

// delete file/directory
func (t *Tree) Delete(filepath string) (*Node, error) {
	if filepath == "/" {
		return nil, fmt.Errorf("cannot delete root file")
	}

	patterns := split(filepath)
	if err := checkFilename(patterns); err != nil {
		return nil, err
	}
	if err := checkHaveDir(patterns); err != nil {
		return nil, err
	}

	lastIndex := len(patterns) - 1
	delFile := patterns[lastIndex] // wait to be deleted
	parentPatterns := patterns[:lastIndex]

	// get its parent's node.
	t.rw.RLock()
	parentDir := search(parentPatterns, t.Root, 0)
	t.rw.RUnlock()
	if parentDir == nil {
		return nil, fmt.Errorf("file not exist")
	}

	// get the index of delfile
	delIndex := -1
	for i := 0; i < len(parentDir.Children); i++ {
		if parentDir.Children[i].FileName == delFile {
			delIndex = i
		}
	}

	if delIndex == -1 {
		return nil, fmt.Errorf("file not exist")
	}

	node := parentDir.Children[delIndex]

	t.rw.Lock()
	defer t.rw.Unlock()
	// delete elmenet from slice
	parentDir.Children = append(parentDir.Children[:delIndex], parentDir.Children[delIndex+1:]...)

	log.Info("delete file success", log.String("file", filepath))
	return node, nil
}

// do not allowed create recursively
func (t *Tree) Mkdir(filepath string) (*Node, error){
	row := split(filepath)
	log.Debugf("row: %+v", row)

	// check directory name
	// if filepath is "/123///" => "/","123","","","", , then skip extra "/", it will become "/123"
	// lastIndex stored the last legal index.
	var lastIndex int
	// find the last non-empty pattern.
	for i := len(row) - 1; i >= 0; i-- {
		lastIndex = i
		if len(row[i]) > 0 {
			break
		}
	}

	patterns := row[:lastIndex+1]
	log.Debugf("row: %+v", patterns)

	if err := checkHaveDir(patterns); err != nil {
		return nil, err
	}

	// check if the directory is legal or not.
	if err := checkFilename(patterns); err != nil {
		return nil, err
	}

	t.rw.RLock()
	parentDir := search(patterns[:lastIndex], t.Root, 0)
	t.rw.RUnlock()

	if parentDir == nil {
		return nil, fmt.Errorf("dir: %s not exist", parentDir.FileName)
	}

	log.Debug("get parentdir", log.String("parentdir", parentDir.FileName))
	node := NewNode(
		IsDirectory(true),
		SetFileName(patterns[lastIndex]),
		SetFilePath(filepath),
	)

	t.rw.Lock()
	defer t.rw.Unlock()
	if err := parentDir.AppendChild(node); err != nil {
		return nil, err
	}

	log.Info("make directory success", log.String("path", filepath))
	return node, nil
}

// rename / move file to another filepath, change the last pattern
func (t *Tree) Rename(src, dest string) error {
	// can not change the root directory name.
	// and can not change file name to root.
	if src == "/" {
		return fmt.Errorf("cannot rename root file")
	}

	if src == dest {
		return nil
	}

	srcPatterns := split(src)
	if err := checkFilename(srcPatterns); err != nil {
		return fmt.Errorf("src file %w", err)
	}

	destPatterns := split(dest)
	if err := checkFilename(destPatterns); err != nil {
		return fmt.Errorf("dest file : %w", err)
	}

	// check if the filepath is same or not.
	if len(srcPatterns) != len(destPatterns) {
		return fmt.Errorf("invaild dest name")
	}

	for i := 0; i < len(srcPatterns)-1; i++ {
		if srcPatterns[i] != destPatterns[i] {
			return fmt.Errorf("invaild dest name")
		}
	}

	if err := checkHaveDir(srcPatterns); err != nil {
		return err
	}

	// check the src file is exist or not.
	// check if the new name is exist or not
	// get the parent directory, and check the children files.
	lastIndex := len(srcPatterns) - 1
	parentPatterns := srcPatterns[:lastIndex]
	t.rw.RLock()
	parentDir := search(parentPatterns, t.Root, 0)
	t.rw.RUnlock()
	if parentDir == nil {
		return fmt.Errorf("src file not exist")
	}
	// if parentdir is a directory. It means that /parentDir/src not exist.
	if !parentDir.IsDirectory() {
		return fmt.Errorf("src file not exist")
	}

	srcName := srcPatterns[lastIndex]
	destName := destPatterns[lastIndex]

	isExist := false
	srcIndex := 0 // store the src index in parentDir slice
	for i := 0; i < len(parentDir.Children); i++ {
		// it comes out that src is exist.
		if parentDir.Children[i].FileName == srcName {
			isExist = true
			srcIndex = i
			continue
		}
		// dest file is exist.
		if parentDir.Children[i].FileName == destName {
			return fmt.Errorf("dest file is exist")
		}
	}
	if !isExist { // src file not exist
		return fmt.Errorf("src file not exist")
	}

	// change the file name
	t.rw.Lock()
	defer t.rw.Unlock()
	parentDir.Children[srcIndex].FileName = destName
	parentDir.Children[srcIndex].UpdateTime = time.Now()

	log.Info("change file name successed", log.String("src", src), log.String("dest", dest))
	return nil
}

// list file's stat
func (t *Tree) Stat(filepath string) (*Node, error) {
	patterns := split(filepath)
	if err := checkFilename(patterns); err != nil {
		return nil, err
	}

	t.rw.RLock()
	node := search(patterns, t.Root, 0)
	t.rw.RUnlock()
	if node == nil {
		return nil, fmt.Errorf("file not exist")
	}

	log.Info("stat file success", log.String("file", filepath))
	return node, nil
}

// list files' stat in directory
func (t *Tree) List(dirpath string) ([]*Node, error) {
	patterns := split(dirpath)
	if err := checkFilename(patterns); err != nil {
		return nil, err
	}

	t.rw.RLock()
	dir := search(patterns, t.Root, 0)
	t.rw.RUnlock()
	if dir == nil {
		return nil, fmt.Errorf("directory not exist")
	}

	if !dir.IsDirectory() {
		return nil, fmt.Errorf("file is not a directory")
	}

	log.Info("list dir success", log.String("dir", dirpath))
	return dir.Children, nil
}

// return node's children node, include child's children
// it will return both directory and normal file.
func (t *Tree) GetChildrenNode(node *Node) []*Node {
	t.rw.RLock()
	defer t.rw.RUnlock()

	var nodes []*Node
	var stack = []*Node{node}
	for len(stack) > 0 {
		lastIndex := len(stack) - 1
		popNode := stack[lastIndex]
		stack = stack[:lastIndex]
		if popNode.IsDirectory() {
			stack = append(stack, popNode.Children...)
		}
		nodes = append(nodes, popNode)
	}

	return nodes
}

// searchFunc file from the prefix tree
// filepath equals filepath + filename
// recursively searchFunc for matching values
func search(patterns []string, node *Node, height int) *Node {
	if patterns == nil {
		return nil
	}

	if patterns[height] != node.FileName {
		return nil
	}

	// it used to test funciton in crud.
	// log.Printf("node: %#v\n\n", node)

	if height == len(patterns)-1 {
		return node
	}

	for i := 0; i < len(node.Children); i++ {
		ans := search(patterns, node.Children[i], height+1)
		if ans != nil {
			return ans
		}
	}
	return nil
}


// Split the file path to a string slice
// The first "/" stands for the root dir, so it need to be add.
// The check of filepath is placed on the gdfs client
func split(filepath string) []string {
	if len(filepath) == 0 {
		return []string{}
	}

	// split("/", "/") ====> []string{"", ""}
	if filepath == "/" {
		return []string{"/"}
	}

	patterns := strings.Split(filepath, "/")
	// example: /tmp/cyb, than the patterns[0] == ""
	patterns[0] = "/"

	return patterns
}

// For example, '/root//t' is a invailed filename. 
func checkFilename(patterns []string) error {
	if len(patterns) == 0 {
		return fmt.Errorf("invaild file name")
	}

	if len(patterns) == 1 && patterns[0] != "/" {
		return fmt.Errorf("invaild file name")
	}

	for i := 0; i < len(patterns); i++ {
		if len(patterns[i]) == 0 {
			return fmt.Errorf("invaild file name")
		}
	}

	return nil
}

// Check if this file is root. For example, its filepath is '/' or 'abc'.
// But not check the dir exist or not.
func checkHaveDir(patterns []string) error {
	// if patterns < 2, then the file can not find the parentDir
	if len(patterns) <= 1 {
		return fmt.Errorf("invailed file path, file must have a dir")
	}
	return nil
}
