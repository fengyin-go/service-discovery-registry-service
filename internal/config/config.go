// Package config 负责从环境变量加载服务配置。
package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config 服务运行配置。
type Config struct {
	Addr                   string
	MaxPageSize            int
	HeartbeatTimeoutSec    int
	HealthCheckIntervalSec int
}

// Load 从环境变量读取配置。
func Load() *Config {
	cfg := &Config{
		Addr:                   ":" + getenv("PORT", "8080"),
		MaxPageSize:            getenvInt("MAX_PAGE_SIZE", 100),
		HeartbeatTimeoutSec:    getenvInt("HEARTBEAT_TIMEOUT_SEC", 30),
		HealthCheckIntervalSec: getenvInt("HEALTH_CHECK_INTERVAL_SEC", 10),
	}
	if v := os.Getenv("ADDR"); v != "" {
		cfg.Addr = v
	}
	return cfg
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return def
	}
	return n
}

func (c *Config) String() string {
	return fmt.Sprintf("addr=%s max_page_size=%d heartbeat_timeout=%ds health_check_interval=%ds",
		c.Addr, c.MaxPageSize, c.HeartbeatTimeoutSec, c.HealthCheckIntervalSec)
}
