package controllers

import (
	"github.com/redis/go-redis/v9"
)

type Database struct {
	RedisClient *redis.Client
}
