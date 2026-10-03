package middlewares

import (
	"bytes"
	"context"
	"fmt"
	meapi "onlystudents/internal/api/v1/me"
	"strconv"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func FrontendMiddleware(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client, indexHTML []byte, scope string, max, window int64) error {
	ctx := context.Background()
	ip := c.IP()
	key := fmt.Sprintf("ratelimit:%s:%s", scope, ip)

	count, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	if count == 1 {
		if err := rdb.Expire(ctx, key, time.Duration(window)*time.Second).Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
		}
	}
	page := indexHTML

	if count > max {
		ttl, err := rdb.TTL(c.Context(), fmt.Sprintf("ratelimit:%s:%s", "api", c.IP())).Result()
		if err == nil && ttl > 0 {
			page = bytes.Replace(page, []byte(`"__INITIAL_STATUS__"`), []byte(strconv.Itoa(int(ttl.Seconds()))), 1)
		}
	} else if statusData, err := meapi.GetStatusData(c, pool, rdb); err == nil {
		payload, _ := json.Marshal(statusData)
		page = bytes.Replace(page, []byte(`"__INITIAL_STATUS__"`), payload, 1)
	}

	c.Set("Content-Type", "text/html; charset=utf-8")
	c.Set("Cache-Control", "no-store")
	return c.Send(page)
}
