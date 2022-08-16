package alive

import (
	"math/rand"
	"fmt"
	"testing"
	"time"
)
var timeout = 3 * time.Second
func TestAlive(t *testing.T) {
	key := "hello"
	a := NewAlive()
	a.Update(key)

	if ok := a.IsAlive(key); !ok {
		t.Fatalf("IsAlive or Updata function is wrong")
	}
	time.Sleep(timeout + time.Second)
	if ok := a.IsAlive(key); ok {
		t.Fatalf("IsAlive is wrong, key is not expired")
	}
	a.Update(key)
	if ok := a.IsAlive(key); !ok {
		t.Fatalf("Updata function is wrong, update failed")
	}
}


func TestBackup(t *testing.T) {
	a := NewAlive()
	a.Update("1")
	a.Update("2")
	a.Update("3")
	a.Update("4")
	a.Update("5")
	a.Update("6")
	a.Update("7")


	s, err := a.Backup()
	if err != nil {
		t.Fatalf("get backup datanode address failed: %s", err.Error())
	}

	fmt.Println(s)
	time.Sleep(timeout + time.Second)
	if s, err := a.Backup(); err == nil {
		t.Fatalf("datanode expired, but balance still get it. %v", s)
	}

	a.Update("1")
	s, err = a.Backup();
	if err != nil {
		t.Fatalf("get backup datanode address failed: %s", err.Error())
	}
	fmt.Println(s)

}

func TestRand(t *testing.T) {
	_ = NewAlive()
	for i := 0; i < 10; i++ {
		d := rand.Intn(100)
		fmt.Printf("%d ", d)
	}
	fmt.Println()
}

