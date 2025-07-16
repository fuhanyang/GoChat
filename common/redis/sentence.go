package redis

import "github.com/gomodule/redigo/redis"

const ()

// 定义 Lua 脚本
var HmsetWithExpireScriptString = `
    if redis.call("EXISTS", KEYS[1]) == 0 then
        if #ARGV < 2 then
            return nil
        end
        redis.call("HMSET", KEYS[1], unpack(ARGV, 2, #ARGV))
        redis.call("EXPIRE", KEYS[1], tonumber(ARGV[1]))
    else
        redis.call("EXPIRE", KEYS[1], tonumber(ARGV[1]))
    end
    return 1
`
var HmsetWithExpireScript = redis.NewScript(1, HmsetWithExpireScriptString)

const (
	HGET    = "HGET"
	HSET    = "HSET"
	HMSET   = "HMSET"
	HMGET   = "HMGET"
	HDEL    = "HDEL"
	HEXISTS = "HEXISTS"
	HGETALL = "HGETALL"
	GET     = "GET"
	SET     = "SET"
	DEL     = "DEL"
	EXPIRE  = "EXPIRE"
	LPUSH   = "LPUSH"
	RPUSH   = "RPUSH"
	LPOP    = "LPOP"
	RPOP    = "RPOP"
	LTRIM   = "LTRIM"
	LLEN    = "LLEN"
	LRANGE  = "LRANGE"

	SADD            = "SADD"
	SREM            = "SREM"
	SRANDMEMBER     = "SRANDMEMBER"
	ONLINE_USER_SET = "online_user_set"
)
