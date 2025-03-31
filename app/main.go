package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func newRouter(ytApiKey string, ytId string) *httprouter.Router {
	mux := httprouter.New()

	if ytApiKey == "" {
		log.Fatal("youtube API key no provider")
	}

	mux.GET("/youtube/channel/stats", getChannelStats(ytApiKey, ytId))

	return mux
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("❌ Could not load .env file:", err)
	} else {
		log.Println("✅ .env file loaded successfully")
	}

	ytApiKey := os.Getenv("YOUTUBE_API_KEY")
	ytId := os.Getenv("YOUTUBE_ID")
	fmt.Println("API KEY FROM ENV:", ytApiKey)
	fmt.Println("YOUTUBE USER ID:", ytId)

	srv := &http.Server{
		Addr:    ":10101",
		Handler: newRouter(ytApiKey, ytId),
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint

		log.Println("service interupt received")

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("http server shutdown error: %v", err)
		}

		log.Println("shutdown complete")
		close(idleConnsClosed)

	}()

	log.Print("Starting server on port 10101")
	if err := srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("fatal http server failed to start: %v", err)
		}
	}

	<-idleConnsClosed
	log.Println("Service Stop")

}
