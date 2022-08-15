package cache

// find which datanode to store this file's chunks and backups
func (c *Cache) Put(filekey, address string) {
	c.rw.Lock()
	defer c.rw.Unlock()

	v, ok := c.mp[filekey]
	if !ok {
		v = NewNode()
		v.Backups = append(v.Backups, address)
		c.mp[filekey] = v
		return
	}

	v.Backups = append(v.Backups, address)
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