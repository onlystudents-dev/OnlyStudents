package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"onlystudents/internal/helpers"
	"time"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
)

type SessionData struct {
	Role      string `json:"role"`
	AccountID int32  `json:"account_id"`
}

type SessionStore struct {
	RedisDB *redis.Client
}

func (s *SessionStore) Create(ctx context.Context, accountID int32, role string) (string, error) {
	buf := make([]byte, helpers.GetUintEnvFallback("SESSION_COOKIE_LEN", 32))
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	secret := base64.RawStdEncoding.EncodeToString(buf)

	ttl := time.Duration(helpers.GetUintEnvFallback("SESSION_TTL", 3600)) * time.Second

	val, err := json.Marshal(SessionData{Role: role, AccountID: accountID})

	if err != nil {
		return "", err
	}

	if err := s.RedisDB.Set(ctx, "session:"+secret, val, ttl).Err(); err != nil {
		return "", err
	}

	return secret, nil
}

func (s *SessionStore) Get(ctx context.Context, secret string) (SessionData, error) {
	val, err := s.RedisDB.Get(ctx, "session:"+secret).Bytes()
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

func (s *SessionStore) Delete(ctx context.Context, secret string) error {
	return s.RedisDB.Del(ctx, "session:"+secret).Err()
}
