package core

import (
	"context"
	"net"
	"time"
)

// Server 服务器接口
type Server interface {
	// Start 启动服务器
	Start() error
	// Stop 停止服务器
	Stop() error
	// RegisterHandler 注册消息处理器
	RegisterHandler(msgID uint16, handler MessageHandler)
	// GetConnection 获取连接
	GetConnection(deviceID string) (Connection, error)
	// GetStats 获取统计信息
	GetStats() ServerStats
	// Use 添加中间件
	Use(middleware Middleware)
	// OnConnect 注册连接建立钩子
	OnConnect(hook func(conn Connection) error)
	// OnDisconnect 注册连接断开钩子
	OnDisconnect(hook func(conn Connection) error)
	// OnError 注册错误处理钩子
	OnError(hook func(conn Connection, err error))
}

// Connection 连接接口
type Connection interface {
	// Send 发送消息
	Send(msg *Message) error
	// Close 关闭连接
	Close() error
	// DeviceID 获取设备ID
	DeviceID() string
	// IsConnected 是否已连接
	IsConnected() bool
	// RemoteAddr 获取远程地址
	RemoteAddr() net.Addr
	// Set 设置连接属性
	Set(key string, value interface{})
	// Get 获取连接属性
	Get(key string) (interface{}, bool)
	// LastActiveTime 获取最后活跃时间
	LastActiveTime() time.Time
	// Context 获取连接上下文
	Context() context.Context
}

// MessageHandler 消息处理器
type MessageHandler func(ctx context.Context, conn Connection, msg *Message) error

// Middleware 中间件类型
type Middleware func(MessageHandler) MessageHandler

// ServerStats 服务器统计信息
type ServerStats struct {
	// ActiveConnections 当前活跃连接数
	ActiveConnections int64
	// TotalConnections 总连接数
	TotalConnections int64
	// ReceivedMessages 接收消息数
	ReceivedMessages int64
	// SentMessages 发送消息数
	SentMessages int64
	// ErrorCount 错误数
	ErrorCount int64
	// StartTime 启动时间
	StartTime time.Time
	// Uptime 运行时间
	Uptime time.Duration
}

// Message 消息结构
type Message struct {
	// Header 消息头
	Header *MessageHeader
	// Body 消息体
	Body []byte
	// Raw 原始数据
	Raw []byte
}

// MessageHeader 消息头结构
type MessageHeader struct {
	// MsgID 消息ID
	MsgID uint16
	// MsgBodyProperty 消息体属性
	MsgBodyProperty uint16
	// ProtocolVersion 协议版本
	ProtocolVersion uint8
	// PhoneNo 终端手机号
	PhoneNo string
	// MsgFlowNo 消息流水号
	MsgFlowNo uint16
	// MsgBodyLength 消息体长度
	MsgBodyLength uint16
	// SubPackage 是否分包
	SubPackage bool
	// Encryption 加密方式
	Encryption uint8
}

// LocationReport 位置信息上报
type LocationReport struct {
	// AlarmFlag 报警标志
	AlarmFlag uint32
	// Status 状态
	Status uint32
	// Latitude 纬度
	Latitude float64
	// Longitude 经度
	Longitude float64
	// Altitude 海拔高度
	Altitude uint16
	// Speed 速度
	Speed uint16
	// Direction 方向
	Direction uint16
	// Time 时间
	Time time.Time
}

// TerminalRegister 终端注册
type TerminalRegister struct {
	// ProvinceID 省域ID
	ProvinceID uint16
	// CityID 市县域ID
	CityID uint16
	// ManufacturerID 制造商ID
	ManufacturerID string
	// TerminalType 终端型号
	TerminalType string
	// TerminalID 终端ID
	TerminalID string
	// PlateColor 车牌颜色
	PlateColor uint8
	// PlateNo 车牌号
	PlateNo string
}

// TerminalRegisterResult 终端注册结果
type TerminalRegisterResult struct {
	// MsgFlowNo 消息流水号
	MsgFlowNo uint16
	// Result 结果
	Result uint8
	// AuthCode 鉴权码
	AuthCode string
}

