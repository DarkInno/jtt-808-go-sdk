package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// Level 日志级别
type Level int

const (
	// DebugLevel 调试级别
	DebugLevel Level = iota
	// InfoLevel 信息级别
	InfoLevel
	// WarnLevel 警告级别
	WarnLevel
	// ErrorLevel 错误级别
	ErrorLevel
	// FatalLevel 致命级别
	FatalLevel
)

// Field 日志字段
type Field struct {
	Key string
	Val interface{}
}

// String 创建字符串字段
func String(key, val string) Field {
	return Field{Key: key, Val: val}
}

// Int 创建整数字段
func Int(key string, val int) Field {
	return Field{Key: key, Val: val}
}

// Int64 创建64位整数字段
func Int64(key string, val int64) Field {
	return Field{Key: key, Val: val}
}

// Float64 创建浮点数字段
func Float64(key string, val float64) Field {
	return Field{Key: key, Val: val}
}

// Bool 创建布尔字段
func Bool(key string, val bool) Field {
	return Field{Key: key, Val: val}
}

// Error 创建错误字段
func Error(key string, val error) Field {
	return Field{Key: key, Val: val}
}

// Duration 创建时间间隔字段
func Duration(key string, val time.Duration) Field {
	return Field{Key: key, Val: val}
}

// Logger 日志记录器
type Logger struct {
	level  Level
	output *os.File
	mu     sync.Mutex
	fields []Field
}

// NewLogger 创建日志记录器
func NewLogger(level Level) *Logger {
	return &Logger{
		level:  level,
		output: os.Stdout,
	}
}

// With 添加字段
func (l *Logger) With(fields ...Field) *Logger {
	newLogger := &Logger{
		level:  l.level,
		output: l.output,
		fields: make([]Field, len(l.fields)+len(fields)),
	}
	copy(newLogger.fields, l.fields)
	copy(newLogger.fields[len(l.fields):], fields)
	return newLogger
}

// Debug 调试日志
func (l *Logger) Debug(msg string, fields ...Field) {
	if l.level <= DebugLevel {
		l.log("DEBUG", msg, fields...)
	}
}

// Info 信息日志
func (l *Logger) Info(msg string, fields ...Field) {
	if l.level <= InfoLevel {
		l.log("INFO", msg, fields...)
	}
}

// Warn 警告日志
func (l *Logger) Warn(msg string, fields ...Field) {
	if l.level <= WarnLevel {
		l.log("WARN", msg, fields...)
	}
}

// Error 错误日志
func (l *Logger) Error(msg string, fields ...Field) {
	if l.level <= ErrorLevel {
		l.log("ERROR", msg, fields...)
	}
}

// Fatal 致命错误日志
func (l *Logger) Fatal(msg string, fields ...Field) {
	if l.level <= FatalLevel {
		l.log("FATAL", msg, fields...)
	}
	os.Exit(1)
}

// Sync 同步日志
func (l *Logger) Sync() error {
	return l.output.Sync()
}

// log 写入日志
func (l *Logger) log(level, msg string, fields ...Field) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 构建日志消息
	logMsg := fmt.Sprintf("%s [%s] %s", time.Now().Format("2006-01-02 15:04:05.000"), level, msg)

	// 添加字段
	allFields := make([]Field, 0, len(l.fields)+len(fields))
	allFields = append(allFields, l.fields...)
	allFields = append(allFields, fields...)

	if len(allFields) > 0 {
		logMsg += " "
		for i, f := range allFields {
			if i > 0 {
				logMsg += " "
			}
			logMsg += fmt.Sprintf("%s=%v", f.Key, f.Val)
		}
	}

	logMsg += "\n"
	l.output.WriteString(logMsg)
}

// SetLevel 设置日志级别
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// SetOutput 设置输出
func (l *Logger) SetOutput(output *os.File) {
	l.output = output
}
