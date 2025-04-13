package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/yourusername/defi-analytics/internal/config"
    "github.com/yourusername/defi-analytics/internal/data"
    "github.com/yourusername/defi-analytics/internal/websocket"
    "github.com/yourusername/defi-analytics/internal/worker"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("failed to load config: %v", err)
    }

    // Initialize MongoDB connection
    mongoClient, err := data.NewMongoClient(cfg.MongoURI)
    if err != nil {
        log.Fatalf("failed to connect to MongoDB: %v", err)
    }
    defer mongoClient.Disconnect(context.Background())

    // Initialize Redis client
    redisClient := worker.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
    defer redisClient.Close()

    // Start distributed worker in background
    go worker.StartWorker(redisClient, mongoClient, cfg)

    // Start WebSocket server
    wsServer := websocket.NewServer(cfg)
    http.HandleFunc("/ws", wsServer.HandleConnections)

    server := &http.Server{
        Addr:         cfg.ServerAddress,
        Handler:      nil,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    go func() {
        log.Printf("WebSocket server started at %s", cfg.ServerAddress)
        if err := server.ListenAndServe(); err != nil {
            log.Fatalf("failed to start server: %v", err)
        }
    }()

    // Graceful shutdown on signal interrupt
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()
    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("could not shutdown server gracefully: %v", err)
    }
    log.Println("Server stopped")
}
