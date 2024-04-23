package redis

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Addr           string `yaml:"addr"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	DB             int    `yaml:"db"`
	SentinelEnable bool   `yaml:"sentinel_enable"`
	SentinelHosts  string `yaml:"sentinel_hosts"`
	SentinelPort   int    `yaml:"sentinel_port"`
	MasterName     string `yaml:"master_name"`
	DefaultQueue   string `yaml:"default_queue"`
}

// InitRedis InitRedis
func InitRedis(cfg *RedisConfig) *redis.Client {
	if cfg.SentinelEnable {
		return initFailoverClient(cfg)
	} else {
		return initClient(cfg)
	}
}

func initClient(cfg *RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Username:     cfg.User,
		Password:     cfg.Password,
		DB:           cfg.DB,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		PoolTimeout:  5 * time.Second,
		MaxConnAge:   0,
	})
}

func initFailoverClient(cfg *RedisConfig) *redis.Client {
	sentinelAddr := strings.Split(cfg.SentinelHosts, ",")
	for i, addr := range sentinelAddr {
		sentinelAddr[i] = fmt.Sprintf("%s:%d", addr, cfg.SentinelPort)
	}
	return redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    cfg.MasterName,
		SentinelAddrs: sentinelAddr,
		Username:      cfg.User,
		Password:      cfg.Password,
		DB:            cfg.DB,
		MinIdleConns:  5,
		DialTimeout:   5 * time.Second,
		WriteTimeout:  5 * time.Second,
		PoolTimeout:   5 * time.Second,
		MaxConnAge:    0,
	})
}
