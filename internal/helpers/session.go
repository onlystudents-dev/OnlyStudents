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
	Role        string `json:"role"`
	AccountID   int32  `json:"account_id"`
	AccountUUID string `json:"account_uuid"`
	DeviceID    string `json:"device_id"`
}

func accountSessionsKey(accountUUID string) string {
	return "account_sessions:" + accountUUID
}

func SessionCreate(ctx context.Context, rdb *redis.Client, accountID int32, accountUUID string, DeviceID string, role string) (string, error) {
	buf := make([]byte, GetUintEnvFallback("SESSION_COOKIE_LEN", 32))
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	secret := base64.RawStdEncoding.EncodeToString(buf)

	ttl := time.Duration(GetInt64EnvFallback("SESSION_TTL", 3600, 2592000)) * time.Second

	val, err := json.Marshal(SessionData{Role: role, AccountID: accountID, DeviceID: DeviceID, AccountUUID: accountUUID})

	if err != nil {
		return "", err
	}

	if err := rdb.Set(ctx, "session:"+secret, val, ttl).Err(); err != nil {
		return "", err
	}

	if err := rdb.SAdd(ctx, accountSessionsKey(accountUUID), secret).Err(); err != nil {
		return "", err
	}

	if err := rdb.Expire(ctx, accountSessionsKey(accountUUID), ttl).Err(); err != nil {
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

func SessionDelete(ctx context.Context, rdb *redis.Client, accountUUID string, secret string) error {
	if err := rdb.Del(ctx, "session:"+secret).Err(); err != nil {
		return err
	}

	if accountUUID == "" {
		return nil
	}

	return rdb.SRem(ctx, accountSessionsKey(accountUUID), secret).Err()
}

func SessionRevokeAll(ctx context.Context, rdb *redis.Client, accountUUID string) error {
	key := accountSessionsKey(accountUUID)

	secrets, err := rdb.SMembers(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	pipe := rdb.Pipeline()
	for _, secret := range secrets {
		pipe.Del(ctx, "session:"+secret)
	}
	pipe.Del(ctx, key)

	_, err = pipe.Exec(ctx)
	return err
}
