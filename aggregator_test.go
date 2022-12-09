package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPartitionByTime(t *testing.T) {
	config := LiveChatAggregatorConfig{
		WindowType:     WindowTypeDuration,
		WindowDuration: 5 * time.Minute,
	}
	aggregator := NewLiveChatAggregator(config)

	t0 := time.Now()
	t1 := t0.Add(3 * time.Minute)
	t2 := t1.Add(3 * time.Minute)

	// ingest some messages
	aggregator.Ingest([]Message{
		{
			Sender:    "grassfan1",
			Content:   "grass",
			Timestamp: t0,
		},
		{
			Sender:    "grassfan2",
			Content:   "grass",
			Timestamp: t0,
		},
		{
			Sender:    "grassfan3",
			Content:   "grass",
			Timestamp: t1,
		},
		{
			Sender:    "waterfan1",
			Content:   "water",
			Timestamp: t1,
		},
		{
			Sender:    "waterfan2",
			Content:   "water",
			Timestamp: t1,
		},
		{
			Sender:    "firefan1",
			Content:   "fire",
			Timestamp: t1,
		},
	})

	// get counts
	assert.Equal(t, []MessageCount{
		{
			Content: "grass",
			Count:   3,
		},
		{
			Content: "water",
			Count:   2,
		},
		{
			Content: "fire",
			Count:   1,
		},
	}, aggregator.TopN(3))

	// ingest some more messages, pushing some messages out of the window
	aggregator.Ingest([]Message{
		{
			Sender:    "electricfan1",
			Content:   "electric",
			Timestamp: t2,
		},
		{
			Sender:    "electricfan2",
			Content:   "electric",
			Timestamp: t2,
		},
	})

	// get counts again
	// two grass votes are no longer in the 5 minute window
	// but two electric votes are now counted
	assert.Equal(t, []MessageCount{
		{
			Content: "water",
			Count:   2,
		},
		{
			Content: "electric",
			Count:   2,
		},
		{
			Content: "fire",
			Count:   1,
		},
		{
			Content: "grass",
			Count:   1,
		},
	}, aggregator.TopN(4))
}
