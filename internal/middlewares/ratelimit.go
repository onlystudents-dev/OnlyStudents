package middlewares

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
)

func RateLimitMiddleware(c fiber.Ctx, rdb *redis.Client, scope string, max, window int64) error {
	ctx := context.Background()
	ip := c.IP()
	key := fmt.Sprintf("ratelimit:%s:%s", scope, ip)

	count, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if count == 1 {
		if err := rdb.Expire(ctx, key, time.Duration(window)*time.Second).Err(); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	if count > max {
		ttl, err := rdb.TTL(ctx, key).Result()
		if err == nil && ttl > 0 {
			c.Set("Retry-After", strconv.Itoa(int(ttl.Seconds())))
		}
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "TOO_MANY_REQUESTS",
		})
	}

	return c.Next()
}
