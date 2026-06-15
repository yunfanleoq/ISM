package packages

import (
	"context"
	"reflect"

	"github.com/mattn/anko/env"
	redis "github.com/redis/go-redis/v9"
)

func RedisGet(Options *redis.Options, key string) (string, error) {
	var ctx = context.Background()
	newOption := redis.Options{}
	newOption.Addr = Options.Addr
	newOption.Password = Options.Password
	newOption.DB = Options.DB
	rdb := redis.NewClient(&newOption)
	return rdb.Get(ctx, key).Result()
}
func RedisRPush(Options *redis.Options, key string, values ...interface{}) error {
	var ctx = context.Background()
	newOption := redis.Options{}
	newOption.Addr = Options.Addr
	newOption.Password = Options.Password
	newOption.DB = Options.DB
	rdb := redis.NewClient(&newOption)
	return rdb.RPush(ctx, key, values).Err()
}

func RedisLPop(Options *redis.Options, key string) (string, error) {
	var ctx = context.Background()
	newOption := redis.Options{}
	newOption.Addr = Options.Addr
	newOption.Password = Options.Password
	newOption.DB = Options.DB
	rdb := redis.NewClient(&newOption)
	return rdb.LPop(ctx, key).Result()
}

func RedisPublish(Options *redis.Options, key string, values interface{}) error {
	var ctx = context.Background()
	newOption := redis.Options{}
	newOption.Addr = Options.Addr
	newOption.Password = Options.Password
	newOption.DB = Options.DB
	rdb := redis.NewClient(&newOption)
	return rdb.Publish(ctx, key, values).Err()
}

func RedisSubscribe(Options *redis.Options, key string, values interface{}) *redis.PubSub {
	var ctx = context.Background()
	newOption := redis.Options{}
	newOption.Addr = Options.Addr
	newOption.Password = Options.Password
	newOption.DB = Options.DB
	rdb := redis.NewClient(&newOption)
	return rdb.Subscribe(ctx, key)
}

func RedisSet(Options *redis.Options, key string, value interface{}) error {
	var ctx = context.Background()
	newOption := redis.Options{}
	newOption.Addr = Options.Addr
	newOption.Password = Options.Password
	newOption.DB = Options.DB
	rdb := redis.NewClient(&newOption)

	return rdb.Set(ctx, key, value, 0).Err()
}

func init() {
	env.Packages["redis"] = map[string]reflect.Value{
		"NewClient":  reflect.ValueOf(redis.NewClient),
		"RedisGet":   reflect.ValueOf(RedisGet),
		"RedisSet":   reflect.ValueOf(RedisSet),
		"RedisRPush": reflect.ValueOf(RedisRPush),
		"RedisLPop":  reflect.ValueOf(RedisLPop),

		"RedisPublish":   reflect.ValueOf(RedisPublish),
		"RedisSubscribe": reflect.ValueOf(RedisSubscribe),
	}

	env.PackageTypes["redis"] = map[string]reflect.Type{
		"Options": reflect.TypeOf(redis.Options{}),
	}
}
