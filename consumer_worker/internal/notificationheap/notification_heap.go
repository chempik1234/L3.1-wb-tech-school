package notificationheap

import (
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/models"
	"time"
)

// NotificationHeap is a slice for sorting incoming notifications
//
// Sorts by PublicationAt (Desc)
//
// { late, ...,  early }
//
// Implements heap.Interface
type NotificationHeap []*models.Notification

// Len returns slice length
//
// required by heap.Interface
func (h *NotificationHeap) Len() int { return len(*h) }

// Less returns if i < j
//
// required by heap.Interface
func (h *NotificationHeap) Less(i, j int) bool {
	timeI, _ := time.Parse(time.RFC3339, (*h)[i].PublicationAt.String())
	timeJ, _ := time.Parse(time.RFC3339, (*h)[j].PublicationAt.String())
	return timeI.After(timeJ)
}

// Swap 2 elements by indices
//
// required by heap.Interface
func (h *NotificationHeap) Swap(i, j int) { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

// Push simply adds a new element to the container
//
// required by heap.Interface
func (h *NotificationHeap) Push(x interface{}) {
	*h = append(*h, x.(*models.Notification))
}

// Pop removes the element with earliest PublicationAt and returns it
//
// Simple slices the slice without allocating a new one
// because space been once occupied is supposed to be filled and freed regularly
//
// required by heap.Interface
func (h *NotificationHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// Peek returns (not pops) the element to be popped
//
// returns nil if Len = 0
func (h *NotificationHeap) Peek() *models.Notification {
	if h.Len() == 0 {
		return nil
	}
	n := len(*h)
	return (*h)[n-1]
}
