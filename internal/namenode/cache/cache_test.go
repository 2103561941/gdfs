package cache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	expired := 3
	c := NewCache(expired)
	if err := c.Update("1"); err != nil {
		t.Fatal(err)
	}

	if err := c.Put("f", "1"); err != nil {
		t.Fatal(err)
	}

	v, err := c.Get("f")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)
	time.Sleep(c.expired)
	if err := c.Put("g", "1"); err == nil {
		t.Fatal("Datanode expired but it can still put filekey.")
	}

	if _, err := c.Get("t"); err == nil {
		t.Fatal("Backups can get unexisted filekey.")
	}
}
