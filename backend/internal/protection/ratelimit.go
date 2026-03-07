package protection

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"shieldpanel/backend/internal/config"
)

var tokenBucketScript = redis.NewScript(`
local key = KEYS[1]
local rate = tonumber(ARGV[1])
local burst = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local requested = tonumber(ARGV[4])

local data = redis.call("HMGET", key, "tokens", "ts")
local tokens = tonumber(data[1])
local ts = tonumber(data[2])
if tokens == nil then tokens = burst end
if ts == nil then ts = now end
local delta = math.max(0, now - ts)
tokens = math.min(burst, tokens + delta * rate)
local allowed = 0
if tokens >= requested then
  tokens = tokens - requested
  allowed = 1
end
redis.call("HMSET", key, "tokens", tokens, "ts", now)
redis.call("EXPIRE", key, math.ceil((burst / rate) * 4))
return {allowed, tokens}
`)

func allowByTokenBucket(ctx context.Context, redisClient *redis.Client, cfg config.Config, domainName string, clientIP string, rate int, burst int) (bool, error) {
	key := fmt.Sprintf("%srate:%s:%s", cfg.Redis.KeyPrefix, domainName, clientIP)
	now := float64(time.Now().UnixNano()) / float64(time.Second)
	result, err := tokenBucketScript.Run(ctx, redisClient, []string{key}, rate, burst, now, 1).Result()
	if err != nil {
		return false, err
	}
	values, ok := result.([]any)
	if !ok || len(values) == 0 {
		return false, nil
	}
	allowed, ok := values[0].(int64)
	if !ok {
		return false, nil
	}
	return allowed == 1, nil
}
