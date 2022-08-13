package cache

// find which datanode to store this file's chunks and backups
func (c *Cache) Put(key string) error {
	c.rw.Lock()
	defer c.rw.Unlock()

	return nil
}

// get which datanode stored this files' chunks and backups.
func (c *Cache) Get(key string) *Node {
	c.rw.RLock()
	defer c.rw.Unlock()

	node, ok := c.mp[key]
	if !ok {
		return nil
	}

	return node
}

// choose which datanode to store this file block
func storeBalance() {}

// chunk a big file to some smaller files.
func chunkfile() {}
