package redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	"databridge/common"

	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel"
)

type RedisConfig struct {
	Addr           string `yaml:"addr"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	DB             int    `yaml:"db"`
	ClusterEnable  bool   `yaml:"cluster_enable"`
	SentinelEnable bool   `yaml:"sentinel_enable"`
	SentinelHosts  string `yaml:"sentinel_hosts"`
	SentinelPort   int    `yaml:"sentinel_port"`
	MasterName     string `yaml:"master_name"`
	DefaultQueue   string `yaml:"default_queue"`
}

// InitRedis InitRedis
func InitRedis(cfg *RedisConfig) any {
	tracer := otel.Tracer(common.ModuleName)
	_, span := tracer.Start(context.Background(), "init redis")
	fmt.Printf("trace_id:%s,span_id:%s\n",
		span.SpanContext().TraceID(), span.SpanContext().SpanID())
	defer span.End()
	if cfg.SentinelEnable {
		return initFailoverClient(cfg)
	} else if cfg.ClusterEnable {
		return initClusterClient(cfg)
	} else {
		return initClient(cfg)
	}
}

func initClusterClient(cfg *RedisConfig) *redis.ClusterClient {
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        []string{},
		Username:     cfg.User,
		Password:     cfg.Password,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		PoolTimeout:  5 * time.Second,
	})
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
	})
}
