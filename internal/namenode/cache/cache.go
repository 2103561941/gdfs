package cache

import (
	"sync"
	"time"
)

// Cache recored the relation between filekey and datanode.
// Notice: The key of map is string and the value of map is struct, not a pointer.  
type Cache struct {
	backups  sync.Map // key: filekey, value: datanodeInfo.
	alive sync.Map	  // key: datanodes, value: filekeyInfo.
	expired time.Duration 
}

type filekeyInfo struct {
	backups []string
}

type datanodeInfo struct {
	alive    time.Time
	filekeys []string
}

// create a new cache
func NewCache(expired int) *Cache {
	return &Cache{
		backups: sync.Map{},
		alive:   sync.Map{},
		expired: time.Second * time.Duration(expired),
	}	
}