// TerminalAuth 终端鉴权
type TerminalAuth struct {
	// AuthCode 鉴权码
	AuthCode string
}

// AlarmReport 报警信息
type AlarmReport struct {
	// AlarmFlag 报警标志
	AlarmFlag uint32
	// Status 状态
	Status uint32
	// Latitude 纬度
	Latitude float64
	// Longitude 经度
	Longitude float64
	// Altitude 海拔高度
	Altitude uint16
	// Speed 速度
	Speed uint16
	// Direction 方向
	Direction uint16
	// Time 时间
	Time time.Time
	// AlarmType 报警类型
	AlarmType uint8
}

// DeviceStatus 设备状态
type DeviceStatus int

const (
	// DeviceStatusOffline 离线
	DeviceStatusOffline DeviceStatus = iota
	// DeviceStatusOnline 在线
	DeviceStatusOnline
	// DeviceStatusAuthenticating 认证中
	DeviceStatusAuthenticating
)

// Storage 存储接口
type Storage interface {
	// SaveLocation 保存位置信息
	SaveLocation(ctx context.Context, loc *LocationReport) error
	// GetLocations 获取位置信息
	GetLocations(ctx context.Context, deviceID string, start, end time.Time) ([]*LocationReport, error)
	// SaveAlarm 保存报警信息
	SaveAlarm(ctx context.Context, alarm *AlarmReport) error
	// GetAlarms 获取报警信息
	GetAlarms(ctx context.Context, deviceID string, start, end time.Time) ([]*AlarmReport, error)
	// SaveDevice 保存设备信息
	SaveDevice(ctx context.Context, deviceID string, info *TerminalRegister) error
	// GetDevice 获取设备信息
	GetDevice(ctx context.Context, deviceID string) (*TerminalRegister, error)
	// UpdateDeviceStatus 更新设备状态
	UpdateDeviceStatus(ctx context.Context, deviceID string, status DeviceStatus) error
	// Close 关闭存储连接
	Close() error
	// Ping 检查连接健康
	Ping() error
}

// Publisher 消息发布者接口
type Publisher interface {
	// Publish 发布消息
	Publish(ctx context.Context, topic string, message []byte) error
	// PublishAsync 异步发布消息
	PublishAsync(ctx context.Context, topic string, message []byte) error
	// Close 关闭发布者
	Close() error
}

// Subscriber 消息订阅者接口
type Subscriber interface {
	// Subscribe 订阅消息
	Subscribe(ctx context.Context, topic string, handler func([]byte)) error
	// Unsubscribe 取消订阅
	Unsubscribe(ctx context.Context, topic string) error
	// Close 关闭订阅者
	Close() error
}

// Config 配置接口
type Config interface {
	// GetString 获取字符串配置
	GetString(key string) string
	// GetInt 获取整数配置
	GetInt(key string) int
	// GetBool 获取布尔配置
	GetBool(key string) bool
	// GetDuration 获取时间间隔配置
	GetDuration(key string) time.Duration
	// Set 设置配置
	Set(key string, value interface{})
	// Unmarshal 反序列化配置
	Unmarshal(target interface{}) error
	// Watch 监听配置变化
	Watch(key string, callback func(interface{}))
}

// Logger 日志接口
type Logger interface {
	// Debug 调试日志
	Debug(msg string, fields ...Field)
	// Info 信息日志
	Info(msg string, fields ...Field)
	// Warn 警告日志
	Warn(msg string, fields ...Field)
	// Error 错误日志
	Error(msg string, fields ...Field)
	// Fatal 致命错误日志
	Fatal(msg string, fields ...Field)
	// With 添加字段
	With(fields ...Field) Logger
	// Sync 同步日志
	Sync() error
}

// Field 日志字段
type Field struct {
	Key string
	Val interface{}
}

// Plugin 插件接口
type Plugin interface {
	// Name 插件名称
	Name() string
	// Version 插件版本
	Version() string
	// Initialize 初始化插件
	Initialize(server Server) error
	// Shutdown 关闭插件
	Shutdown() error
}
