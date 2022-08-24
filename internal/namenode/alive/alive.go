// keep alive with datanode

package alive

import (
	"container/list"
	"sync"
	"time"
)

// Function:
// Alive is used to update datanode's alived status.
// Check is datanode alived.
// Give suitable addresses of datanodes to stored new filekeys.
type Alive struct {
	mu sync.Mutex
	mp map[string]*list.Element

	// Used to loadbalance.
	// Stored the sort of a datanodes' capcity.
	ll      *list.List
	expired time.Duration // heartbeat timeout.
}

type entry struct {
	addr string
	updateTime time.Time
	cap int64 // the capacity of the datanode. 
}

func NewAlive(expired int) *Alive {
	return &Alive{
		mp:      make(map[string]*list.Element),
		expired: time.Duration(expired) * time.Second,
		ll:      list.New(),
	}
}
