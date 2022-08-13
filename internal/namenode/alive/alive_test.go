package alive

import "testing"

func TestAlive(t *testing.T) {
	key := "hello"
	a := NewAlive()
	a.Add(key)
	if ok := a.IsAlive(key); !ok {
		t.Fatalf("IsAlive or Add function is wrong")
	}
	a.Remove(key)
	if ok := a.IsAlive(key); ok {
		t.Fatalf("Remove is wrong")
	}
	a.Remove(key)
}
