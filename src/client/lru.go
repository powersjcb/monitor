package client

type lruData struct {
	data PingRequest
	prev *lruData
	next *lruData
	sentinel bool
}

type LRU struct {
	head *lruData
	tail *lruData
	m map[PingID]*lruData
	size int
}

func NewLRU(size int) LRU {
	// setup sentinel nodes
	h := &lruData{sentinel: true}
	t := &lruData{sentinel: true}
	h.prev = t
	t.next = h
	return LRU{
		size: size,
		head: h,
		tail: t,
		m: make(map[PingID]*lruData),
	}
}

func (h *LRU) Len() int { return len(h.m) }

func (h *LRU) Add(request PingRequest) {
	r, exists := h.m[request.ID]
	if exists {
		h.remove(r)
	}
	if h.Len() >= h.size {
		h.remove(h.head.prev)
	}
	data := &lruData{
		data: request,
		prev: h.tail,
		next: h.tail.next,
	}
	h.head.prev = data
	h.tail.next = data
	h.m[request.ID] = data
}

func (h *LRU) remove(ejected *lruData) *lruData {
	if ejected == nil {
		return nil
	}
	if ejected.sentinel {
		panic("ejected sentinel node")
	}
	delete(h.m, ejected.data.ID)
	next := ejected.next
	prev := ejected.prev
	next.prev = prev
	prev.next = next
	return ejected
}

func (h *LRU) Remove(id PingID) *PingRequest {
	ejected, exists := h.m[id]
	if exists {
		h.remove(ejected)
		return &ejected.data
	}
	return nil
}
