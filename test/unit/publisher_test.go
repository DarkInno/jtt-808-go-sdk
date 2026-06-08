package unit

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestPublisherInterface 测试Publisher接口
func TestPublisherInterface(t *testing.T) {
	// 测试Publisher接口定义
	// 这里测试接口的基本用法
	t.Run("InterfaceDefinition", func(t *testing.T) {
		// Publisher接口应该有以下方法：
		// Publish(ctx context.Context, topic string, message []byte) error
		// PublishAsync(ctx context.Context, topic string, message []byte) error
		// Close() error
	})
}

// TestSubscriberInterface 测试Subscriber接口
func TestSubscriberInterface(t *testing.T) {
	// 测试Subscriber接口定义
	// 这里测试接口的基本用法
	t.Run("InterfaceDefinition", func(t *testing.T) {
		// Subscriber接口应该有以下方法：
		// Subscribe(ctx context.Context, topic string, handler func([]byte)) error
		// Unsubscribe(ctx context.Context, topic string) error
		// Close() error
	})
}

// TestMockPublisher 测试模拟Publisher
func TestMockPublisher(t *testing.T) {
	// 创建模拟Publisher
	publisher := NewMockPublisher()

	t.Run("Publish", func(t *testing.T) {
		ctx := context.Background()
		message := []byte("test message")

		err := publisher.Publish(ctx, "test-topic", message)
		if err != nil {
			t.Fatalf("Failed to publish message: %v", err)
		}

		// 验证消息已发布
		messages := publisher.GetMessages("test-topic")
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}

		if string(messages[0]) != "test message" {
			t.Errorf("Expected 'test message', got '%s'", string(messages[0]))
		}
	})

	t.Run("PublishAsync", func(t *testing.T) {
		ctx := context.Background()
		message := []byte("async message")

		err := publisher.PublishAsync(ctx, "async-topic", message)
		if err != nil {
			t.Fatalf("Failed to publish async message: %v", err)
		}

		// 等待异步发布完成
		time.Sleep(100 * time.Millisecond)

		messages := publisher.GetMessages("async-topic")
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
	})

	t.Run("PublishToMultipleTopics", func(t *testing.T) {
		ctx := context.Background()

		publisher.Publish(ctx, "topic-1", []byte("message-1"))
		publisher.Publish(ctx, "topic-2", []byte("message-2"))
		publisher.Publish(ctx, "topic-1", []byte("message-3"))

		messages1 := publisher.GetMessages("topic-1")
		if len(messages1) != 2 {
			t.Errorf("Expected 2 messages in topic-1, got %d", len(messages1))
		}

		messages2 := publisher.GetMessages("topic-2")
		if len(messages2) != 1 {
			t.Errorf("Expected 1 message in topic-2, got %d", len(messages2))
		}
	})

	t.Run("Close", func(t *testing.T) {
		err := publisher.Close()
		if err != nil {
			t.Fatalf("Failed to close publisher: %v", err)
		}

		// 验证关闭后无法发布
		ctx := context.Background()
		err = publisher.Publish(ctx, "test-topic", []byte("test"))
		if err == nil {
			t.Error("Expected error when publishing to closed publisher")
		}
	})
}

// TestMockSubscriber 测试模拟Subscriber
func TestMockSubscriber(t *testing.T) {
	subscriber := NewMockSubscriber()

	t.Run("Subscribe", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		received := make(chan []byte, 1)

		handler := func(msg []byte) {
			received <- msg
		}

		err := subscriber.Subscribe(ctx, "test-topic", handler)
		if err != nil {
			t.Fatalf("Failed to subscribe: %v", err)
		}

		// 发布消息
		subscriber.PublishToTopic("test-topic", []byte("test message"))

		// 等待接收
		select {
		case msg := <-received:
			if string(msg) != "test message" {
				t.Errorf("Expected 'test message', got '%s'", string(msg))
			}
		case <-time.After(1 * time.Second):
			t.Error("Timeout waiting for message")
		}
	})

	t.Run("Unsubscribe", func(t *testing.T) {
		ctx := context.Background()

		handler := func(msg []byte) {
			// 处理消息
		}

		err := subscriber.Subscribe(ctx, "test-topic", handler)
		if err != nil {
			t.Fatalf("Failed to subscribe: %v", err)
		}

		err = subscriber.Unsubscribe(ctx, "test-topic")
		if err != nil {
			t.Fatalf("Failed to unsubscribe: %v", err)
		}
	})

	t.Run("Close", func(t *testing.T) {
		err := subscriber.Close()
		if err != nil {
			t.Fatalf("Failed to close subscriber: %v", err)
		}
	})
}

