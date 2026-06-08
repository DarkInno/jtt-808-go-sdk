package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publisher RabbitMQ消息发布者
type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// Subscriber RabbitMQ消息订阅者
type Subscriber struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	handler     func([]byte)
	consumerTag string
}

// Config RabbitMQ配置
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	VHost    string
}

// NewPublisher 创建RabbitMQ发布者
func NewPublisher(config *Config) (*Publisher, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		config.User, config.Password, config.Host, config.Port, config.VHost)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &Publisher{
		conn:    conn,
		channel: channel,
	}, nil
}

// Publish 发布消息
func (p *Publisher) Publish(ctx context.Context, topic string, message []byte) error {
	// 声明队列
	_, err := p.channel.QueueDeclare(
		topic, // 队列名称
		true,  // 持久化
		false, // 自动删除
		false, // 排他
		false, // 无等待
		nil,   // 参数
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// 发布消息
	return p.channel.PublishWithContext(ctx,
		"",    // 交换机
		topic, // 路由键
		false, // 强制
		false, // 立即
		amqp.Publishing{
			ContentType:  "application/octet-stream",
			DeliveryMode: amqp.Persistent,
			Body:         message,
		})
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
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

// NewSubscriber 创建RabbitMQ订阅者
func NewSubscriber(config *Config, handler func([]byte)) (*Subscriber, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		config.User, config.Password, config.Host, config.Port, config.VHost)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &Subscriber{
		conn:    conn,
		channel: channel,
		handler: handler,
	}, nil
}

// Subscribe 订阅消息
func (s *Subscriber) Subscribe(ctx context.Context, topic string, handler func([]byte)) error {
	// 声明队列
	_, err := s.channel.QueueDeclare(
		topic, // 队列名称
		true,  // 持久化
		false, // 自动删除
		false, // 排他
		false, // 无等待
		nil,   // 参数
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// 消费消息
	if s.consumerTag == "" {
		s.consumerTag = fmt.Sprintf("jt808-%p", s)
	}
	msgs, err := s.channel.Consume(
		topic,         // 队列
		s.consumerTag, // 消费者
		false,         // 自动确认
		false,         // 排他
		false,         // 本地
		false,         // 无等待
		nil,           // 参数
	)
	if err != nil {
		return fmt.Errorf("failed to consume: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-msgs:
			if !ok {
				return nil
			}

			if handler != nil {
				handler(msg.Body)
			} else if s.handler != nil {
				s.handler(msg.Body)
			}

			// 确认消息
			if err := msg.Ack(false); err != nil {
				return fmt.Errorf("failed to ack message: %w", err)
			}
		}
	}
}

// Unsubscribe 取消订阅
func (s *Subscriber) Unsubscribe(ctx context.Context, topic string) error {
	if s.consumerTag == "" {
		return nil
	}
	return s.channel.Cancel(s.consumerTag, false)
}

// Close 关闭订阅者
func (s *Subscriber) Close() error {
	if s.channel != nil {
		s.channel.Close()
	}
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
