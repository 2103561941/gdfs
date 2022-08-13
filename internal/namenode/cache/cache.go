package cache

import "sync"

// map string stored file's uuid, which is stored in filetree and datanode.
type Cache struct {
	rw sync.RWMutex
	mp map[string]*Node
}

// create a new cache
func NewCache() *Cache {
	return &Cache{
		mp: make(map[string]*Node, 0),
	}
}

// datanode infomation
type Node struct {
	Backups []*Backup
}

// file backups, stroed the message of datanode
type Backup struct {
	Address string // ip + port
}