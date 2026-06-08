# JT/T 808 Go SDK 自定义插件示例

这个示例展示了如何创建和使用自定义插件来扩展 JT/T 808 服务器的功能。

## 功能说明

- 实现自定义插件接口
- 位置信息处理插件
- 报警检测插件
- 中间件与插件结合使用

## 运行方式

```bash
# 进入示例目录
cd examples/custom

# 运行示例
go run main.go
```

服务器将在 `:8080` 端口监听连接。

## 插件架构

### 插件接口

```go
type Plugin interface {
    Name() string
    Version() string
    Initialize(server Server) error
    Shutdown() error
}
```

### 插件类型

1. **位置信息插件 (LocationPlugin)**
   - 处理位置信息上报
   - 支持扩展业务逻辑

2. **报警检测插件 (AlarmPlugin)**
   - 实时检测超速、疲劳驾驶等报警
   - 支持自定义报警规则

## 代码结构

### 1. 插件实现

```go
type LocationPlugin struct {
    server core.Server
    name   string
}

func (p *LocationPlugin) Name() string {
    return p.name
}

func (p *LocationPlugin) Version() string {
    return "1.0.0"
}

func (p *LocationPlugin) Initialize(server core.Server) error {
    p.server = server
    server.RegisterHandler(types.MsgIDLocationReport, p.handleLocationReport)
    return nil
}

func (p *LocationPlugin) Shutdown() error {
    return nil
}
```

### 2. 中间件插件

```go
func (p *AlarmPlugin) alarmDetectionMiddleware() core.Middleware {
    return func(next core.MessageHandler) core.MessageHandler {
        return func(ctx context.Context, conn core.Connection, msg *core.Message) error {
            // 前置处理：报警检测
            if msg.Header.MsgID == types.MsgIDLocationReport {
                report, err := protocol.ParseLocationReport(msg.Body)
                if err == nil {
                    p.checkAlarms(conn, report)
                }
            }
            
            // 调用下一个处理器
            return next(ctx, conn, msg)
        }
    }
}
```

### 3. 插件注册

```go
// 创建插件
locationPlugin := NewLocationPlugin()
alarmPlugin := NewAlarmPlugin()

// 初始化插件
if err := locationPlugin.Initialize(server); err != nil {
    log.Fatalf("初始化位置插件失败: %v", err)
}

if err := alarmPlugin.Initialize(server); err != nil {
    log.Fatalf("初始化报警插件失败: %v", err)
}
```

## 报警检测规则

### 超速报警
- 阈值：120 km/h
- 检测位置：位置信息上报时

### 疲劳驾驶
- 阈值：连续驾驶超过4小时
- 需要更复杂的实现逻辑

### 区域报警
- 支持地理围栏检测
- 需要额外的地理数据支持

## 扩展建议

1. **数据持久化插件**
   - 支持多种数据库（MySQL、PostgreSQL、MongoDB）
   - 支持数据分片和归档

2. **消息推送插件**
   - WebSocket实时推送
   - MQTT消息转发
   - HTTP Webhook

3. **数据分析插件**
   - 实时轨迹分析
   - 驾驶行为分析
   - 统计报表生成

4. **安全认证插件**
   - 终端身份验证
   - 数据加密传输
   - 访问控制

## 注意事项

- 插件初始化顺序很重要，中间件插件需要在消息处理器插件之前初始化
- 插件之间应该解耦，避免循环依赖
- 生产环境建议使用依赖注入框架管理插件
- 插件应该支持热插拔和动态配置