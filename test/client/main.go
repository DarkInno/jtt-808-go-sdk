package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	deviceCount   = 10 // 模拟设备数量
	serverAddr    = "127.0.0.1:8080"
	messageRate   = 1  // 每个设备每秒发送消息数
	testDuration  = 30 // 测试持续时间（秒）
	wg            sync.WaitGroup
	totalSent     int64
	totalReceived int64
	totalErrors   int64
	connectedDevs int64
)

func main() {
	fmt.Printf("=== JT808 国标设备模拟器 ===\n")
	fmt.Printf("服务器地址: %s\n", serverAddr)
	fmt.Printf("模拟设备数: %d\n", deviceCount)
	fmt.Printf("消息频率: %d/秒/设备\n", messageRate)
	fmt.Printf("测试时长: %d秒\n\n", testDuration)

	// 启动统计打印
	go printStats()

	// 启动设备模拟
	stopCh := make(chan struct{})
	for i := 0; i < deviceCount; i++ {
		wg.Add(1)
		go simulateDevice(i, stopCh)
		time.Sleep(10 * time.Millisecond) // 错开连接时间
	}

	// 等待测试结束
	time.Sleep(time.Duration(testDuration) * time.Second)
	close(stopCh)
	wg.Wait()

	printFinalResults()
}

// simulateDevice 模拟单个设备
func simulateDevice(id int, stopCh <-chan struct{}) {
	defer wg.Done()

	deviceID := fmt.Sprintf("DEVICE_%06d", id)
	phoneNo := fmt.Sprintf("138%08d", id)

	// 连接服务器
	conn, err := net.DialTimeout("tcp", serverAddr, 5*time.Second)
	if err != nil {
		atomic.AddInt64(&totalErrors, 1)
		return
	}
	defer conn.Close()

	atomic.AddInt64(&connectedDevs, 1)
	defer atomic.AddInt64(&connectedDevs, -1)

	// 发送注册消息
	if err := sendRegister(conn, deviceID, phoneNo); err != nil {
		atomic.AddInt64(&totalErrors, 1)
		return
	}

	// 等待注册应答
	if err := readResponse(conn); err != nil {
		atomic.AddInt64(&totalErrors, 1)
		return
	}

	// 发送鉴权消息
	if err := sendAuth(conn, deviceID, phoneNo); err != nil {
		atomic.AddInt64(&totalErrors, 1)
		return
	}

	// 等待鉴权应答
	if err := readResponse(conn); err != nil {
		atomic.AddInt64(&totalErrors, 1)
		return
	}

	// 定期发送心跳和位置信息
	heartbeatTicker := time.NewTicker(30 * time.Second)
	locationTicker := time.NewTicker(time.Second / time.Duration(messageRate))
	defer heartbeatTicker.Stop()
	defer locationTicker.Stop()

	for {
		select {
		case <-stopCh:
			return
		case <-heartbeatTicker.C:
			if err := sendHeartbeat(conn, phoneNo); err != nil {
				return
			}
		case <-locationTicker.C:
			if err := sendLocation(conn, deviceID, phoneNo); err != nil {
				return
			}
			// 读取应答
			if err := readResponse(conn); err != nil {
				return
			}
		}
	}
}

// sendRegister 发送注册消息
func sendRegister(conn net.Conn, deviceID, phoneNo string) error {
	// 构造注册消息体
	body := make([]byte, 48)
	// 省域ID
	binary.BigEndian.PutUint16(body[0:2], 1)
	// 市县域ID
	binary.BigEndian.PutUint16(body[2:4], 1)
	// 制造商ID
	copy(body[4:9], "TEST0")
	// 终端型号
	copy(body[9:39], "JT808_TEST_DEVICE")
	// 终端ID
	copy(body[39:46], deviceID[:7])
	// 车牌颜色
	body[46] = 1
	// 车牌号
	copy(body[47:], "TEST123")

	msg := buildMessage(0x0100, phoneNo, 1, body)
	_, err := conn.Write(msg)
	if err != nil {
		return err
	}

	atomic.AddInt64(&totalSent, 1)
	return nil
}

// sendAuth 发送鉴权消息
func sendAuth(conn net.Conn, deviceID, phoneNo string) error {
	// 构造鉴权消息体
	authCode := "AUTH123456"
	body := []byte(authCode)

	msg := buildMessage(0x0102, phoneNo, 2, body)
	_, err := conn.Write(msg)
	if err != nil {
		return err
	}

	atomic.AddInt64(&totalSent, 1)
	return nil
}

