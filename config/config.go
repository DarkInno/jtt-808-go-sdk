package config

import (
	"sync"
	"time"
)

// Config 配置管理器
type Config struct {
	mu       sync.RWMutex
	values   map[string]interface{}
	watchers map[string][]func(interface{})
}

// NewConfig 创建配置管理器
func NewConfig() *Config {
	return &Config{
		values:   make(map[string]interface{}),
		watchers: make(map[string][]func(interface{})),
	}
}

// GetString 获取字符串配置
func (c *Config) GetString(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if val, ok := c.values[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// GetInt 获取整数配置
func (c *Config) GetInt(key string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if val, ok := c.values[key]; ok {
		if i, ok := val.(int); ok {
			return i
		}
	}
	return 0
}

// GetBool 获取布尔配置
func (c *Config) GetBool(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if val, ok := c.values[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}

// GetDuration 获取时间间隔配置
func (c *Config) GetDuration(key string) time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if val, ok := c.values[key]; ok {
		if d, ok := val.(time.Duration); ok {
			return d
		}
	}
	return 0
}

// Set 设置配置
func (c *Config) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.values[key] = value

	// 通知观察者
	if watchers, ok := c.watchers[key]; ok {
		for _, watcher := range watchers {
			go watcher(value)
		}
	}
}

// Unmarshal 反序列化配置
func (c *Config) Unmarshal(target interface{}) error {
	// 简化实现，实际应该使用反射或JSON序列化
	return nil
}

// Watch 监听配置变化
func (c *Config) Watch(key string, callback func(interface{})) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.watchers[key] = append(c.watchers[key], callback)
}

// LoadFromFile 从文件加载配置
func (c *Config) LoadFromFile(path string) error {
	// 简化实现，实际应该读取YAML/JSON/TOML文件
	return nil
}

// LoadFromEnv 从环境变量加载配置
func (c *Config) LoadFromEnv() error {
	// 简化实现，实际应该读取环境变量
	return nil
}

// ServerConfig 服务器配置
type ServerConfig struct {
	// ListenAddr 监听地址
	ListenAddr string `json:"listen_addr" yaml:"listen_addr"`
	// MaxConnections 最大连接数
	MaxConnections int `json:"max_connections" yaml:"max_connections"`
	// ReadTimeout 读超时
	ReadTimeout time.Duration `json:"read_timeout" yaml:"read_timeout"`
	// WriteTimeout 写超时
	WriteTimeout time.Duration `json:"write_timeout" yaml:"write_timeout"`
	// IdleTimeout 空闲超时
	IdleTimeout time.Duration `json:"idle_timeout" yaml:"idle_timeout"`
	// ReadBufferSize 读缓冲区大小
	ReadBufferSize int `json:"read_buffer_size" yaml:"read_buffer_size"`
	// WriteBufferSize 写缓冲区大小
	WriteBufferSize int `json:"write_buffer_size" yaml:"write_buffer_size"`
	// WorkerMinWorkers 最小Worker数
	WorkerMinWorkers int `json:"worker_min_workers" yaml:"worker_min_workers"`
	// WorkerMaxWorkers 最大Worker数
	WorkerMaxWorkers int `json:"worker_max_workers" yaml:"worker_max_workers"`
	// WorkerQueueSize Worker队列大小
	WorkerQueueSize int `json:"worker_queue_size" yaml:"worker_queue_size"`
	// PoolShardCount 连接池分片数
	PoolShardCount int `json:"pool_shard_count" yaml:"pool_shard_count"`
}

// DefaultServerConfig 默认服务器配置
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		ListenAddr:       ":8080",
		MaxConnections:   1000000,
		ReadTimeout:      30 * time.Second,
		WriteTimeout:     30 * time.Second,
		IdleTimeout:      300 * time.Second,
		ReadBufferSize:   4096,
		WriteBufferSize:  4096,
		WorkerMinWorkers: 100,
		WorkerMaxWorkers: 10000,
		WorkerQueueSize:  100000,
		PoolShardCount:   16,
	}
}
