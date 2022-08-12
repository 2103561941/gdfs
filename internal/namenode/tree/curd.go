package tree

import (
	"log"
	"strings"
)

var (
	search = getSearch()
	mkdir  = GetMkdir()
)

// get file's information.
func (t *Tree) Get(filepath string) *Node {
	t.rw.RLock()
	defer t.rw.RUnlock()

	patterns := split(filepath)
	node := search(patterns, t.Root)
	return node
}

// put a node to the perfix tree.
// it can not only put a normal file but a directory, it depends on the file Type of node.
// user should pass the whole filepath(filepath + filename) when insert node.
func (t *Tree) Put(filepath string, node *Node) {
	t.rw.Lock()
	defer t.rw.Unlock()
	patterns := split(filepath)

	lastIndex := len(patterns) - 1
	pattern := patterns[lastIndex]
	parentDir := mkdir(patterns[:lastIndex], t.Root)

	// patentDir may not exist, it can be design as linux.
	// for example, I pass /tmp/cyb/demo.txt, but there is not have directory named 'cyb'
	// So, I can create the directory 'cyb', instead of return an error.

	node.Partten = pattern
	node.FileMeta.Filename = pattern
	parentDir.Children = append(parentDir.Children, node)
}

// searchFunc file from the prefix tree
// filepath equals filepath + filename
// recursively searchFunc for matching values
// here, I use closure to store height
type searchFunc func(patterns []string, node *Node) *Node

func getSearch() searchFunc {
	var fn searchFunc
	height := 0
	fn = func(patterns []string, node *Node) *Node {
		if patterns[height] != node.Partten {
			return nil
		}

		// it used to test funciton in crud.
		log.Printf("node: %#v\n\n", node.FileMeta)

		if height == len(patterns)-1 {
			return node
		}

		for i := 0; i < len(node.Children); i++ {
			height++
			ans := fn(patterns, node.Children[i])
			if ans != nil {
				return ans
			}
			height--
		}
		return nil
	}
	return fn
}

// this function is use to create directories recursively.

func GetMkdir() searchFunc {
	height := 0
	var fn searchFunc
	fn = func(patterns []string, node *Node) *Node {
		if patterns[height] != node.Partten {
			return nil
		}

		if height == len(patterns)-1 {
			return node
		}

		// Notice!
		// height must be add at the out of for-loop.
		height++
		for i := 0; i < len(node.Children); i++ {
			ans := fn(patterns, node.Children[i])
			if ans != nil {
				return ans
			}
		}
		// here is the difference with search
		// it chooses to make a directory, instead of return nil

		// Notice!
		// patterns[height] is the children of node, and we need to make a

		dir := &Node{
			Partten: patterns[height],
			Children: make([]*Node, 0),
			FileMeta: &Meta{
				Filename: patterns[height],
				FileType: Direcotry,
			},
		}
		
		node.Children = append(node.Children, dir)
		return fn(patterns, dir)
	}
	return fn
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