// sendHeartbeat 发送心跳消息
func sendHeartbeat(conn net.Conn, phoneNo string) error {
	msg := buildMessage(0x0002, phoneNo, 3, []byte{})
	_, err := conn.Write(msg)
	if err != nil {
		return err
	}

	atomic.AddInt64(&totalSent, 1)
	return nil
}

// sendLocation 发送位置信息
func sendLocation(conn net.Conn, deviceID, phoneNo string) error {
	// 构造位置消息体
	body := make([]byte, 28)

	// 报警标志
	binary.BigEndian.PutUint32(body[0:4], 0)

	// 状态（已定位）
	binary.BigEndian.PutUint32(body[4:8], 0x02)

	// 纬度（39.9042 * 100000 = 3990420）
	binary.BigEndian.PutUint32(body[8:12], 3990420)

	// 经度（116.4074 * 100000 = 11640740）
	binary.BigEndian.PutUint32(body[12:16], 11640740)

	// 海拔
	binary.BigEndian.PutUint16(body[16:18], 50)

	// 速度（60km/h）
	binary.BigEndian.PutUint16(body[18:20], 600)

	// 方向
	binary.BigEndian.PutUint16(body[20:22], 90)

	// 时间（BCD编码）
	now := time.Now()
	body[22] = byte((now.Year()%100)/10)<<4 | byte((now.Year()%100)%10)
	body[23] = byte(int(now.Month())/10)<<4 | byte(int(now.Month())%10)
	body[24] = byte(now.Day()/10)<<4 | byte(now.Day()%10)
	body[25] = byte(now.Hour()/10)<<4 | byte(now.Hour()%10)
	body[26] = byte(now.Minute()/10)<<4 | byte(now.Minute()%10)
	body[27] = byte(now.Second()/10)<<4 | byte(now.Second()%10)

	msg := buildMessage(0x0200, phoneNo, 4, body)
	_, err := conn.Write(msg)
	if err != nil {
		return err
	}

	atomic.AddInt64(&totalSent, 1)
	return nil
}

// readResponse 读取应答
func readResponse(conn net.Conn) error {
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		return err
	}

	if n > 0 {
		atomic.AddInt64(&totalReceived, 1)
	}

	return nil
}

// buildMessage 构建JT808消息
func buildMessage(msgID uint16, phoneNo string, flowNo uint16, body []byte) []byte {
	// 消息头
	header := make([]byte, 12)
	binary.BigEndian.PutUint16(header[0:2], msgID)
	binary.BigEndian.PutUint16(header[2:4], uint16(len(body)))
	copy(header[4:10], encodeBCD(phoneNo, 6))
	binary.BigEndian.PutUint16(header[10:12], flowNo)

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

// encodeBCD BCD编码
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

// escape 转义
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

// printStats 打印统计信息
func printStats() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	startTime := time.Now()
	for range ticker.C {
		sent := atomic.LoadInt64(&totalSent)
		received := atomic.LoadInt64(&totalReceived)
		errors := atomic.LoadInt64(&totalErrors)
		connected := atomic.LoadInt64(&connectedDevs)
		elapsed := time.Since(startTime)

		fmt.Printf("\r[%s] Connected: %d | Sent: %d | Received: %d | Errors: %d | Rate: %.1f msg/s",
			elapsed.Truncate(time.Second),
			connected,
			sent,
			received,
			errors,
			float64(sent)/elapsed.Seconds(),
		)
	}
}

// printFinalResults 打印最终结果
func printFinalResults() {
	sent := atomic.LoadInt64(&totalSent)
	received := atomic.LoadInt64(&totalReceived)
	errors := atomic.LoadInt64(&totalErrors)

	fmt.Printf("\n\n=== 测试结果 ===\n")
	fmt.Printf("发送消息数: %d\n", sent)
	fmt.Printf("接收消息数: %d\n", received)
	fmt.Printf("错误数: %d\n", errors)
	fmt.Printf("成功率: %.2f%%\n", float64(received)/float64(sent)*100)
	fmt.Printf("平均速率: %.1f msg/s\n", float64(sent)/float64(testDuration))
}
