package middlewares

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
)

func RateLimitMiddleware(c fiber.Ctx, rdb *redis.Client) error {
	max := helpers.GetUintEnvFallback("AUTH_RATE_LIMIT_MAX", 5)
	window := helpers.GetUintEnvFallback("AUTH_RATE_LIMIT_WINDOW", 60) // seconds

	ip := c.IP()
	key := fmt.Sprintf("ratelimit:auth:%s", ip)

	ctx := context.Background()

	count, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if count == 1 {
		if err := rdb.Expire(ctx, key, time.Duration(window)*time.Second).Err(); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	if count > int64(max) {
		ttl, err := rdb.TTL(ctx, key).Result()
		if err == nil && ttl > 0 {
			c.Set("Retry-After", strconv.Itoa(int(ttl.Seconds())))
		}
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "too many requests, please try again later",
		})
	}

	return c.Next()
}
