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
		// w.Write([]byte("response!"))
		ctx := context.Background()
		yts, err := youtube.NewService(ctx, option.WithAPIKey(k))
		if err != nil {
			// fmt.Println("failed to create service")
			w.WriteHeader((http.StatusBadRequest))
			return
		}

		call := yts.Channels.List([]string{"snippet,contentDetails,statistics"})
		response, err := call.Id(id).Do()
		if err != nil {
			// fmt.Println(err)
			w.WriteHeader((http.StatusBadRequest))
			return
		}

		var yt YoutubeStats
		fmt.Println(response.Items[0].Snippet.Title)
		if len(response.Items) > 0 {
			val := response.Items[0]
			yt := YoutubeStats{
				Subscriber:  int(val.Statistics.SubscriberCount),
				ChannelName: val.Snippet.Title,
				Views:       int(val.Statistics.ViewCount),
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(yt); err != nil {
				w.WriteHeader((http.StatusBadRequest))
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(yt); err != nil {
			// panic(err)
			w.WriteHeader((http.StatusBadRequest))
			return
		}
	}

}
