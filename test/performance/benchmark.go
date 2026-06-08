package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// 压测参数
	concurrency = 1000    // 并发连接数
	duration    = 60      // 压测持续时间（秒）
	messageRate = 10      // 每个连接每秒发送消息数
	serverAddr  = ":8080" // 服务器地址
)

var (
	totalConnections  int64
	activeConnections int64
	totalMessages     int64
	successMessages   int64
	failedMessages    int64
	totalLatency      int64
)

func main() {
	fmt.Printf("Starting benchmark test...\n")
	fmt.Printf("Concurrency: %d\n", concurrency)
	fmt.Printf("Duration: %d seconds\n", duration)
	fmt.Printf("Message rate: %d per second per connection\n", messageRate)
	fmt.Printf("Server: %s\n\n", serverAddr)

	var wg sync.WaitGroup
	stopCh := make(chan struct{})

	// 启动统计打印
	go printStats()

	// 启动压测
	startTime := time.Now()
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			runWorker(id, stopCh)
		}(i)
	}

	// 等待压测结束
	time.Sleep(time.Duration(duration) * time.Second)
	close(stopCh)
	wg.Wait()

	// 打印最终结果
	printFinalResults(time.Since(startTime))
}

func runWorker(id int, stopCh <-chan struct{}) {
	// 连接服务器
	conn, err := net.DialTimeout("tcp", serverAddr, 5*time.Second)
	if err != nil {
		atomic.AddInt64(&failedMessages, 1)
		return
	}
	defer conn.Close()

	atomic.AddInt64(&totalConnections, 1)
	atomic.AddInt64(&activeConnections, 1)
	defer atomic.AddInt64(&activeConnections, -1)

	// 生成设备ID
	deviceID := fmt.Sprintf("DEVICE_%06d", id)

	// 发送注册消息
	if err := sendRegisterMessage(conn, deviceID); err != nil {
		atomic.AddInt64(&failedMessages, 1)
		return
	}

	// 循环发送位置消息
	ticker := time.NewTicker(time.Second / time.Duration(messageRate))
	defer ticker.Stop()

	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			start := time.Now()
			if err := sendLocationMessage(conn, deviceID); err != nil {
				atomic.AddInt64(&failedMessages, 1)
				return
			}
			latency := time.Since(start)
			atomic.AddInt64(&totalLatency, int64(latency))
			atomic.AddInt64(&successMessages, 1)
			atomic.AddInt64(&totalMessages, 1)
		}
	}
}

func sendRegisterMessage(conn net.Conn, deviceID string) error {
	// 构造注册消息
	body := buildRegisterBody(deviceID)
	msg := buildMessage(0x0100, deviceID, body)
	_, err := conn.Write(msg)
	return err
}

func sendLocationMessage(conn net.Conn, deviceID string) error {
	// 构造位置消息
	body := buildLocationBody()
	msg := buildMessage(0x0200, deviceID, body)
	_, err := conn.Write(msg)
	return err
}

func buildMessage(msgID uint16, deviceID string, body []byte) []byte {
	// 消息头
	header := make([]byte, 12)
	binary.BigEndian.PutUint16(header[0:2], msgID)
	binary.BigEndian.PutUint16(header[2:4], uint16(len(body)))
	copy(header[4:10], encodeBCD(deviceID, 6))
	binary.BigEndian.PutUint16(header[10:12], 1) // 流水号

	// 计算校验码
	var checksum byte
	for _, b := range header {
		checksum ^= b
	}
	for _, b := range body {
		checksum ^= b
	}

	// 组装消息
	msg := make([]byte, 0, len(header)+len(body)+3)
	msg = append(msg, 0x7E)
	msg = append(msg, escape(header)...)
	msg = append(msg, escape(body)...)
	msg = append(msg, checksum)
	msg = append(msg, 0x7E)

	return msg
}

