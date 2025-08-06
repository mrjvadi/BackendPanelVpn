package cache

import (
	"context"
	"encoding/json"
	"fmt"
	type_redis "github.com/mrjvadi/BackendPanelVpn/types/type-redis"
	"github.com/redis/go-redis/v9"
	"time"
)

// SetSessionData سشن را به صورت JSON ذخیره می‌کند
func (r *RedisCache) SetSessionData(ctx context.Context, session *type_redis.SessionData, ttl time.Duration) error {
	key := fmt.Sprintf("session_data:%s", session.Token)
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("marshal session data: %w", err)
	}
	if err := r.Client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("redis set session data: %w", err)
	}
	return nil
}

// GetSessionData بر اساس توکن، سشن را برمی‌گرداند
func (r *RedisCache) GetSessionData(ctx context.Context, token string) (*type_redis.SessionData, error) {
	key := fmt.Sprintf("session_data:%s", token)
	raw, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("redis get session data: %w", err)
	}

	var session type_redis.SessionData
	if err := json.Unmarshal([]byte(raw), &session); err != nil {
		return nil, fmt.Errorf("unmarshal session data: %w", err)
	}
	return &session, nil
}

// DeleteSessionData سشن را از ردیس حذف می‌کند
func (r *RedisCache) DeleteSessionData(ctx context.Context, token string) error {
	key := fmt.Sprintf("session_data:%s", token)
	if err := r.Client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis delete session data: %w", err)
	}
	return nil
}

// (اختیاری) Index بر اساس username برای lookup
func (r *RedisCache) SetUserIndex(ctx context.Context, username, token string, ttl time.Duration) error {
	key := fmt.Sprintf("user_index:%s", username)
	if err := r.Client.Set(ctx, key, token, ttl).Err(); err != nil {
		return fmt.Errorf("redis set user index: %w", err)
	}
	return nil
}

func (r *RedisCache) GetUserIndex(ctx context.Context, username string) (string, error) {
	key := fmt.Sprintf("user_index:%s", username)
	token, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", fmt.Errorf("redis get user index: %w", err)
	}
	return token, nil
}

func (r *RedisCache) DeleteUserIndex(ctx context.Context, username string) error {
	key := fmt.Sprintf("user_index:%s", username)
	if err := r.Client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis delete user index: %w", err)
	}
	return nil
}
