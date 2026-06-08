package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:8080", 5*time.Second)
	if err != nil {
		fmt.Printf("连接失败: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("连接成功")

	// 构造注册消息
	phoneNo := "13800138000"
	body := make([]byte, 53)
	binary.BigEndian.PutUint16(body[0:2], 1100)
	binary.BigEndian.PutUint16(body[2:4], 1101)
	copy(body[4:9], "TEST0")
	copy(body[9:39], "JT808-G-2024")
	copy(body[39:46], "TESTDEV")
	body[46] = 1
	copy(body[47:], "TEST1234")

	msg := buildJT808Message(0x0100, phoneNo, 1, body)
	fmt.Printf("发送注册消息: %d bytes, hex=%X\n", len(msg), msg)

	_, err = conn.Write(msg)
	if err != nil {
		fmt.Printf("发送失败: %v\n", err)
		return
	}

	// 读取应答
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("读取应答失败: %v\n", err)
		return
	}

	fmt.Printf("收到应答: %d bytes, hex=%X\n", n, buf[:n])

	// 发送心跳
	hbMsg := buildJT808Message(0x0002, phoneNo, 2, []byte{})
	fmt.Printf("发送心跳消息: %d bytes\n", len(hbMsg))
	conn.Write(hbMsg)

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, err = conn.Read(buf)
	if err != nil {
		fmt.Printf("读取心跳应答失败: %v\n", err)
		return
	}

	fmt.Printf("收到心跳应答: %d bytes, hex=%X\n", n, buf[:n])
	fmt.Println("测试完成!")
}

func buildJT808Message(msgID uint16, phoneNo string, flowNo uint16, body []byte) []byte {
	header := make([]byte, 12)
	binary.BigEndian.PutUint16(header[0:2], msgID)
	binary.BigEndian.PutUint16(header[2:4], uint16(len(body)))
	copy(header[4:10], encodeBCD(phoneNo, 6))
	binary.BigEndian.PutUint16(header[10:12], flowNo)

	var checksum byte
	for _, b := range header {
		checksum ^= b
	}
	for _, b := range body {
		checksum ^= b
	}

	msg := make([]byte, 0, len(header)+len(body)+3)
	msg = append(msg, 0x7E)
	msg = append(msg, escape(header)...)
	msg = append(msg, escape(body)...)
	msg = append(msg, checksum)
	msg = append(msg, 0x7E)
	return msg
}

func encodeBCD(s string, length int) []byte {
	result := make([]byte, length)
	for i := 0; i < length && i*2 < len(s); i++ {
		high, low := byte(0), byte(0)
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
		switch b {
		case 0x7E:
			result = append(result, 0x7D, 0x02)
		case 0x7D:
			result = append(result, 0x7D, 0x01)
		default:
			result = append(result, b)
		}
	}
	return result
}
