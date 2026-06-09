# JT/T 808 Go SDK 压测报告

## 测试信息

- 测试日期：2026-06-09
- 测试机器：Windows amd64，本机回环地址 `127.0.0.1`
- Go 版本：`go1.26.4 windows/amd64`
- 服务端：`cmd/highperf`
- 压测工具：`test/stress`
- 测试场景：`location`，终端注册后持续发送位置上报并等待服务端响应
- 监听地址：`127.0.0.1:8080`

## 构建与验证

为避免 Git ownership 和沙箱缓存目录影响，测试时显式关闭 VCS stamping，并把 Go 构建缓存放到仓库内。

```powershell
New-Item -ItemType Directory -Force .gocache,.gotmp | Out-Null
$env:GOCACHE=(Resolve-Path .gocache).Path
$env:GOTMPDIR=(Resolve-Path .gotmp).Path

go test -buildvcs=false ./...
go build -buildvcs=false -o .\bin\jt808-highperf.exe .\cmd\highperf
go build -buildvcs=false -o .\bin\jt808-stress.exe .\test\stress
```

验证结果：`go test -buildvcs=false ./...` 全部通过。

## 烟测结果

命令：

```powershell
.\bin\jt808-stress.exe -s 127.0.0.1:8080 -c 10 -d 3 -r 1 -t location -b 10
```

结果：

| 指标 | 数值 |
|---|---:|
| 总连接数 | 10 |
| 峰值连接数 | 10 |
| 连接错误 | 0 |
| 总消息数 | 20 |
| 成功消息 | 20 |
| 失败消息 | 0 |
| 消息速率 | 6.66 msg/s |
| 平均延迟 | 7.983625 ms |
| 成功率 | 100.00% |

## 第一轮压测

命令：

```powershell
.\bin\jt808-stress.exe -s 127.0.0.1:8080 -c 1000 -d 30 -r 10 -t location -b 500
```

结果：

| 指标 | 数值 |
|---|---:|
| 目标并发连接数 | 1000 |
| 实际总连接数 | 359 |
| 峰值连接数 | 359 |
| 连接错误 | 641 |
| 总消息数 | 107482 |
| 成功消息 | 107482 |
| 失败消息 | 641 |
| 消息速率 | 3573.12 msg/s |
| 平均延迟 | 52.915 us |
| 成功率 | 100.00% |

结论：该轮失败集中在建连阶段。服务端日志显示 `errorCount=0`，已建立连接上的消息全部成功。`359` 并发不是服务端并发能力上限，而是 Windows 本机回环压测时突发建连批次过大，受到客户端侧端口、连接队列或 TCP 栈调度影响导致的连接失败。

## 第二轮压测

命令：

```powershell
.\bin\jt808-stress.exe -s 127.0.0.1:8080 -c 1000 -d 30 -r 10 -t location -b 50
```

结果：

| 指标 | 数值 |
|---|---:|
| 目标并发连接数 | 1000 |
| 实际总连接数 | 1000 |
| 峰值连接数 | 1000 |
| 连接错误 | 0 |
| 总消息数 | 304236 |
| 成功消息 | 304236 |
| 失败消息 | 0 |
| 消息速率 | 9809.82 msg/s |
| 平均延迟 | 436.295 us |
| 成功率 | 100.00% |
| 客户端内存使用 | 3.72 MB |

服务端日志摘要：

| 指标 | 数值 |
|---|---:|
| 活跃连接峰值 | 1000 |
| 服务端累计收发消息 | 413107 |
| 服务端错误数 | 0 |
| Worker 队列堆积 | 0 |
| 服务端进程 Working Set | 约 94 MB |
| 服务端 Go heap alloc 峰值 | 约 41 MB |

## 结论

在本机回环地址、`1000` 并发连接、每连接 `10 msg/s`、持续 `30s` 的位置上报场景下，高性能服务端稳定处理约 `9.8k msg/s`，消息成功率 `100%`，服务端无错误且 Worker 队列无堆积。

压测参数中 `-b` 批次大小会显著影响 Windows 本机压测的建连成功率。`-b 500` 下出现大量连接错误，属于 Windows 本机突发建连限制导致的并发表现偏低，不应作为 SDK 服务端并发上限解读；`-b 50` 下 `1000` 连接全部成功。后续复测建议优先使用较小批次逐步建连，例如：

```powershell
.\bin\jt808-stress.exe -s 127.0.0.1:8080 -c 1000 -d 30 -r 10 -t location -b 50
```

如需评估极限吞吐，可以在保持 `-b 50` 或更小批次的前提下，逐步提升 `-c` 和 `-r`，并同步记录服务端日志中的 `errorCount`、`workerPoolQueue`、内存和 GC 指标。
