package tree

import (
	"fmt"
	"log"
	"strings"
)

// get file's information.
func (t *Tree) Get(filepath string) *Node {
	t.rw.RLock()
	defer t.rw.RUnlock()

	patterns := split(filepath)
	node := search(patterns, t.Root, 0)
	return node
}

// put a node to the perfix tree.
// it can not only put a normal file but a directory, it depends on the file Type of node.
// user should pass the whole filepath(filepath + filename) when insert node.
func (t *Tree) Put(filepath string, node *Node) error {
	t.rw.Lock()
	defer t.rw.Unlock()
	patterns := split(filepath)

	lastIndex := len(patterns) - 1
	pattern := patterns[lastIndex]
	parentDir := mkdir(patterns[:lastIndex], t.Root, 0)
	if parentDir == nil {
		return fmt.Errorf("make directory failed")
	}
	// patentDir may not exist, it can be design as linux.
	// for example, I pass /tmp/cyb/demo.txt, but there is not have directory named 'cyb'
	// So, I can create the directory 'cyb', instead of return an error.

	node.FileName = pattern
	// check if the file is exist or not,
	// if file is already existed, then return an error
	for i := 0; i < len(parentDir.Children); i++ {
		if parentDir.Children[i].FileName == pattern {
			return fmt.Errorf("file: %s is already existed.", filepath)
		}
	}
	parentDir.Children = append(parentDir.Children, node)

	return nil
}

// searchFunc file from the prefix tree
// filepath equals filepath + filename
// recursively searchFunc for matching values

func search(patterns []string, node *Node, height int) *Node {
	if patterns[height] != node.FileName {
		return nil
	}

	// it used to test funciton in crud.
	log.Printf("node: %#v\n\n", node)

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

// this function is use to create directories recursively.
func mkdir(patterns []string, node *Node, height int) *Node {
	if patterns[height] != node.FileName {
		return nil
	}

	if height == len(patterns)-1 {
		return node
	}

	// Notice!
	// height must be add at the out of for-loop.
	height++
	for i := 0; i < len(node.Children); i++ {
		ans := mkdir(patterns, node.Children[i], height)
		if ans != nil {
			return ans
		}
	}
	// here is the difference with search
	// it chooses to make a directory, instead of return nil

	// Notice!
	// patterns[height] is the children of node, and we need to make a
	dir := NewNode(
		SetFileName(patterns[height]),
		IsDirectory(true),
	)

	node.Children = append(node.Children, dir)
	return mkdir(patterns, dir, height)
}

// split the file path to a string slice
// the first "/" stands for the root dir, so it need to be add.
// the check of filepath is placed on the gdfs client
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


// chunck a big file to some smaller files.
func (t *Tree)chunkfile() {

}
