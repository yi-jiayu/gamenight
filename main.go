package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var (
	apiKey     string
	videoID    string
	liveChatID string
)

func init() {
	flag.StringVar(&apiKey, "apiKey", "", "API key")
	flag.StringVar(&videoID, "videoID", "", "Video ID")
	flag.StringVar(&liveChatID, "liveChatID", "", "Live chat ID")
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

	req := yt.LiveChatMessages.List(liveChatID, []string{"snippet"})
	for {
		res, err := req.Do()
		if err != nil {
			panic(err)
		}
		for _, message := range res.Items {
			fmt.Println(message.Snippet.TextMessageDetails.MessageText)
		}
		fmt.Fprintf(os.Stderr, "wait %d ms\n", res.PollingIntervalMillis)
		time.Sleep(time.Duration(res.PollingIntervalMillis) * time.Millisecond)
		req = yt.LiveChatMessages.List(liveChatID, []string{"snippet"}).PageToken(res.NextPageToken)
	}
}
