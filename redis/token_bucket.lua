local tokens_key = KEYS[1]
local timestamp_key = KEYS[2]
local capacity = tonumber(ARGV[1])
local refill_per_ms = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local cost = tonumber(ARGV[4])
local tokens = tonumber(redis.call('GET', tokens_key) or capacity)
local updated = tonumber(redis.call('GET', timestamp_key) or now)
tokens = math.min(capacity, tokens + (now - updated) * refill_per_ms)
local allowed = tokens >= cost
if allowed then tokens = tokens - cost end
redis.call('SET', tokens_key, tokens, 'PX', 86400000)
redis.call('SET', timestamp_key, now, 'PX', 86400000)
return { allowed and 1 or 0, math.floor(tokens) }
