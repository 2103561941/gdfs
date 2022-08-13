package alive

import (
	"testing"
	"time"
)

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
