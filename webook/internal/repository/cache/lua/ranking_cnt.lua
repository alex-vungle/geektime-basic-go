local zsetName = KEY[1]
local member2incr = ARGV[1]

local exist = redis.call("EXISTS", zsetName)
if exist == 1 then
    newScore = redis.call("ZINCRBY", zsetName, 1, member2incr)
    return 1
else
    return 0
end