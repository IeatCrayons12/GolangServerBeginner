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

	mux.GET("/youtube/channel/stats", getChannelStats(ytApiKey, ytId))
	return mux
}

func main() {
	// ✅ Load .env only in dev (not production/docker)
	if os.Getenv("GO_ENV") != "production" {
		if err := godotenv.Load(".env"); err != nil {
			log.Println("⚠️  .env file not found, continuing without it")
		} else {
			log.Println("✅ .env file loaded successfully")
		}
	}

	ytApiKey := os.Getenv("YOUTUBE_API_KEY")
	ytId := os.Getenv("YOUTUBE_ID")

	if ytApiKey == "" || ytId == "" {
		log.Fatal("❌ Missing required environment variables. Make sure YOUTUBE_API_KEY and YOUTUBE_ID are set.")
	}

	fmt.Println("📺 API KEY FROM ENV:", ytApiKey)
	fmt.Println("📺 YOUTUBE USER ID:", ytId)

	srv := &http.Server{
		Addr:    ":10101",
		Handler: newRouter(ytApiKey, ytId),
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("🚦 Shutdown signal received, exiting...")

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("❌ HTTP server Shutdown Error: %v", err)
		}

		log.Println("✅ Shutdown complete")
		close(idleConnsClosed)
	}()

	log.Println("🚀 Starting server on port 10101")
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("❌ Fatal server error: %v", err)
	}

	<-idleConnsClosed
	log.Println("👋 Server stopped")
}
