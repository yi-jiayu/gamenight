// This example demonstrates an integer heap built using the heap interface.
package main

type MessageHeap []MessageCount

func (h MessageHeap) Len() int           { return len(h) }
func (h MessageHeap) Less(i, j int) bool { return h[i].Count > h[j].Count }
func (h MessageHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MessageHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(MessageCount))
}

func (h *MessageHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