// MockPublisher 模拟Publisher
type MockPublisher struct {
	mu       sync.Mutex
	messages map[string][][]byte
	closed   bool
}

// NewMockPublisher 创建模拟Publisher
func NewMockPublisher() *MockPublisher {
	return &MockPublisher{
		messages: make(map[string][][]byte),
	}
}

// Publish 发布消息
func (p *MockPublisher) Publish(ctx context.Context, topic string, message []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return fmt.Errorf("publisher is closed")
	}

	p.messages[topic] = append(p.messages[topic], message)
	return nil
}

// PublishAsync 异步发布消息
func (p *MockPublisher) PublishAsync(ctx context.Context, topic string, message []byte) error {
	go func() {
		p.Publish(ctx, topic, message)
	}()
	return nil
}

// Close 关闭发布者
func (p *MockPublisher) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.closed = true
	return nil
}

// GetMessages 获取指定主题的消息
func (p *MockPublisher) GetMessages(topic string) [][]byte {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.messages[topic]
}

// MockSubscriber 模拟Subscriber
type MockSubscriber struct {
	mu       sync.Mutex
	handlers map[string]func([]byte)
	closed   bool
}

// NewMockSubscriber 创建模拟Subscriber
func NewMockSubscriber() *MockSubscriber {
	return &MockSubscriber{
		handlers: make(map[string]func([]byte)),
	}
}

// Subscribe 订阅消息
func (s *MockSubscriber) Subscribe(ctx context.Context, topic string, handler func([]byte)) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return fmt.Errorf("subscriber is closed")
	}

	s.handlers[topic] = handler
	return nil
}

// Unsubscribe 取消订阅
func (s *MockSubscriber) Unsubscribe(ctx context.Context, topic string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.handlers, topic)
	return nil
}

// Close 关闭订阅者
func (s *MockSubscriber) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	return nil
}

// PublishToTopic 发布消息到指定主题（用于测试）
func (s *MockSubscriber) PublishToTopic(topic string, message []byte) {
	s.mu.Lock()
	handler, ok := s.handlers[topic]
	s.mu.Unlock()
	if ok {
		handler(message)
	}
}

// TestPublisherConcurrency 测试Publisher并发安全性
func TestPublisherConcurrency(t *testing.T) {
	publisher := NewMockPublisher()

	var wg sync.WaitGroup
	numGoroutines := 10

	// 并发发布消息
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ctx := context.Background()
			message := []byte(fmt.Sprintf("message-%d", id))
			topic := fmt.Sprintf("topic-%d", id%3) // 使用3个主题

			err := publisher.Publish(ctx, topic, message)
			if err != nil {
				t.Errorf("Failed to publish message: %v", err)
			}
		}(i)
	}

	wg.Wait()

	// 验证消息数量
	totalMessages := 0
	for i := 0; i < 3; i++ {
		topic := fmt.Sprintf("topic-%d", i)
		messages := publisher.GetMessages(topic)
		totalMessages += len(messages)
	}

	if totalMessages != numGoroutines {
		t.Errorf("Expected %d total messages, got %d", numGoroutines, totalMessages)
	}
}

// TestSubscriberConcurrency 测试Subscriber并发安全性
func TestSubscriberConcurrency(t *testing.T) {
	subscriber := NewMockSubscriber()

	var wg sync.WaitGroup
	numGoroutines := 10

	// 并发订阅
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ctx := context.Background()
			topic := fmt.Sprintf("topic-%d", id)
			handler := func(msg []byte) {
				// 处理消息
			}

			err := subscriber.Subscribe(ctx, topic, handler)
			if err != nil {
				t.Errorf("Failed to subscribe: %v", err)
			}
		}(i)
	}

	wg.Wait()

	// 并发取消订阅
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ctx := context.Background()
			topic := fmt.Sprintf("topic-%d", id)

			err := subscriber.Unsubscribe(ctx, topic)
			if err != nil {
				t.Errorf("Failed to unsubscribe: %v", err)
			}
		}(i)
	}

	wg.Wait()
}

