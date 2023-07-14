package database

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	controllers "github.com/vendz/bitsnip/controllers"
)

func NewDatabase() controllers.Database {

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_DB_ADDR_LOCAL"),
		Password: "",
		DB:       0,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}
	fmt.Println("connected to Redis...")

	return controllers.Database{
		RedisClient: rdb,
	}
}
