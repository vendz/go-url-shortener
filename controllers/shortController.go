package controllers

import (
	"context"
	"os"
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

	err := databaseClient.RedisClient.Set(context.Background(), alias, payload.Url, payload.Expiry*time.Hour*24).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	res := models.Response{
		Url:         payload.Url,
		CustomShort: os.Getenv("DOMAIN") + "/" + alias,
		Expiry:      payload.Expiry,
	}

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
