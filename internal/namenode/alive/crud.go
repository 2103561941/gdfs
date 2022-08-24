package alive

import (
	"container/list"
	"time"
)

// Update datanode time. (keep alive)
// If datanode is not register before, then create it in list and map.
func (a *Alive) Update(addr string, cap int64) {
	a.mu.Lock()
	defer a.mu.Unlock()

	v, ok := a.mp[addr]
	// The datanode haven't register before.
	// register in map and list.
	if !ok {
		val := &entry{
			addr:       addr,
			updateTime: time.Now(),
			cap:        cap,
		}

		// Traverse the list to find suitable index to stored this element.
		// List is sorted by cap from small to large.
		var ele *list.Element
		for ele = a.ll.Front(); ele != nil; ele = ele.Next() {
			entry := ele.Value.(*entry)
			if entry.cap > val.cap {
				ele = a.ll.InsertBefore(val, ele)
				break
			}
		}

		if ele == nil {
			ele = a.ll.PushBack(val)
		}

		a.mp[addr] = ele
		return
	}

	val := v.Value.(*entry)
	val.updateTime = time.Now()
	if val.cap == cap {
		return
	}
	val.cap = cap

	// Resort by move element to another location.
	var e *list.Element
	for e = a.ll.Front(); e != nil; e = e.Next() {
		ele := e.Value.(*entry)
		if ele.cap > val.cap {
			a.ll.MoveBefore(v, e)
			return
		}
	}
	// Move to the back
	a.ll.MoveToBack(v)

}

// Check if datanode is alived.
func (a *Alive) IsAlive(addr string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	e, ok := a.mp[addr]
	// datanode haven't register
	if !ok {
		return false
	}

	val := e.Value.(*entry)
	expireTime := time.Since(val.updateTime)
	// Datanode is expired.
	// Then delete it.
	if expireTime >= a.expired {
		delete(a.mp, addr)
		a.ll.Remove(e)
		return false
	}
	return true
}

// LoadBalance: sorted by the capacity of datanode.
// Perfer datanodes with small cap and addresses of returns are unique.
// Input how many datanode should be returned. But it also depends on how many alived datanodes I have.
// Return a serious of datanode to stroed related file.
func (a *Alive) LoadBalance(n int) ([]string, error) {
	backups := []string{}
	i := 0
	for ele := a.ll.Front(); ele != nil; ele = ele.Next() {
		i++
		if i > n {
			break
		}
		val := ele.Value.(*entry)
		expireTime := time.Since(val.updateTime)
		// Datanode is expired.
		// Then delete it.
		if expireTime >= a.expired {
			delete(a.mp, val.addr)
			a.ll.Remove(ele)
			continue
		}

		backups = append(backups, val.addr)
	}

	return backups, nil
}
