package helpers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
)

func UserRateLimit(c fiber.Ctx, rdb *redis.Client, scope string, max int64, window int64, user_id int32, role string) (error, bool) {
	switch role {
	case "student":
	case "guardian":
	case "teacher":
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"}), true
	}

	ctx := context.Background()
	key := fmt.Sprintf("user_ratelimit:%s:%s:%d", scope, role, user_id)

	count, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"}), true
	}

	if count == 1 {
		if err := rdb.Expire(ctx, key, time.Duration(window)*time.Second).Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"}), true
		}
	}

	if count > max {
		ttl, err := rdb.TTL(ctx, key).Result()
		if err == nil && ttl > 0 {
			c.Set("Retry-After", strconv.Itoa(int(ttl.Seconds())))
		}
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "TOO_MANY_REQUESTS",
		}), true
	}

	return nil, false
}