func buildRegisterBody(deviceID string) []byte {
	body := make([]byte, 37)
	binary.BigEndian.PutUint16(body[0:2], 0) // 省域ID
	binary.BigEndian.PutUint16(body[2:4], 0) // 市县域ID
	copy(body[4:9], "TEST0")                 // 制造商ID
	copy(body[9:39], "TEST_TERMINAL")        // 终端型号
	copy(body[39:46], deviceID[:7])          // 终端ID
	body[46] = 1                             // 车牌颜色
	return body
}

func buildLocationBody() []byte {
	body := make([]byte, 28)
	binary.BigEndian.PutUint32(body[0:4], 0)           // 报警标志
	binary.BigEndian.PutUint32(body[4:8], 0x02)        // 状态（已定位）
	binary.BigEndian.PutUint32(body[8:12], 39904200)   // 纬度
	binary.BigEndian.PutUint32(body[12:16], 116407400) // 经度
	binary.BigEndian.PutUint16(body[16:18], 50)        // 海拔
	binary.BigEndian.PutUint16(body[18:20], 600)       // 速度
	binary.BigEndian.PutUint16(body[20:22], 90)        // 方向
	// 时间（BCD编码）
	now := time.Now()
	body[22] = byte((now.Year()%100)/10)<<4 | byte((now.Year()%100)%10)
	body[23] = byte(int(now.Month())/10)<<4 | byte(int(now.Month())%10)
	body[24] = byte(now.Day()/10)<<4 | byte(now.Day()%10)
	body[25] = byte(now.Hour()/10)<<4 | byte(now.Hour()%10)
	body[26] = byte(now.Minute()/10)<<4 | byte(now.Minute()%10)
	body[27] = byte(now.Second()/10)<<4 | byte(now.Second()%10)
	return body
}

func encodeBCD(s string, length int) []byte {
	result := make([]byte, length)
	for i := 0; i < length && i*2 < len(s); i++ {
		high := byte(0)
		low := byte(0)
		if i*2 < len(s) {
			high = s[i*2] - '0'
		}
		if i*2+1 < len(s) {
			low = s[i*2+1] - '0'
		}
		result[i] = (high << 4) | low
	}
	return result
}

func escape(data []byte) []byte {
	var result []byte
	for _, b := range data {
		if b == 0x7E {
			result = append(result, 0x7D, 0x02)
		} else if b == 0x7D {
			result = append(result, 0x7D, 0x01)
		} else {
			result = append(result, b)
		}
	}
	return result
}

func printStats() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		total := atomic.LoadInt64(&totalMessages)
		success := atomic.LoadInt64(&successMessages)
		failed := atomic.LoadInt64(&failedMessages)
		active := atomic.LoadInt64(&activeConnections)
		latency := atomic.LoadInt64(&totalLatency)

		var avgLatency time.Duration
		if total > 0 {
			avgLatency = time.Duration(latency / total)
		}

		fmt.Printf("\rActive: %d | Total: %d | Success: %d | Failed: %d | Avg Latency: %v",
			active, total, success, failed, avgLatency)
	}
}

func printFinalResults(duration time.Duration) {
	total := atomic.LoadInt64(&totalMessages)
	success := atomic.LoadInt64(&successMessages)
	failed := atomic.LoadInt64(&failedMessages)
	latency := atomic.LoadInt64(&totalLatency)

	var avgLatency time.Duration
	if total > 0 {
		avgLatency = time.Duration(latency / total)
	}

	fmt.Printf("\n\n=== Benchmark Results ===\n")
	fmt.Printf("Duration: %v\n", duration)
	fmt.Printf("Total Messages: %d\n", total)
	fmt.Printf("Success Messages: %d\n", success)
	fmt.Printf("Failed Messages: %d\n", failed)
	fmt.Printf("Messages per Second: %.2f\n", float64(total)/duration.Seconds())
	fmt.Printf("Average Latency: %v\n", avgLatency)
	fmt.Printf("Success Rate: %.2f%%\n", float64(success)/float64(total)*100)
}
