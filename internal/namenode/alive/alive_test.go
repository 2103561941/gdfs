package alive

import (
	"fmt"
	"testing"
)

func TestAlive(t *testing.T) {
	expired := 100
	a := NewAlive(expired)
	a.Update("1", 1000)
	a.Update("2", 2000)
	a.Update("3", 3000)
	a.Update("4", 4000)

	for e := a.ll.Front(); e != nil; e = e.Next() {
		t.Logf("%+v \n", e.Value.(*entry))
	}

	
	v, err := a.LoadBalance(10)
	if err != nil {
		t.Fatalf("get loadbalance failed: %s", err.Error())
	}
	fmt.Printf("%+v\n", v)

	fmt.Println()
	a.Update("1", 10000)
	for e := a.ll.Front(); e != nil; e = e.Next() {
		t.Logf("%+v \n", e.Value.(*entry))
	}

	v, err = a.LoadBalance(3)
	if err != nil {
		t.Fatalf("get loadbalance failed: %s", err.Error())
	}
	fmt.Printf("%+v\n", v)

	// time.Sleep(a.expired)
	// if ok := a.IsAlive("1"); ok {
	// 	t.Fatal("isalived is error")
	// }

}
