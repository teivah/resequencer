// Package resequencer handles resequencing operations.
package resequencer

import (
	"context"
	"sync"

	"github.com/emirpasic/gods/trees/binaryheap"
)

// Handler is the parent struct handling the resequencing operations.
type Handler struct {
	expected int
	push     chan struct{}
	ch       chan []int
	mu       sync.Mutex
	minHeap  *binaryheap.Heap
}

// NewHandler creates a new Handler.
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

// Push adds new sequence identifiers.
func (h *Handler) Push(sequenceIDs ...int) {
	h.mu.Lock()
	for _, sequenceID := range sequenceIDs {
		h.minHeap.Push(sequenceID)
	}
	h.mu.Unlock()
	h.push <- struct{}{}
}

// Messages returns the channel of messages.
func (h *Handler) Messages() <-chan []int {
	return h.ch
}
