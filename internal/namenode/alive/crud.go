package alive

import (
	"fmt"
	"time"
)

// Update datanode time. (keep alive)
func (a *Alive) Update(key string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	v, ok := a.mp[key]
	// The datanode haven't register before.
	if !ok {
		val := &datanodeInfo{
			updateTime: time.Now(),
		}
		a.mp[key] = val
	}
	if v == nil {
		v = &datanodeInfo{}
	}
	v.updateTime = time.Now()
}

// Check if datanode is alived.
func (a *Alive) IsAlive(key string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	val, ok := a.mp[key]
	// datanode haven't register
	if !ok {
		return false
	}

	expireTime := time.Since(val.updateTime)
	// Datanode is expired.
	// Then delete it.
	if expireTime >= a.expired {
		delete(a.mp, key)
		return false
	}
	return true
}

// loadBalance algorithm.
// Simple traversal of the map to achieve load balancing.
func (a *Alive) LoadBalance(backups int) ([]string, error) {
	length := len(a.mp)

	if length == 0 {
		return nil, fmt.Errorf("There is no alived datanode")
	}

	str := make([]string, 0)

	if backups > length {
		backups = length
	}

	i := 0
	for k := range a.mp {
		if i >= backups {
			break
		}
		if ok := a.IsAlive(k); ok {
			str = append(str, k)
			i++
		}
	}

	if len(str) == 0 {
		return nil, fmt.Errorf("datanode are all died")
	}
	return str, nil
}
