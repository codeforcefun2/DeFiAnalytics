package data

import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoClient sets up a MongoDB client.
func NewMongoClient(uri string) (*mongo.Client, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        return nil, err
    }
    // Verify connection with a ping
    if err := client.Ping(ctx, nil); err != nil {
        return nil, err
    }
    return client, nil
}

// Additional helper functions (such as InsertTransaction, QueryAnalytics, etc.) can be added here.
