package cache

import "fmt"

// find which datanode to store this file's chunks and backups
func (c *Cache) Put(filekey, address string) error {
	c.rw.Lock()
	defer c.rw.Unlock()

	v, ok := c.mp[filekey]
	
	// it turns out that the file in namenode is deleted.
	// so here return error, tell datanode delete the extra file.
	if !ok { 
		return fmt.Errorf("file not exist in namenode")
	}

	v.Backups = append(v.Backups, address)
	return nil
}

// create filekey in cache
func (c *Cache) Create(filekey string) error {
	c.rw.Lock()
	defer c.rw.Unlock()

	_, ok := c.mp[filekey]
	if ok {
		return fmt.Errorf("filekey already exist")
	}

	c.mp[filekey] = NewNode()
	return nil
}


// get which datanode stored this files' chunks and backups.
func (c *Cache) Get(filekey string) *Node {
	c.rw.RLock()
	defer c.rw.RUnlock()

	node, ok := c.mp[filekey]
	if !ok {
		return nil
	}

	return node
}

func (c *Cache) Delete(filekey string) {
	c.rw.Lock()
	defer c.rw.Unlock()

	delete(c.mp, filekey)
}