package helpers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
)

type SessionData struct {
	Role      string `json:"role"`
	AccountID int32  `json:"account_id"`
}

func SessionCreate(ctx context.Context, rdb *redis.Client, accountID int32, role string) (string, error) {
	buf := make([]byte, GetUintEnvFallback("SESSION_COOKIE_LEN", 32))
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	secret := base64.RawStdEncoding.EncodeToString(buf)

	ttl := time.Duration(GetInt64EnvFallback("SESSION_TTL", 3600, 2592000)) * time.Second

	val, err := json.Marshal(SessionData{Role: role, AccountID: accountID})

	if err != nil {
		return "", err
	}

	if err := rdb.Set(ctx, "session:"+secret, val, ttl).Err(); err != nil {
		return "", err
	}

	return secret, nil
}

func SessionGet(ctx context.Context, rdb *redis.Client, secret string) (SessionData, error) {
	val, err := rdb.Get(ctx, "session:"+secret).Bytes()
	if errors.Is(err, redis.Nil) {
		return SessionData{}, errors.New("session not found")
	}

	if err != nil {
		return SessionData{}, err
	}

	var data SessionData
	if err := json.Unmarshal(val, &data); err != nil {
		return SessionData{}, err
	}

	return data, nil
}

func SessionDelete(ctx context.Context, rdb *redis.Client, secret string) error {
	return rdb.Del(ctx, "session:"+secret).Err()
}
