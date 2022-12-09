package main

import (
	"container/heap"
	"sort"
	"time"
)

type Message struct {
	Sender    string
	Content   string
	Timestamp time.Time
}

type WindowType int

const (
	WindowTypeDuration WindowType = iota
	WindowTypeCount
)

type LiveChatAggregatorConfig struct {
	Unique         bool
	Normalizer     func(string) string
	WindowType     WindowType
	WindowDuration time.Duration
	WindowSize     int
}

type LiveChatAggregator struct {
	config LiveChatAggregatorConfig

	// message counts
	counts map[string]int

	// map of senders to message
	votes map[string]string

	window []Message
}

type MessageCount struct {
	Content string
	Count   int
}

func NewLiveChatAggregator(config LiveChatAggregatorConfig) *LiveChatAggregator {
	return &LiveChatAggregator{
		config: config,
	}
}

func (a *LiveChatAggregator) Ingest(messages []Message) {
	// todo: optimize using sliding window algorithm
	// keep a running count here instead of when we call TopN
	// adjust the count based on incoming and outgoing elements
	end := messages[len(messages)-1].Timestamp
	cutoff := end.Add(-a.config.WindowDuration)
	a.window = append(a.window, messages...)
	i := sort.Search(len(a.window), func(i int) bool { return a.window[i].Timestamp.After(cutoff) })
	a.window = a.window[i:]
}

func (a *LiveChatAggregator) TopN(n int) []MessageCount {
	a.counts = make(map[string]int)
	for _, message := range a.window {
		a.counts[message.Content]++
	}
	h := &MessageHeap{}
	for content, count := range a.counts {
		heap.Push(h, MessageCount{Content: content, Count: count})
	}
	heap.Init(h)
	var top []MessageCount
	for i := 0; i < n; i++ {
		top = append(top, heap.Pop(h).(MessageCount))
	}
	return top
}
