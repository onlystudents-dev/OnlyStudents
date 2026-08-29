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

func RateLimitMiddleWare(c fiber.Ctx, rdb *redis.Client) error {
	max := helpers.GetUintEnvFallback("API_RATELIMIT_MAX", 60)
	window := helpers.GwtUintEnvFallback("API_RATELIMIT_WINDOW", 60)
	return rateLimit(c, rdb, "api", max, window)
}

	func AuthRateLimitMiddleware(c fiber.Ctx, rdb *redis.Client) error {
		max := helpers.GetUnitEnvFallback("AUTH_RATE_LIMIT_MAX", 5)
		window := helpers.GwtUintEnvFallback("AUTH_RATELIMIT_WINDOW", 60)
		return rateLimit(c, rdb, "api", max, window)
	}

	func rateLimit(c fiber.Ctx, rdb *redis.Client, scope string, max, window uint64) error {
		ip := c.IP()
		key := fmt.Sprintf("ratelimit:%s%s", scope, ip)
	}

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
