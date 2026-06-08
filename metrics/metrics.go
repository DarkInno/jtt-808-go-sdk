package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

// Metrics 监控指标
type Metrics struct {
	// 连接指标
	activeConnections int64
	totalConnections  int64
	connectionErrors  int64

	// 消息指标
	receivedMessages int64
	sentMessages     int64
	messageErrors    int64

	// 性能指标
	messageLatencySum   int64
	messageLatencyCount int64

	// 系统指标
	goroutineCount int64
	memoryUsage    int64

	// 时间指标
	startTime time.Time
	mu        sync.RWMutex
}

// NewMetrics 创建监控指标
func NewMetrics() *Metrics {
	return &Metrics{
		startTime: time.Now(),
	}
}

// IncrActiveConnections 增加活跃连接数
func (m *Metrics) IncrActiveConnections() {
	atomic.AddInt64(&m.activeConnections, 1)
	atomic.AddInt64(&m.totalConnections, 1)
}

// DecrActiveConnections 减少活跃连接数
func (m *Metrics) DecrActiveConnections() {
	atomic.AddInt64(&m.activeConnections, -1)
}

// GetActiveConnections 获取活跃连接数
func (m *Metrics) GetActiveConnections() int64 {
	return atomic.LoadInt64(&m.activeConnections)
}

// GetTotalConnections 获取总连接数
func (m *Metrics) GetTotalConnections() int64 {
	return atomic.LoadInt64(&m.totalConnections)
}

// IncrConnectionErrors 增加连接错误数
func (m *Metrics) IncrConnectionErrors() {
	atomic.AddInt64(&m.connectionErrors, 1)
}

// GetConnectionErrors 获取连接错误数
func (m *Metrics) GetConnectionErrors() int64 {
	return atomic.LoadInt64(&m.connectionErrors)
}

// IncrReceivedMessages 增加接收消息数
func (m *Metrics) IncrReceivedMessages() {
	atomic.AddInt64(&m.receivedMessages, 1)
}

// GetReceivedMessages 获取接收消息数
func (m *Metrics) GetReceivedMessages() int64 {
	return atomic.LoadInt64(&m.receivedMessages)
}

// IncrSentMessages 增加发送消息数
func (m *Metrics) IncrSentMessages() {
	atomic.AddInt64(&m.sentMessages, 1)
}

// GetSentMessages 获取发送消息数
func (m *Metrics) GetSentMessages() int64 {
	return atomic.LoadInt64(&m.sentMessages)
}

// IncrMessageErrors 增加消息错误数
func (m *Metrics) IncrMessageErrors() {
	atomic.AddInt64(&m.messageErrors, 1)
}

// GetMessageErrors 获取消息错误数
func (m *Metrics) GetMessageErrors() int64 {
	return atomic.LoadInt64(&m.messageErrors)
}

// RecordMessageLatency 记录消息延迟
func (m *Metrics) RecordMessageLatency(latency time.Duration) {
	atomic.AddInt64(&m.messageLatencySum, int64(latency))
	atomic.AddInt64(&m.messageLatencyCount, 1)
}

// GetMessageLatency 获取平均消息延迟
func (m *Metrics) GetMessageLatency() time.Duration {
	sum := atomic.LoadInt64(&m.messageLatencySum)
	count := atomic.LoadInt64(&m.messageLatencyCount)
	if count == 0 {
		return 0
	}
	return time.Duration(sum / count)
}

// SetGoroutineCount 设置goroutine数量
func (m *Metrics) SetGoroutineCount(count int) {
	atomic.StoreInt64(&m.goroutineCount, int64(count))
}

// GetGoroutineCount 获取goroutine数量
func (m *Metrics) GetGoroutineCount() int64 {
	return atomic.LoadInt64(&m.goroutineCount)
}

// SetMemoryUsage 设置内存使用量
func (m *Metrics) SetMemoryUsage(usage int64) {
	atomic.StoreInt64(&m.memoryUsage, usage)
}

// GetMemoryUsage 获取内存使用量
func (m *Metrics) GetMemoryUsage() int64 {
	return atomic.LoadInt64(&m.memoryUsage)
}

// GetStartTime 获取启动时间
func (m *Metrics) GetStartTime() time.Time {
	return m.startTime
}

// GetUptime 获取运行时间
func (m *Metrics) GetUptime() time.Duration {
	return time.Since(m.startTime)
}

// GetStats 获取统计信息
func (m *Metrics) GetStats() Stats {
	return Stats{
		ActiveConnections: m.GetActiveConnections(),
		TotalConnections:  m.GetTotalConnections(),
		ConnectionErrors:  m.GetConnectionErrors(),
		ReceivedMessages:  m.GetReceivedMessages(),
		SentMessages:      m.GetSentMessages(),
		MessageErrors:     m.GetMessageErrors(),
		MessageLatency:    m.GetMessageLatency(),
		GoroutineCount:    m.GetGoroutineCount(),
		MemoryUsage:       m.GetMemoryUsage(),
		StartTime:         m.startTime,
		Uptime:            m.GetUptime(),
	}
}

// Stats 统计信息
type Stats struct {
	ActiveConnections int64
	TotalConnections  int64
	ConnectionErrors  int64
	ReceivedMessages  int64
	SentMessages      int64
	MessageErrors     int64
	MessageLatency    time.Duration
	GoroutineCount    int64
	MemoryUsage       int64
	StartTime         time.Time
	Uptime            time.Duration
}

// Reset 重置指标
func (m *Metrics) Reset() {
	atomic.StoreInt64(&m.activeConnections, 0)
	atomic.StoreInt64(&m.totalConnections, 0)
	atomic.StoreInt64(&m.connectionErrors, 0)
	atomic.StoreInt64(&m.receivedMessages, 0)
	atomic.StoreInt64(&m.sentMessages, 0)
	atomic.StoreInt64(&m.messageErrors, 0)
	atomic.StoreInt64(&m.messageLatencySum, 0)
	atomic.StoreInt64(&m.messageLatencyCount, 0)
	m.startTime = time.Now()
}
