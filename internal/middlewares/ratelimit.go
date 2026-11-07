package middlewares

import (
	"context"
	"crypto/rand"
	"fmt"
	"onlystudents/internal/helpers"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
)

func RateLimitMiddleware(c fiber.Ctx, rdb *redis.Client, scope string, max, window int64) error {
	ctx := context.Background()
	ip := c.IP()
	key := fmt.Sprintf("ip_ratelimit:%s:%s", scope, ip)

	cmd := rdb.Eval(ctx, helpers.SLIDING_WINDOW_LUA, []string{key}, []string{strconv.FormatInt(time.Now().UnixMilli(), 10), strconv.FormatInt(window*1000, 10), strconv.FormatInt(max, 10), rand.Text()})

	if cmd.Err() != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "SERVER_ERROR",
		})
	}

	out, err := cmd.Int64Slice()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "SERVER_ERROR",
		})
	}

	ratelimited, retry_after := out[0], out[1]

	if ratelimited == 0 {
		c.Set("Retry-After", strconv.FormatInt(retry_after/1000, 10))
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "TOO_MANY_REQUESTS",
		})
	}

	return c.Next()
}
