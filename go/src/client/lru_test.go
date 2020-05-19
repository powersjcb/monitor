package client_test

import (
	"github.com/powersjcb/monitor/go/src/client"
	"testing"
	"time"
)

var r1 = client.PingRequest{
	ID:     123,
	Seq:    0,
	Target: "192.168.1.1",
	SentAt: time.Now(),
}

var r2 = client.PingRequest{
ID:     543,
Seq:    0,
Target: "192.168.1.1",
SentAt: time.Now(),
}

func TestLRU_Add(t *testing.T) {
	lru := client.NewLRU(1)

	// inserts one item
	lru.Add(r1)
	if lru.Len() != 1 {
		t.Error("failed to insert new item")
		return
	}

	// restricts size
	lru.Add(r2)
	if lru.Len() > 1 {
		t.Errorf("does not restrict size: lru.Len() == %d", lru.Len())
		return
	}

	// try ejecting missing element
	r := lru.Remove(r1.ID)
	if r != nil || lru.Len() != 1 {
		t.Errorf("should not eject anything: %v", r)
		return
	}
}

func TestLRU_Remove(t *testing.T) {
	lru := client.NewLRU(2)

	if lru.Len() != 0 {
		t.Error("lru should initialize empty")
	}

	r := lru.Remove(1)
	if r != nil {
		t.Error("should return nil when empty")
	}

	lru.Add(r1)
	if lru.Len() != 1 {
		t.Errorf("lru should contain 1 item: lru.Len() == %d", lru.Len())
	}

	r = lru.Remove(r1.ID)
	if r == nil {
		t.Error("Remove did not return item")
	}
	if lru.Len() != 0 {
		t.Error("lru should be empty")
	}
}