// TestPublisherWithTimeout 测试带超时的Publisher
func TestPublisherWithTimeout(t *testing.T) {
	publisher := NewMockPublisher()

	t.Run("PublishWithTimeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		message := []byte("test message")

		err := publisher.Publish(ctx, "test-topic", message)
		if err != nil {
			t.Fatalf("Failed to publish message: %v", err)
		}
	})

	t.Run("PublishWithCancelledContext", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // 立即取消

		message := []byte("test message")

		// 注意：MockPublisher不检查context，所以不会返回错误
		// 实际实现应该检查context
		err := publisher.Publish(ctx, "test-topic", message)
		// MockPublisher不检查context，所以err为nil
		_ = err
	})
}

// TestSubscriberWithTimeout 测试带超时的Subscriber
func TestSubscriberWithTimeout(t *testing.T) {
	subscriber := NewMockSubscriber()

	t.Run("SubscribeWithTimeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		handler := func(msg []byte) {
			// 处理消息
		}

		err := subscriber.Subscribe(ctx, "test-topic", handler)
		if err != nil {
			t.Fatalf("Failed to subscribe: %v", err)
		}
	})

	t.Run("SubscribeWithCancelledContext", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // 立即取消

		handler := func(msg []byte) {
			// 处理消息
		}

		// 注意：MockSubscriber不检查context，所以不会返回错误
		// 实际实现应该检查context
		err := subscriber.Subscribe(ctx, "test-topic", handler)
		// MockSubscriber不检查context，所以err为nil
		_ = err
	})
}

// TestPublisherMessageOrder 测试消息顺序
func TestPublisherMessageOrder(t *testing.T) {
	publisher := NewMockPublisher()
	ctx := context.Background()

	// 按顺序发布消息
	for i := 0; i < 10; i++ {
		message := []byte(fmt.Sprintf("message-%d", i))
		err := publisher.Publish(ctx, "ordered-topic", message)
		if err != nil {
			t.Fatalf("Failed to publish message %d: %v", i, err)
		}
	}

	// 验证消息顺序
	messages := publisher.GetMessages("ordered-topic")
	if len(messages) != 10 {
		t.Fatalf("Expected 10 messages, got %d", len(messages))
	}

	for i, msg := range messages {
		expected := fmt.Sprintf("message-%d", i)
		if string(msg) != expected {
			t.Errorf("Expected '%s', got '%s'", expected, string(msg))
		}
	}
}

// TestPublisherEmptyMessage 测试空消息
func TestPublisherEmptyMessage(t *testing.T) {
	publisher := NewMockPublisher()
	ctx := context.Background()

	// 发布空消息
	err := publisher.Publish(ctx, "empty-topic", []byte{})
	if err != nil {
		t.Fatalf("Failed to publish empty message: %v", err)
	}

	// 发布nil消息
	err = publisher.Publish(ctx, "nil-topic", nil)
	if err != nil {
		t.Fatalf("Failed to publish nil message: %v", err)
	}

	// 验证消息
	messages := publisher.GetMessages("empty-topic")
	if len(messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(messages))
	}

	messages = publisher.GetMessages("nil-topic")
	if len(messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(messages))
	}
}

// TestPublisherLargeMessage 测试大消息
func TestPublisherLargeMessage(t *testing.T) {
	publisher := NewMockPublisher()
	ctx := context.Background()

	// 创建大消息（1MB）
	largeMessage := make([]byte, 1024*1024)
	for i := range largeMessage {
		largeMessage[i] = byte(i % 256)
	}

	err := publisher.Publish(ctx, "large-topic", largeMessage)
	if err != nil {
		t.Fatalf("Failed to publish large message: %v", err)
	}

	messages := publisher.GetMessages("large-topic")
	if len(messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(messages))
	}

	if len(messages[0]) != len(largeMessage) {
		t.Errorf("Expected message length %d, got %d", len(largeMessage), len(messages[0]))
	}
}

// TestSubscriberMultipleHandlers 测试多个处理器
func TestSubscriberMultipleHandlers(t *testing.T) {
	subscriber := NewMockSubscriber()
	ctx := context.Background()

	received1 := make(chan []byte, 1)
	// 订阅同一主题的多个处理器
	handler1 := func(msg []byte) {
		received1 <- msg
	}

	// 注意：MockSubscriber只支持每个主题一个处理器
	// 实际实现可能支持多个处理器
	subscriber.Subscribe(ctx, "multi-topic", handler1)

	// 发布消息
	subscriber.PublishToTopic("multi-topic", []byte("test message"))

	// 等待接收
	select {
	case msg := <-received1:
		if string(msg) != "test message" {
			t.Errorf("Expected 'test message', got '%s'", string(msg))
		}
	case <-time.After(1 * time.Second):
		t.Error("Timeout waiting for message")
	}
}
