package cache

// Get which datanodes stored this filekey.
func (c *Cache) Get(filekey string) []string {
	c.rw.RLock()
	defer c.rw.RUnlock()

	val, ok := c.mp[filekey]
	if !ok {
		return nil
	}

	return val.addrs
}

// Put address to the filekey.
func (c *Cache) Put(filekey, address string) {
	c.rw.Lock()
	defer c.rw.Unlock()

	v, ok := c.mp[filekey]
	if !ok {
		v = newNode()
		v.addrs = append(v.addrs, address)
		c.mp[filekey] = v
		return
	}

	// Check if the address is already exist.
	for i := 0; i < len(v.addrs); i++ {
		if v.addrs[i] == address {
			return
		}
	}
	v.addrs = append(v.addrs, address)
	c.mp[filekey] = v
}

// Delete such filekey in map.
func (c *Cache) Delete(filekey string) {
	c.rw.Lock()
	defer c.rw.Unlock()

	delete(c.mp, filekey)
}
