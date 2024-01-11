// Redis client
package db

import (
	"github.com/go-redis/redis/v8"
)

var (
	RedisClientLocal = "127.0.0.1:6379"
)
var RdsClientLocal *redis.Client

func InitRedisDB() {
	RdsClientLocal = redis.NewClient(&redis.Options{
		Addr:     RedisClientLocal,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

}
