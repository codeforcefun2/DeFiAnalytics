package worker

import (
    "context"
    "encoding/json"
    "log"
    "time"

    "github.com/go-redis/redis/v8"
    "go.mongodb.org/mongo-driver/mongo"

    "github.com/yourusername/defi-analytics/internal/config"
)

// Transaction represents a blockchain transaction record.
type Transaction struct {
    Hash      string    `json:"hash"`
    From      string    `json:"from"`
    To        string    `json:"to"`
    Value     float64   `json:"value"`
    Timestamp time.Time `json:"timestamp"`
}

// NewRedisClient sets up and returns a new Redis client.
func NewRedisClient(addr, password string, db int) *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })
}

// StartWorker listens for new transactions on a Redis pub/sub channel and processes them.
func StartWorker(redisClient *redis.Client, mongoClient *mongo.Client, cfg *config.Config) {
    ctx := context.Background()
    pubsub := redisClient.Subscribe(ctx, "transactions")
    defer pubsub.Close()

    for {
        msg, err := pubsub.ReceiveMessage(ctx)
        if err != nil {
            log.Printf("error receiving message from Redis: %v", err)
            time.Sleep(time.Second)
            continue
        }
        var tx Transaction
        err = json.Unmarshal([]byte(msg.Payload), &tx)
        if err != nil {
            log.Printf("failed to unmarshal transaction: %v", err)
            continue
        }
        log.Printf("Processed transaction: %s", tx.Hash)

        // Process transaction asynchronously (for example, storing it in MongoDB)
        go storeTransaction(ctx, mongoClient, tx)
    }
}

func storeTransaction(ctx context.Context, client *mongo.Client, tx Transaction) {
    collection := client.Database("defianalytics").Collection("transactions")
    if _, err := collection.InsertOne(ctx, tx); err != nil {
        log.Printf("failed to store transaction %s: %v", tx.Hash, err)
    }
}
