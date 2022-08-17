// keep alive with datanode

package alive

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/cyb0225/gdfs/internal/namenode/config"
	"github.com/cyb0225/gdfs/pkg/log"
)

// record alived datanode
type Alive struct {
	mu sync.Mutex

	// string records the address of datanode
	// it stored last := time datanode heartbeat to namenode, when I need to use

	// this datanode, I need to check whether the := time is expired or not.
	mp map[string]time.Time //

	balance []string
}

func init() {
}

func NewAlive() *Alive {
	rand.Seed(time.Now().UnixMicro())
	return &Alive{
		mp:      make(map[string]time.Time),
		balance: make([]string, 0),
	}
}

// datanode online
func (a *Alive) Update(key string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	_, ok := a.mp[key]
	if !ok { // not register before
		a.balance = append(a.balance, key)
		a.mp[key] = time.Now()
	}

	a.mp[key] = time.Now()
}

// check datanode is alive or not
// if datanode is not alive, then kick it out
func (a *Alive) IsAlive(key string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	t, ok := a.mp[key]
	if !ok { // datanode haven't register
		// log.Debug("key not exist", log.String("key", key))
		delete(a.mp, key)
		return false
	}

	expireTime := time.Since(t)
	// timeout := time.Duration(config.Cfg.Timeout) * time.Second
	timeout := time.Duration(config.Cfg.Timeout) * time.Second

	return expireTime < timeout
}

// backup needs to think about the number of datanode and backups,
// because, stored the same backups in one datanode is useless,
// so choose Backup datanode should ompare the number of servers and the number of backups which is smaller
func (a *Alive) Backup() ([]string, error) {
	length := len(a.mp)
	log.Debug("alive datanode", log.Int("datanode num", length))

	if length == 0 {
		return nil, fmt.Errorf("there is no alived datanode")
	}

	str := make([]string, 0)

	var n int
	if config.Cfg.BackupN < length {
		n = config.Cfg.BackupN
	} else {
		n = length
	}

	i := 0
	log.Debug("backup num", log.Int("n", n), log.Int("backups", config.Cfg.BackupN))
	for k, _ := range a.mp {
		if i >= n {
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
