package redisClient

import (
	"context"
	"time"

	"service/utils"

	"github.com/go-redis/redis/v8"
)

var (
	ctx = context.Background()
	Rdb = redis.NewClient(&redis.Options{
		Addr:     utils.GetEnv("REDIS_IP"),
		Password: utils.GetEnv("REDIS_PWD"),
		DB:       0, // use default DB
	})
	RedisKeyLimitGlobal             = "GlobalInfo"
	RedisKeyPicbot                  = "Picbot:"
	RedisKeyPicbotAlive             = RedisKeyPicbot + "Alive:" // + picbot name
	RedisKeyPicbotState             = RedisKeyPicbot + "State:" // + picbot name
	RedisKeySocketWaitReplyMessages = "SocketWaitReplyMessages:"
	RedisKeyToken                   = "Token:"
	RedisKeyMember                  = "MemberId:"
	GenerateIdKey                   = "GenerateIdKey"
)

type GlobalInfo struct {
	SendPicbotProductErrorEmailCount int
}

func InitRedis() bool {
	utils.PrintObj(utils.GetEnv("REDIS_IP"), "InitRedis")
	ticker := time.NewTicker(1 * time.Second)
	count := 0
	success := false

	for range ticker.C {
		initRedisClient()

		if err := PingRedis(); err == nil {
			success = true
		}

		utils.PrintObj("InitRedis retry:"+utils.ToString(count), "")

		if count == 60 || success {
			return success
		}

		count++
	}

	return true
}

func initRedisClient() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     utils.GetEnv("REDIS_IP"),
		Password: utils.GetEnv("REDIS_PWD"),
		DB:       0, // use default DB
	})
}

func PingRedis() error {
	_, err := Rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}

	return nil
}

func SetRedis(key string, value string, hours int16) error {
	time := time.Duration(hours) * time.Hour
	err := Rdb.Set(ctx, key, value, time).Err()
	utils.PrintObj([]string{key, utils.FixStringLen(value)}, "SetRedis")

	if err != nil {
		utils.PrintObj(err, "SetRedis err")
	}

	return err
}

func SetRedisForever(key string, value string) error {
	err := Rdb.Set(ctx, key, value, -1).Err()
	utils.PrintObj([]string{key, utils.FixStringLen(value)}, "SetRedis")

	if err != nil {
		utils.PrintObj(err, "SetRedis err")
	}

	return err
}

func GetRedis(key string) string {
	value, err := Rdb.Get(ctx, key).Result()
	utils.PrintObj([]string{key, utils.FixStringLen(value)}, "GetRedis")

	if err != nil {
		utils.PrintObj(err, "GetRedis err")
		return ""
	}

	return value
}

func GetRedisKeys(prefix string) []string {
	var cursor uint64
	res := []string{}

	for {
		var keys []string
		var err error
		keys, cursor, err = Rdb.Scan(ctx, cursor, prefix, 0).Result()
		if err != nil {
			panic(err)
		}

		for _, key := range keys {
			res = append(res, key)
		}

		if cursor == 0 { // no more keys
			break
		}
	}

	utils.PrintObj(res, "GetRedisKeys")

	return res
}

func HasRedis(key string) bool {
	val, err := Rdb.Exists(ctx, key).Result()
	utils.PrintObj([]interface{}{key, val}, "HasRedis")

	if err != nil || val == 0 {
		utils.PrintObj("redis not found", "")
		return false
	}

	return true
}

func DeleteRedis(key string) {
	utils.PrintObj(key, "DeleteRedis")
	_, err := Rdb.Del(ctx, key).Result()

	if err != nil {
		utils.PrintObj(err, "DeleteRedis err")
	}
}
