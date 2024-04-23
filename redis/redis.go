package redis

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

type RedisConfig struct {
	Addr           string `yaml:"addr"`
	Password       string `yaml:"password"`
	DB             int    `yaml:"db"`
	SentinelEnable bool   `yaml:"sentinel_enable"`
	SentinelHosts  string `yaml:"sentinel_hosts"`
	SentinelPort   int    `yaml:"sentinel_port"`
	MasterName     string `yaml:"master_name"`

	// to use
	User          string `yaml:"user"`
	DefaultQueue  string `yaml:"default_queue"`
	ResultBackend string `yaml:"result_backend"`
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
		Password:      cfg.Password,
		DB:            cfg.DB,
		MinIdleConns:  5,
		DialTimeout:   5 * time.Second,
		WriteTimeout:  5 * time.Second,
		PoolTimeout:   5 * time.Second,
		MaxConnAge:    0,
	})
}
