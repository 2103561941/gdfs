package cache

import (
	"fmt"
	"time"

	"github.com/cyb0225/gdfs/pkg/log"
)

// Get which datanodes stored this file. 
// And check if datanodes are alived.
// If some datanodes is expired, then delete it.
// Input filename.
// Here is an explanation why I named the filename filekey:
// Datanode stored parts of file which user pushed through clinent.
// So the "real" file's name called filename, and the part of file's name called filekey.

func (c *Cache) Get(filekey string) ([]string, error) {
	val, ok := c.backups.Load(filekey)
	if !ok {
		return nil, fmt.Errorf("filekey not exist")
	}

	backups, ok := val.(filekeyInfo)
	if !ok {
		// refelact failed, check the type of val.
		log.Error("type refelact failed", log.String("filekey", filekey))
		return nil, fmt.Errorf("type refelact failed")
	}

	// select alived backups.  
	alivedBackups := []string{}
	for i := 0; i < len(backups.backups); i++ {
		addr := backups.backups[i]
		val, ok := c.alive.Load(addr)
		if !ok {
			// Datanode haven't register before. 
			continue
		}
		datanode, ok := val.(datanodeInfo)
		if !ok {
			log.Error("type refelact failed", log.String("filekey", filekey), log.String("address", addr))
			continue		
		}
		if time.Since(datanode.alive) >= c.expired { // expired
			c.alive.Delete(addr)
		}
		alivedBackups = append(alivedBackups, addr)
	}
	// Indicates that some datanode have expired
	if len(alivedBackups) != len(backups.backups) {
		c.backups.Store(filekey, alivedBackups)
	}

	return alivedBackups, nil
}

// Bind the relation between filekey and the address of datanode.
// Input filename and the address of datanode.
func (c *Cache) Put(filekey string, address string) error {
	// Check if datanode is register.
	aliveVal, ok := c.alive.Load(address)
	if !ok { // not register
		return fmt.Errorf("datanode [%s] haven't register before.", address)
	}

	datanode, ok := aliveVal.(datanodeInfo)
	if !ok {
		log.Error("type refelact failed", log.String("filekey", filekey), log.String("address", address))
		return fmt.Errorf("type refelact failed")
	}
	// Check if the datanode is alive.
	if time.Since(datanode.alive) >= c.expired { // expired
		return fmt.Errorf("datanode [%s] is expired", address)
	}

	// Do some nil check.
	if datanode.filekeys == nil {
		datanode.filekeys = make([]string, 0)
	}
	datanode.filekeys = append(datanode.filekeys, filekey)
	// Update datanode's alive time
	datanode.alive = time.Now()
	c.alive.Store(address, datanode)

	backupsVal, ok := c.backups.Load(filekey)
	if !ok { // Filekey haven't register before.
		// Make a new store.
		info := filekeyInfo{
			backups: []string{address},
		}
		c.backups.Store(filekey, info)
		return nil
	}

	// Append address to filekey's backups.
	backups, ok := backupsVal.(filekeyInfo)
	if !ok {
		log.Error("type refelact failed", log.String("filekey", filekey), log.String("address", address))
		return fmt.Errorf("type refelact failed")
	}

	if backups.backups == nil {
		backups.backups = make([]string, 0)
	}
	backups.backups = append(backups.backups, address)
	c.backups.Store(filekey, backups)

	return nil
}

// Update Datanode. (keep alive)
// Input the address of datanode.
// Update time corresponding datanode in alive map to current time.
// If there is no such datanode, then create it.
func (c *Cache) Update(address string) error {
	aliveVal, ok := c.alive.Load(address)
	if !ok { // Datanode haven't register before.
		// Register a new address.
		datanode := datanodeInfo{
			alive:    time.Now(),
			filekeys: make([]string, 0),
		}
		c.alive.Store(address, datanode)
		return nil
	}

	datanode := aliveVal.(datanodeInfo)
	datanode.alive = time.Now()
	c.alive.Store(address, datanode)

	return nil
}

// Delete filekey
// Input filekey, and delete in map.
// Return addresses of datanodes used to tell datanodes to delete the file named filekey.
// Delete datanode's filekey in memory(alive map).
func (c *Cache) Delete(filekey string) ([]string, error) {
	backupsVal, ok := c.backups.Load(filekey)
	if !ok { // There didn't store such filekey.
		return nil, nil
	}

	backups, ok := backupsVal.(filekeyInfo)
	if !ok {
		log.Error("type refelact failed", log.String("filekey", filekey))
		return nil, fmt.Errorf("type refelact failed")
	}

	// addrs is addresses of datanodes which stored this filekey.
	addrs := backups.backups

	// delete the filekey in datanodes.
	for i := 0; i < len(addrs); i++ {
		addr := addrs[i]
		aliveVal, ok := c.alive.Load(addr)
		if !ok { // Datanode not exist.
			continue
		}

		datanode, ok := aliveVal.(datanodeInfo)
		if !ok {
			log.Error("type refelact failed", log.String("filekey", filekey), log.String("address", addr))
			return nil, fmt.Errorf("type refelact failed")
		}
		// Search the index of filekey in datanode
		index := -1 
		for j := 0; j < len(datanode.filekeys); j++ {
			if datanode.filekeys[j] == filekey {
				index = j
				break
			}
		}
		// Delete filekey in datanode's filekeys.
		datanode.filekeys = append(datanode.filekeys[:index], datanode.filekeys[index+1:]...)
		c.alive.Store(addr, datanode)
	}

	// Delete filekey
	c.backups.Delete(filekey) 

	return addrs, nil
}

// Load Balance.
// According to the inputparam, give addresses of datanodes.
// if there are not enough suitable datanodes to given, then it will return as many datanode as possible.

// Design:
// Load balance algorithm needs to take into account of datanode's current capacity and its partition.
// If a new datanode connects to gdfs,  it will be seletced as the storaged node first.(Provided that its partition is suitable)
// The partition is used to distinguish different engine room at real world.
// The same backups are best stored in different partition.
func (c *Cache) LoadBalance(backups int) ([]string, error) {
	// if c.alive. < backups {
	// 	backups = len(c.alive)
	// }

	addrs := []string{"127.0.0.1:5000", "127.0.0.1:5001"}
	return addrs, nil
}