// keep alive with datanode

package alive

import "sync"

// record alived datanode
type Alive struct {
	rw sync.RWMutex
	// string records the address of datanode, useing empty struct makes map a set
	mp map[string]struct{} 
}


func NewAlive() *Alive {
	return &Alive{
		mp: make(map[string]struct{}),
	}
}

// datanode online
func(a *Alive) Add(key string) {
	a.rw.Lock()
	defer a.rw.Unlock()

	a.mp[key] = struct{}{}
}

// datanode expired
func (a *Alive) Remove(key string) {
	a.rw.Lock()
	defer a.rw.Unlock()
	
	delete(a.mp, key)
}

// check datanode is alive or not
func (a *Alive) IsAlive(key string) bool {
	a.rw.RLock()
	defer a.rw.RUnlock()

	_, ok := a.mp[key]
	return ok
} 