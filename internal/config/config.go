package config

import (
    "fmt"
    "os"
)

type Config struct {
    ServerAddress string
    MongoURI      string
    RedisAddr     string
    RedisPassword string
    RedisDB       int
}

func LoadConfig() (*Config, error) {
    return &Config{
        ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
        MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
        RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
        RedisPassword: getEnv("REDIS_PASSWORD", ""),
        RedisDB:       getEnvInt("REDIS_DB", 0),
    }, nil
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value, exists := os.LookupEnv(key); exists {
        var v int
        fmt.Sscanf(value, "%d", &v)
        return v
    }
    return defaultValue
}
