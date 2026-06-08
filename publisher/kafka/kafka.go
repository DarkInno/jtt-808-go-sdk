package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

// Publisher Kafka消息发布者
type Publisher struct {
	writer *kafka.Writer
}

// Subscriber Kafka消息订阅者
type Subscriber struct {
	reader  *kafka.Reader
	handler func([]byte)
}

// Config Kafka配置
type Config struct {
	Brokers []string
	Topic   string
	GroupID string
}

// NewPublisher 创建Kafka发布者
func NewPublisher(config *Config) *Publisher {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(config.Brokers...),
		Topic:        config.Topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		BatchSize:    100,
	}

	return &Publisher{writer: writer}
}

// Publish 发布消息
func (p *Publisher) Publish(ctx context.Context, topic string, message []byte) error {
	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", time.Now().UnixNano())),
		Value: message,
	}

	if topic != "" {
		msg.Topic = topic
	}

	return p.writer.WriteMessages(ctx, msg)
}

// PublishAsync 异步发布消息
func (p *Publisher) PublishAsync(ctx context.Context, topic string, message []byte) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	msg := append([]byte(nil), message...)
	go func() {
		_ = p.Publish(ctx, topic, msg)
	}()
	return nil
}

// Close 关闭发布者
func (p *Publisher) Close() error {
	return p.writer.Close()
}

// NewSubscriber 创建Kafka订阅者
func NewSubscriber(config *Config, handler func([]byte)) *Subscriber {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        config.Brokers,
		Topic:          config.Topic,
		GroupID:        config.GroupID,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset,
	})

	return &Subscriber{
		reader:  reader,
		handler: handler,
	}
}

// Subscribe 订阅消息
func (s *Subscriber) Subscribe(ctx context.Context, topic string, handler func([]byte)) error {
	for {
		msg, err := s.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		if handler != nil {
			handler(msg.Value)
		} else if s.handler != nil {
			s.handler(msg.Value)
		}
	}
}

// Unsubscribe 取消订阅
func (s *Subscriber) Unsubscribe(ctx context.Context, topic string) error {
	return s.reader.Close()
}

// Close 关闭订阅者
func (s *Subscriber) Close() error {
	return s.reader.Close()
}
