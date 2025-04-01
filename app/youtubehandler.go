package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YoutubeStats struct {
	Subscriber     int    `json:"subscribers"`
	ChannelName    string `json:"channelName"`
	MinutesWatched int    `json:"minutesWatched"`
	Views          int    `json:"views"`
}

func getChannelStats(k string, id string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		fmt.Println("🔑 API Key:", k)
		fmt.Println("📺 Channel ID:", id)

		ctx := context.Background()
		yts, err := youtube.NewService(ctx, option.WithAPIKey(k))
		if err != nil {
			fmt.Println("❌ Failed to create YouTube service:", err)
			http.Error(w, "Failed to connect to YouTube service", http.StatusInternalServerError)
			return
		}

		call := yts.Channels.List([]string{"snippet,contentDetails,statistics"})
		response, err := call.Id(id).Do()
		if err != nil {
			fmt.Println("❌ YouTube API error:", err)
			http.Error(w, "Failed to fetch channel stats", http.StatusBadRequest)
			return
		}

		if len(response.Items) == 0 {
			fmt.Println("❌ No channel found with that ID")
			http.Error(w, "Channel not found", http.StatusNotFound)
			return
		}

		val := response.Items[0]
		fmt.Println("✅ Channel Found:", val.Snippet.Title)

		yt := YoutubeStats{
			Subscriber:  int(val.Statistics.SubscriberCount),
			ChannelName: val.Snippet.Title,
			Views:       int(val.Statistics.ViewCount),
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(yt); err != nil {
			fmt.Println("❌ Failed to encode response:", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
