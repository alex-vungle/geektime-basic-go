local zsetName = KEYS[1]
local member2incr = ARGV[1]
local newScore = ARGV[2]

local exist = redis.call("EXISTS", zsetName)

if exist == 0 then
    redis.call("ZADD", zsetName, newScore, member2incr)
end

local currentScore = redis.call("ZSCORE", zsetName, member2incr)

if currentScore then
    local nweScore = newScore +1
    redis.call("ZADD", zsetName, newScore, member2incr)
else
     redis.call("ZADD", zsetName, newScore, member2incr)
end