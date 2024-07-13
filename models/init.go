package models

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB = Init()

var RDB = InitRedisDB()

func Init() *gorm.DB {
	dsn := "root:root@tcp(127.0.0.1:3306)/online_judge?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("gorm Init Error : %s", err)
	}
	return db
}

func InitRedisDB() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "47.121.184.191:6379",
		Password: "hxf@.2790", // no password set
		DB:       0,           // use default DB
	})
}
