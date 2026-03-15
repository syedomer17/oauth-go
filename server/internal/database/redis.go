package database

import (
	"context"
	"crypto/tls"
	"log"
	"oauth/internal/config"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var RedisNil = redis.Nil

func ConnectRedis() {
	db := 0
	if config.Config.RedisDB != "" {
		parsedDB, err := strconv.Atoi(config.Config.RedisDB)
		if err != nil {
			log.Fatal("Invalid REDIS_DB value:", err)
		}
		db = parsedDB
	}

	addr := config.Config.RedisAddr
	password := config.Config.RedisPassword
	if addr == "" {
		addr = config.Config.UpstashRedisRestURL
	}
	if password == "" {
		password = config.Config.UpstashRedisRestToken
	}

	opts := &redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	if strings.HasPrefix(addr, "rediss://") {
		opts.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	RedisClient = redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		log.Fatal("Redis connection error:", err)
	}

	log.Println("Redis connected")
}
