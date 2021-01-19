package kredis

import (
	"log"

	"github.com/go-redis/redis/v8"
)

var rClient *redis.Client

func InitClient(host, port string, db int) *redis.Client {
	rClient = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: "",
		DB:       db,
	})
	return rClient
}

func GetClient() *redis.Client {
	if rClient == nil {
		log.Println("Redis client is not initialized - First run InitClient")
	}
	return rClient
}
