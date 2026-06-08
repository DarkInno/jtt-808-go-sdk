package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Publisher Redis消息发布者
type Publisher struct {
	client *redis.Client
}

// Subscriber Redis消息订阅者
type Subscriber struct {
	client  *redis.Client
	pubsub  *redis.PubSub
	handler func([]byte)
}

// Config Redis配置
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// NewPublisher 创建Redis发布者
func NewPublisher(config *Config) (*Publisher, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Publisher{client: client}, nil
}

// Publish 发布消息
func (p *Publisher) Publish(ctx context.Context, topic string, message []byte) error {
	return p.client.Publish(ctx, topic, message).Err()
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
	return p.client.Close()
}

// NewSubscriber 创建Redis订阅者
func NewSubscriber(config *Config, handler func([]byte)) (*Subscriber, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Subscriber{
		client:  client,
		handler: handler,
	}, nil
}

// Subscribe 订阅消息
func (s *Subscriber) Subscribe(ctx context.Context, topic string, handler func([]byte)) error {
	s.pubsub = s.client.Subscribe(ctx, topic)
	defer s.pubsub.Close()

	// 等待订阅确认
	_, err := s.pubsub.Receive(ctx)
	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	// 接收消息
	ch := s.pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-ch:
			if !ok {
				return nil
			}

			if handler != nil {
				handler([]byte(msg.Payload))
			} else if s.handler != nil {
				s.handler([]byte(msg.Payload))
			}
		}
	}
}

// Unsubscribe 取消订阅
func (s *Subscriber) Unsubscribe(ctx context.Context, topic string) error {
	if s.pubsub != nil {
		return s.pubsub.Unsubscribe(ctx, topic)
	}
	return nil
}

// Close 关闭订阅者
func (s *Subscriber) Close() error {
	if s.pubsub != nil {
		s.pubsub.Close()
	}
	return s.client.Close()
}
