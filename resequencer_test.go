package resequencer

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNominal(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	handler := NewHandler(ctx, -1)
	handler.Pop(0)
	handler.Pop(1)
	handler.Pop(2)
	cancel()
	time.Sleep(100 * time.Millisecond)
	assertHandler(t, handler, []int{0, 1, 2})
}

func TestOutOfOrder(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	handler := NewHandler(ctx, -1)
	handler.Pop(2)
	handler.Pop(1)
	handler.Pop(0)
	cancel()
	time.Sleep(100 * time.Millisecond)
	assertHandler(t, handler, []int{0, 1, 2})
}

func TestMissing(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	handler := NewHandler(ctx, -1)
	handler.Pop(2)
	handler.Pop(3)
	handler.Pop(0)
	cancel()
	time.Sleep(100 * time.Millisecond)
	assertHandler(t, handler, []int{0})
}

func assertHandler(t *testing.T, h *Handler, expected []int) {
	got := make([]int, 0)
	for v := range h.Messages() {
		got = append(got, v...)
	}
	assert.Equal(t, expected, got)
}
