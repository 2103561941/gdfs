// keep alive with datanode

package alive

import (
	"sync"
	"time"
)

var (
	timeout = time.Second * 3 // 20s
)

// record alived datanode
type Alive struct {
	rw sync.RWMutex

	// string records the address of datanode
	// it stored last := time datanode heartbeat to namenode,
	// when I need to use

	// this datanode, I need to check whether the := time is expired or not.
	mp map[string]time.Time //

}

func NewAlive() *Alive {
	return &Alive{
		mp: make(map[string]time.Time),
	}
}

// datanode online
func (a *Alive) Update(key string) {
	a.rw.Lock()
	defer a.rw.Unlock()

	a.mp[key] = time.Now()
}

// check datanode is alive or not
func (a *Alive) IsAlive(key string) bool {
	a.rw.RLock()
	defer a.rw.RUnlock()

	t, ok := a.mp[key]
	if !ok { // datanode haven't register
		return false
	}

	// datanode expired
	duration := time.Since(t)
	if duration > timeout {
		return false
	}

	return true
}

// datanode expired, the time will expired, don't need to remove it
// func (a *Alive) Remove(key string) {
// 	a.rw.Lock()
// 	defer a.rw.Unlock()

// 	delete(a.mp, key)
// }
