package caching

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client *redis.Client
}

func NewRateLimiter(client *redis.Client) *RateLimiter {
	return &RateLimiter{client: client}
}

func (rl *RateLimiter) RateLimiting(ctx context.Context, ip string, windowMinutes int) (int, error) {
	key := "rate_limit:" + ip
	count, err := rl.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment rate limit for ip %s: %w", ip, err)
	}
	if count == 1 {
		err := rl.client.Expire(ctx, key, time.Duration(windowMinutes)*time.Minute).Err()
		if err != nil {
			return 0, fmt.Errorf("failed to set rate limit expiration for ip %s: %w", ip, err)
		}
	}
	return int(count), nil
}
