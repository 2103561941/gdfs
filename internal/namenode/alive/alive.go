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
	mu      sync.Mutex
	mp      map[string]*datanodeInfo

	alive   list.List  // Used to loadbalance.
	expired time.Duration // heartbeat timeout.
}

type datanodeInfo struct {
	updateTime time.Time
}


func NewAlive(expired int) *Alive {
	return &Alive{
		mp:      make(map[string]*datanodeInfo),
		expired: time.Duration(expired) * time.Second,
	}
}
