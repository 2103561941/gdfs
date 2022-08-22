package cache

import "sync"

// map string stored file's uuid, which is stored in filetree and datanode.
type Cache struct {
	rw sync.RWMutex
	mp map[string]*filekeyInfo
}

// create a new cache
func NewCache() *Cache {
	return &Cache{
		mp: make(map[string]*filekeyInfo),
	}
}

// Filekey's Information
type filekeyInfo struct {
	addrs []string // stored which datanodes stored this filekey.
}

func newNode() *filekeyInfo {
	return &filekeyInfo{
		addrs: make([]string, 0),
	}
}
