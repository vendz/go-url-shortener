package controllers

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/vendz/bitsnip/helper"
	"github.com/vendz/bitsnip/models"
)

func (databaseClient Database) ShortenUrl(c *fiber.Ctx) error {
	var payload models.Request
	var alias string

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	val, err := databaseClient.RedisClient.Get(context.Background(), c.IP()).Result()
	if err == redis.Nil {
		_ = databaseClient.RedisClient.Set(context.Background(), c.IP(), os.Getenv("RATE_LIMIT"), time.Hour).Err()
	} else {
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			limit, _ := databaseClient.RedisClient.TTL(context.Background(), c.IP()).Result()
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"status": "fail", "error": "Rate Limit Exceeded", "retry-after": limit / time.Nanosecond / time.Minute})
		}
	}

	if payload.Url == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "error": "Url is required"})
	}

	if payload.Alias == "" {
		for {
			alias = helper.Sug()
			val, _ := databaseClient.RedisClient.Get(context.Background(), alias).Result()
			if val == "" {
				break
			}
		}
	} else {
		val, _ := databaseClient.RedisClient.Get(context.Background(), payload.Alias).Result()
		if val != "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "error": "Alias is already taken"})
		}
		alias = payload.Alias
	}

	if payload.Expiry == 0 {
		payload.Expiry = 60 // days
	}

	err = databaseClient.RedisClient.Set(context.Background(), alias, payload.Url, payload.Expiry*time.Hour*24).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	res := models.Response{
		Url:                 payload.Url,
		CustomShort:         os.Getenv("DOMAIN") + "/" + alias,
		Expiry:              payload.Expiry,
		XrateLimitRemaining: 10,
		XrateLimitReset:     60,
	}

	databaseClient.RedisClient.Decr(context.Background(), c.IP()).Err()
	val, _ = databaseClient.RedisClient.Get(context.Background(), c.IP()).Result()
	res.XrateLimitRemaining, _ = strconv.Atoi(val)
	ttl, _ := databaseClient.RedisClient.TTL(context.Background(), c.IP()).Result()
	res.XrateLimitReset = ttl / time.Nanosecond / time.Minute
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": res})
}

func (databaseClient Database) GetUrl(c *fiber.Ctx) error {
	alias := c.Params("alias")
	val, err := databaseClient.RedisClient.Get(context.Background(), alias).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "error": "Url not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	// return c.Redirect(val, fiber.StatusMovedPermanently)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "url": val})
}

func (databaseClient Database) DeleteUrl(c *fiber.Ctx) error {
	alias := c.Params("alias")
	_, err := databaseClient.RedisClient.Get(context.Background(), alias).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "error": "Url not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	err = databaseClient.RedisClient.Del(context.Background(), alias).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}
