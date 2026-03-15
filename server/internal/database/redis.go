package database

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/url"
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

	opts := &redis.Options{
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	configureRedisAddress(opts)

	RedisClient = redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		log.Fatal("Redis connection error:", err)
	}

	log.Println("Redis connected")
}

func configureRedisAddress(opts *redis.Options) {
	addrInput := strings.TrimSpace(config.Config.RedisAddr)
	password := strings.TrimSpace(config.Config.RedisPassword)

	if addrInput == "" {
		addrInput = strings.TrimSpace(config.Config.UpstashRedisRestURL)
	}
	if password == "" {
		password = strings.TrimSpace(config.Config.UpstashRedisRestToken)
	}

	if addrInput == "" {
		log.Fatal("Redis address is empty; set REDIS_ADDR (host:port, redis://, or rediss://)")
	}

	addr := addrInput
	useTLS := false

	if strings.Contains(addrInput, "://") {
		parsed, err := url.Parse(addrInput)
		if err != nil {
			log.Fatal("Invalid Redis address URL:", err)
		}

		switch parsed.Scheme {
		case "redis", "rediss", "tcp":
			addr = parsed.Host
			if parsed.Scheme == "rediss" {
				useTLS = true
			}

			if password == "" {
				if parsedPassword, ok := parsed.User.Password(); ok {
					password = parsedPassword
				}
			}
		case "http", "https":
			// Upstash REST URL fallback: derive Redis TCP endpoint from host and default TLS.
			addr = parsed.Hostname()
			if addr == "" {
				log.Fatal("Invalid Upstash REST URL: missing hostname")
			}
			if parsed.Port() == "" {
				addr = net.JoinHostPort(addr, "6379")
			}
			if parsed.Scheme == "https" {
				useTLS = true
			}
		default:
			log.Fatal("Unsupported Redis address scheme:", parsed.Scheme)
		}
	}

	if !strings.Contains(addr, ":") {
		addr = net.JoinHostPort(addr, "6379")
	}

	opts.Addr = addr
	opts.Password = password

	if useTLS {
		opts.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}
}
