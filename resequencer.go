package resequencer

import (
	"context"
	"sync"

	"github.com/emirpasic/gods/trees/binaryheap"
)

type Handler struct {
	expected int
	push     chan struct{}
	ch       chan []int
	once     sync.Once
	closed   bool
	mu       sync.Mutex
	minHeap  *binaryheap.Heap
}

func NewHandler(ctx context.Context, current int) *Handler {
	handler := &Handler{
		expected: current + 1,
		push:     make(chan struct{}, 1),
		ch:       make(chan []int, 1),
		minHeap: binaryheap.NewWith(
			func(a, b interface{}) int {
				return a.(int) - b.(int)
			},
		),
	}
	go handler.process(ctx)
	return handler
}

func (h *Handler) process(ctx context.Context) {
	var send []int
	for {
		select {
		case <-h.push:
			h.mu.Lock()
			for !h.minHeap.Empty() {
				v, _ := h.minHeap.Peek()
				peek := v.(int)
				if peek != h.expected {
					break
				}
				send = append(send, peek)
				h.expected++
				h.minHeap.Pop()
			}
			h.mu.Unlock()
			if len(send) != 0 {
				h.ch <- send
				send = make([]int, 0)
			}
		case <-ctx.Done():
			close(h.ch)
			return
		}
	}
}

func (h *Handler) Pop(i int) {
	h.mu.Lock()
	h.minHeap.Push(i)
	h.mu.Unlock()
	h.push <- struct{}{}
}

func (h *Handler) Messages() <-chan []int {
	return h.ch
}
