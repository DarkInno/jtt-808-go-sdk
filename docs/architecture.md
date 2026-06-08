# JT/T 808 Go SDK 架构设计文档

本文档详细介绍 JT/T 808 Go SDK 的架构设计、模块划分和设计理念。

## 目录

- [架构概述](#架构概述)
- [分层架构](#分层架构)
- [模块设计](#模块设计)
- [核心接口](#核心接口)
- [数据流设计](#数据流设计)
- [并发模型](#并发模型)
- [扩展机制](#扩展机制)
- [目录结构](#目录结构)

## 架构概述

JT/T 808 Go SDK 采用分层架构设计，将系统划分为多个独立的模块，各模块之间通过接口进行交互，实现高内聚、低耦合的设计目标。

### 设计原则

1. **接口驱动**：所有核心功能通过接口定义，支持多种实现
2. **可扩展性**：通过中间件、钩子、插件机制支持功能扩展
3. **高性能**：采用异步处理、对象复用等技术优化性能
4. **易用性**：提供简洁的 API 和丰富的示例

## 分层架构

```
┌─────────────────────────────────────────────────────┐
│                   应用层 (Application)                │
│   用户业务逻辑、自定义消息处理器、扩展功能              │
├─────────────────────────────────────────────────────┤
│                   集成层 (Integration)                │
│   消息队列、外部系统集成、推送服务                      │
├─────────────────────────────────────────────────────┤
│                   业务层 (Business)                   │
│   消息路由、业务规则、数据转换                          │
├─────────────────────────────────────────────────────┤
│                   协议层 (Protocol)                   │
│   JT/T 808协议编解码、消息验证、会话管理                │
├─────────────────────────────────────────────────────┤
│                   传输层 (Transport)                  │
│   TCP/UDP连接管理、连接池、网络IO                     │
├─────────────────────────────────────────────────────┤
│                   基础设施层 (Infrastructure)          │
│   配置管理、日志、监控、存储接口                        │
└─────────────────────────────────────────────────────┘
```

### 各层职责

| 层次 | 职责 | 模块 |
|------|------|------|
| 应用层 | 用户业务逻辑 | cmd/, examples/ |
| 集成层 | 外部系统集成 | publisher/ |
| 业务层 | 消息路由处理 | middleware/, hooks/ |
| 协议层 | 协议编解码 | protocol/ |
| 传输层 | 网络通信 | transport/ |
| 基础设施层 | 基础服务 | core/, storage/, logger/, metrics/, config/ |

## 模块设计

### 核心模块 (core/)

核心模块定义了系统的基础接口和数据结构。

```
core/
├── interface.go    # 核心接口定义
├── pool.go         # 连接池实现
└── worker.go       # Worker Pool实现
```

#### 核心接口

```go
// Server 服务器接口
type Server interface {
    Start() error
    Stop() error
    RegisterHandler(msgID uint16, handler MessageHandler)
    GetConnection(deviceID string) (Connection, error)
    GetStats() ServerStats
    Use(middleware Middleware)
    OnConnect(hook func(conn Connection) error)
    OnDisconnect(hook func(conn Connection) error)
    OnError(hook func(conn Connection, err error))
}

// Connection 连接接口
type Connection interface {
    Send(msg *Message) error
    Close() error
    DeviceID() string
    IsConnected() bool
    RemoteAddr() net.Addr
    Set(key string, value interface{})
    Get(key string) (interface{}, bool)
    LastActiveTime() time.Time
    Context() context.Context
}

// Storage 存储接口
type Storage interface {
    SaveLocation(ctx context.Context, loc *LocationReport) error
    GetLocations(ctx context.Context, deviceID string, start, end time.Time) ([]*LocationReport, error)
    SaveAlarm(ctx context.Context, alarm *AlarmReport) error
    GetAlarms(ctx context.Context, deviceID string, start, end time.Time) ([]*AlarmReport, error)
    SaveDevice(ctx context.Context, deviceID string, info *TerminalRegister) error
    GetDevice(ctx context.Context, deviceID string) (*TerminalRegister, error)
    UpdateDeviceStatus(ctx context.Context, deviceID string, status DeviceStatus) error
    Close() error
    Ping() error
}
```

### 协议模块 (protocol/)

协议模块负责 JT/T 808 协议的编解码。

```
protocol/
├── codec.go        # 编解码器实现
├── codec_test.go   # 单元测试
└── types/
    └── constants.go # 常量定义
```

#### 编解码器

```go
// Codec 编解码器
type Codec struct {
    pool sync.Pool
}

// 编码消息
func (c *Codec) Encode(msg *Message) ([]byte, error)

// 解码消息
func (c *Codec) Decode(data []byte) (*Message, error)
```

#### 消息解析

```go
// 解析位置信息上报
func ParseLocationReport(body []byte) (*core.LocationReport, error)

// 解析终端注册
func ParseTerminalRegister(body []byte) (*core.TerminalRegister, error)
```

### 传输模块 (transport/)

传输模块负责网络连接管理。

```
transport/
└── tcp.go          # TCP服务器实现
```

#### TCP服务器

```go
// TCPServer TCP服务器
type TCPServer struct {
    listener    net.Listener
    codec       *protocol.Codec
    connections sync.Map
    handlers    map[uint16]core.MessageHandler
    middleware  []core.Middleware
    hooks       *Hooks
    stats       *Stats
    config      *Config
}
```

### 存储模块 (storage/)

存储模块提供数据持久化支持。

```
storage/
├── memory.go       # 内存存储实现
├── mysql/
│   └── mysql.go    # MySQL实现
├── postgres/
│   └── postgres.go # PostgreSQL实现
└── redis/
    └── redis.go    # Redis实现
```

### 消息队列模块 (publisher/)

消息队列模块提供消息推送支持。

```
publisher/
├── kafka/
│   └── kafka.go    # Kafka实现
├── rabbitmq/
│   └── rabbitmq.go # RabbitMQ实现
└── redis/
    └── redis.go    # Redis Pub/Sub实现
```

### 中间件模块 (middleware/)

中间件模块提供通用的消息处理逻辑。

```
middleware/
├── auth.go         # 认证中间件
├── logging.go      # 日志中间件
├── ratelimit.go    # 限流中间件
├── recovery.go     # 恢复中间件
└── timeout.go      # 超时中间件
```

### 日志模块 (logger/)

日志模块提供结构化日志支持。

```
logger/
├── logger.go       # 日志接口定义
└── zap/
    └── zap.go      # Zap日志实现
```

### 监控模块 (metrics/)

监控模块提供性能指标统计。

```
metrics/
└── metrics.go      # 监控指标定义
```

## 核心接口

### 消息处理器

```go
// MessageHandler 消息处理器
type MessageHandler func(ctx context.Context, conn Connection, msg *Message) error
```

### 中间件

```go
// Middleware 中间件类型
type Middleware func(MessageHandler) MessageHandler
```

### 插件

```go
// Plugin 插件接口
type Plugin interface {
    Name() string
    Version() string
    Initialize(server Server) error
    Shutdown() error
}
```

## 数据流设计

### 消息接收流程

```
客户端连接 → 传输层接收 → 协议层解码 → 中间件处理 → 业务处理器
                                                      ↓
                                              存储层持久化
                                                      ↓
                                              集成层推送
```

### 消息发送流程

```
业务逻辑 → 构建消息 → 协议层编码 → 传输层发送 → 客户端接收
```

### 详细流程

1. **连接建立**
   - 客户端发起 TCP 连接
   - 传输层接受连接，创建 Connection 对象
   - 执行 OnConnect 钩子

2. **消息接收**
   - 传输层读取数据流
   - 协议层解码消息
   - 执行中间件链
   - 路由到消息处理器

3. **消息处理**
   - 消息处理器处理业务逻辑
   - 可选：保存到存储
   - 可选：推送到消息队列
   - 可选：发送响应消息

4. **连接断开**
   - 执行 OnDisconnect 钩子
   - 清理连接资源
   - 更新统计信息

## 并发模型

### Goroutine per Connection

每个 TCP 连接使用独立的 goroutine 处理读取：

```go
func (s *TCPServer) handleConn(netConn net.Conn) {
    defer s.wg.Done()
    
    // 创建连接对象
    conn := NewTCPConnection(netConn, s)
    
    // 读取消息循环
    for {
        msg, err := s.readMessage(reader)
        if err != nil {
            return
        }
        
        // 处理消息
        s.handleMessage(conn, msg)
    }
}
```

### Worker Pool

使用 Worker Pool 处理业务逻辑，避免 goroutine 过多：

```go
type WorkerPool struct {
    taskQueue chan Task
    workers   []*Worker
}

func (p *WorkerPool) Submit(task Task) {
    p.taskQueue <- task
}
```

### 连接池

使用分片连接池减少锁竞争：

```go
type ShardedConnPool struct {
    shards     []*connShard
    shardCount int
}

type connShard struct {
    mu    sync.RWMutex
    conns map[string]*Connection
}
```

## 扩展机制

### 中间件

中间件用于在消息处理前后执行通用逻辑：

```go
// 定义中间件
func LoggingMiddleware(next core.MessageHandler) core.MessageHandler {
    return func(ctx context.Context, conn core.Connection, msg *core.Message) error {
        log.Printf("处理消息: %d", msg.Header.MsgID)
        return next(ctx, conn, msg)
    }
}

// 使用中间件
server.Use(LoggingMiddleware)
```

### 钩子函数

钩子函数用于在特定事件发生时执行自定义逻辑：

```go
// 连接建立钩子
server.OnConnect(func(conn core.Connection) error {
    log.Printf("设备连接: %s", conn.DeviceID())
    return nil
})

// 连接断开钩子
server.OnDisconnect(func(conn core.Connection) error {
    log.Printf("设备断开: %s", conn.DeviceID())
    return nil
})

// 错误处理钩子
server.OnError(func(conn core.Connection, err error) {
    log.Printf("错误: %s, %v", conn.DeviceID(), err)
})
```

### 插件系统

插件用于扩展服务器功能：

```go
type MyPlugin struct{}

func (p *MyPlugin) Name() string { return "my-plugin" }
func (p *MyPlugin) Version() string { return "1.0.0" }
func (p *MyPlugin) Initialize(server core.Server) error {
    // 初始化插件
    return nil
}
func (p *MyPlugin) Shutdown() error {
    // 清理资源
    return nil
}

// 注册插件
server.RegisterPlugin(&MyPlugin{})
```

## 目录结构

```
jt808-go-sdk/
├── cmd/                    # 命令行工具
│   └── server/
│       └── main.go         # 服务端入口
├── config/                 # 配置管理
│   └── config.go
├── core/                   # 核心模块
│   ├── interface.go        # 核心接口定义
│   ├── pool.go             # 连接池实现
│   └── worker.go           # Worker Pool实现
├── docs/                   # 文档
│   ├── api/                # API文档
│   ├── architecture.md     # 架构设计文档
│   ├── deployment.md       # 部署指南
│   ├── getting-started.md  # 快速开始指南
│   └── technical-design.md # 技术设计文档
├── examples/               # 使用示例
│   ├── basic/              # 基础示例
│   ├── advanced/           # 高级示例
│   └── custom/             # 自定义示例
├── hooks/                  # 钩子函数
│   └── hooks.go
├── logger/                 # 日志模块
│   ├── logger.go           # 日志接口定义
│   └── zap/
│       └── zap.go          # Zap日志实现
├── metrics/                # 监控指标
│   └── metrics.go
├── middleware/             # 中间件
│   ├── auth.go             # 认证中间件
│   ├── logging.go          # 日志中间件
│   ├── ratelimit.go        # 限流中间件
│   ├── recovery.go         # 恢复中间件
│   └── timeout.go          # 超时中间件
├── protocol/               # 协议编解码
│   ├── codec.go            # 编解码器实现
│   ├── codec_test.go       # 单元测试
│   └── types/
│       └── constants.go    # 常量定义
├── publisher/              # 消息队列
│   ├── kafka/
│   │   └── kafka.go        # Kafka实现
│   ├── rabbitmq/
│   │   └── rabbitmq.go     # RabbitMQ实现
│   └── redis/
│       └── redis.go        # Redis Pub/Sub实现
├── storage/                # 存储模块
│   ├── memory.go           # 内存存储实现
│   ├── mysql/
│   │   └── mysql.go        # MySQL实现
│   ├── postgres/
│   │   └── postgres.go     # PostgreSQL实现
│   └── redis/
│       └── redis.go        # Redis实现
├── test/                   # 测试工具
│   ├── client/             # 测试客户端
│   ├── stress/             # 压力测试
│   └── integration/        # 集成测试
├── transport/              # 传输层
│   └── tcp.go              # TCP服务器实现
├── go.mod                  # Go模块定义
├── go.sum                  # 依赖校验
├── LICENSE                 # 开源协议
└── README.md               # 项目说明
```

## 设计决策

### 为什么使用接口驱动设计？

1. **可测试性**：可以轻松创建 mock 实现进行单元测试
2. **可扩展性**：用户可以根据需求实现自定义版本
3. **解耦**：模块之间通过接口交互，降低耦合度

### 为什么使用分片连接池？

1. **减少锁竞争**：每个分片独立加锁
2. **支持高并发**：百万连接场景下性能更优
3. **动态扩容**：可以根据负载动态调整分片数量

### 为什么使用 Worker Pool？

1. **资源控制**：限制并发 goroutine 数量
2. **任务排队**：高负载时任务排队处理
3. **动态调整**：可以根据负载动态调整 worker 数量

## 性能优化

### 内存优化

```go
// 使用 sync.Pool 复用对象
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 4096)
    },
}

// 使用对象复用
buf := bufferPool.Get().([]byte)
defer bufferPool.Put(buf)
```

### 并发优化

```go
// 使用分片减少锁竞争
type ShardedMap struct {
    shards [256]struct {
        mu    sync.RWMutex
        items map[string]interface{}
    }
}

func (m *ShardedMap) getShard(key string) *shard {
    hash := fnv.New32a()
    hash.Write([]byte(key))
    return &m.shards[hash.Sum32()%256]
}
```

## 下一步

- [快速开始指南](getting-started.md) - 快速上手使用
- [API 文档](api/README.md) - 完整的 API 参考
- [部署指南](deployment.md) - 生产环境部署
