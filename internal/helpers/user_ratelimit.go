package helpers

import (
	"context"
	"crypto/rand"
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

	cmd := rdb.Eval(ctx, SLIDING_WINDOW_LUA, []string{key}, []string{strconv.FormatInt(time.Now().UnixMilli(), 10), strconv.FormatInt(window*1000, 10), strconv.FormatInt(max, 10), rand.Text()})

	if cmd.Err() != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "SERVER_ERROR",
		}), true
	}

	out, err := cmd.Int64Slice()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "SERVER_ERROR",
		}), true
	}

	ratelimited, retry_after := out[0], out[1]

	if ratelimited == 0 {
		c.Set("Retry-After", strconv.FormatInt(retry_after/1000, 10))
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "TOO_MANY_REQUESTS",
		}), true
	}

	return nil, false
}
