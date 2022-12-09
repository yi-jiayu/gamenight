package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var (
	apiKey        string
	videoID       string
	liveChatID    string
	windowSeconds int
)

func init() {
	flag.StringVar(&apiKey, "apiKey", "", "API key")
	flag.StringVar(&videoID, "videoID", "", "Video ID")
	flag.StringVar(&liveChatID, "liveChatID", "", "Live chat ID")
	flag.IntVar(&windowSeconds, "windowSeconds", 60, "Sliding window duration in seconds")
}

func main() {
	flag.Parse()

	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "apiKey is required")
		os.Exit(1)
	}

	if videoID == "" && liveChatID == "" {
		fmt.Fprintln(os.Stderr, "one of videoID or liveChatID is required")
		os.Exit(1)
	}

	ctx := context.Background()
	yt, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		panic(err)
	}

	if liveChatID == "" {
		fmt.Fprintln(os.Stderr, "getting liveChatID from videoID")
		res, err := yt.Videos.List([]string{"liveStreamingDetails"}).Id(videoID).Do()
		if err != nil {
			panic(err)
		}
		liveChatID = res.Items[0].LiveStreamingDetails.ActiveLiveChatId
	}

	req := yt.LiveChatMessages.List(liveChatID, []string{"authorDetails", "snippet"})

	config := LiveChatAggregatorConfig{
		WindowDuration: time.Duration(windowSeconds) * time.Second,
	}
	aggregator := NewLiveChatAggregator(config)

	for {
		res, err := req.Do()
		if err != nil {
			panic(err)
		}

		var messages []Message
		for _, message := range res.Items {
			if message.Snippet.TextMessageDetails == nil {
				// ignore non-text messages
				continue
			}
			timestamp, err := time.Parse("2006-01-02T15:04:05.999999-07:00", message.Snippet.PublishedAt)
			if err != nil {
				log.Printf("ignoring invalid timestamp: %v", message.Snippet.PublishedAt)
				continue
			}
			messages = append(messages, Message{
				Sender:    message.AuthorDetails.ChannelId,
				Content:   message.Snippet.TextMessageDetails.MessageText,
				Timestamp: timestamp,
			})
		}
		aggregator.Ingest(messages)
		top10 := aggregator.TopN(10)
		fmt.Print("\033[H\033[2J")
		for _, message := range top10 {
			fmt.Printf("%s: %d\n", message.Content, message.Count)
		}
		fmt.Printf("window_size=%d", len(aggregator.window))

		time.Sleep(time.Duration(res.PollingIntervalMillis) * time.Millisecond)
		req = yt.LiveChatMessages.List(liveChatID, []string{"authorDetails", "snippet"}).PageToken(res.NextPageToken)
	}
}
