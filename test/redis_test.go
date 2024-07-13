package test

import (
	"context"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
	Addr:     "47.121.184.191:6379",
	Password: "hxf@.2790", // no password set
	DB:       0,           // use default DB
})

func TestRedisGet(t *testing.T) {
	get, _ := rdb.Get(ctx, "name").Result()
	println(get)
}

func TestRedisSet(t *testing.T) {
	rdb.Set(ctx, "name", "va", time.Second*1630)
}
