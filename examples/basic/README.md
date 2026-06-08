# JT/T 808 Go SDK 基础示例

这是一个最简单的 JT/T 808 服务器示例，展示如何快速搭建一个基础的 GPS 定位服务器。

## 功能说明

- 接收终端注册消息
- 接收位置信息上报
- 使用内存存储数据
- 支持连接管理

## 运行方式

```bash
# 进入示例目录
cd examples/basic

# 运行示例
go run main.go
```

服务器将在 `:8080` 端口监听 JT/T 808 协议连接。

## 测试方法

可以使用 JT/T 808 模拟客户端发送测试数据：

```bash
# 使用 netcat 发送测试数据（需要先进行转义处理）
echo -n "7e010000..." | xxd -r -p | nc localhost 8080
```

## 代码说明

1. **存储初始化**: 使用 `storage.NewMemoryStorage()` 创建内存存储
2. **服务器配置**: 使用 `transport.DefaultConfig()` 获取默认配置
3. **消息处理**: 通过 `RegisterHandler` 注册不同消息类型的处理器
4. **连接钩子**: 通过 `OnConnect` 和 `OnDisconnect` 监听连接状态

## 支持的消息类型

| 消息ID | 说明 | 处理函数 |
|--------|------|----------|
| 0x0100 | 终端注册 | `MsgIDTerminalRegister` |
| 0x0200 | 位置信息汇报 | `MsgIDLocationReport` |

## 注意事项

- 本示例使用内存存储，重启后数据会丢失
- 生产环境建议使用 Redis 或数据库存储
- 未包含鉴权功能，仅用于演示基本流程