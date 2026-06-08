# JT/T 808 Go SDK 高级功能示例

这个示例展示了 JT/T 808 Go SDK 的高级功能，包括中间件使用、自定义存储、监控指标等。

## 功能说明

- 使用日志中间件记录请求日志
- 使用自定义中间件实现消息统计
- 自定义存储实现（基于内存存储扩展）
- 超速检测和报警
- 详细的监控指标统计
- 连接管理和错误处理

## 运行方式

```bash
# 进入示例目录
cd examples/advanced

# 运行示例
go run main.go
```

服务器将在 `:8080` 端口监听连接。

## 代码结构

### 1. 中间件使用

```go
// 添加日志中间件
server.Use(middleware.Logging(log))

// 添加自定义中间件
server.Use(func(next core.MessageHandler) core.MessageHandler {
    return func(ctx context.Context, conn core.Connection, msg *core.Message) error {
        // 前置处理
        start := time.Now()
        
        // 调用下一个处理器
        err := next(ctx, conn, msg)
        
        // 后置处理
        duration := time.Since(start)
        // ...
        
        return err
    }
})
```

### 2. 自定义存储

```go
type CustomStorage struct {
    *storage.MemoryStorage
}

func (s *CustomStorage) SaveLocation(ctx context.Context, loc *core.LocationReport) error {
    // 添加自定义逻辑
    // ...
    
    // 调用原始存储方法
    return s.MemoryStorage.SaveLocation(ctx, loc)
}
```

### 3. 监控指标

- 活跃连接数
- 总连接数
- 接收消息数
- 错误消息数
- 设备数量
- 位置记录数量

## 支持的消息类型

| 消息ID | 说明 | 处理功能 |
|--------|------|----------|
| 0x0100 | 终端注册 | 设备信息保存 |
| 0x0200 | 位置信息汇报 | 位置保存、超速检测 |

## 配置参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| ListenAddr | :8080 | 监听地址 |
| MaxConnections | 50000 | 最大连接数 |
| ReadTimeout | 60s | 读超时 |
| WriteTimeout | 60s | 写超时 |

## 监控指标

服务器每30秒打印一次统计信息：

```
服务器统计 active_connections=10 total_connections=150 received_messages=5000 error_count=5 device_count=8 location_count=4500
```

## 注意事项

- 本示例使用自定义内存存储，重启后数据会丢失
- 生产环境建议使用 Redis 或数据库存储
- 超速阈值可在代码中调整（当前为120km/h）
- 中间件按注册顺序执行