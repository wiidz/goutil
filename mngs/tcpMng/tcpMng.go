package tcpMng

import (
	"bytes"
	"fmt"
	"github.com/wiidz/goutil/structs/configStruct"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"net"
	"time"
)

// TCPClient 是一个 TCP 客户端结构体
type TCPClient struct {
	conn   net.Conn
	Config *configStruct.TcpConfig
}

// NewTCPClient 创建一个新的 TCP 客户端实例
func NewTCPClient(config *configStruct.TcpConfig) (client *TCPClient, err error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.IP, config.Port))
	if err != nil {
		return nil, err
	}

	// 设置读取操作的超时时间
	if config.ReadTimeOut != 0 {
		_ = conn.SetReadDeadline(time.Now().Add(config.ReadTimeOut))
	}

	// 设置写入操作的超时时间
	if config.WriteTimeOut != 0 {

		_ = conn.SetWriteDeadline(time.Now().Add(config.WriteTimeOut))
	}

	return &TCPClient{conn: conn, Config: config}, nil
}

// SendMessage 发送消息到服务器
func (c *TCPClient) SendMessage(message string) error {
	_, err := fmt.Fprintf(c.conn, message)
	return err
}

// SendBytes 发送消息到服务器
func (c *TCPClient) SendBytes(bytes []byte) error {
	_, err := c.conn.Write(bytes)
	return err
}

// ReceiveMessage 从服务器接收消息
func (c *TCPClient) ReceiveMessage() (returnByte []byte, msgLen int, err error) {
	// 创建一个缓冲区来存储读取的数据
	buffer := make([]byte, 1024)

	// 从连接中读取数据
	msgLen, err = c.conn.Read(buffer)
	if err != nil {
		return nil, 0, err
	}

	// 将读取的数据转换为字符串并打印出来
	//fmt.Println("Received data:", string(buffer[:n]))
	return buffer[:msgLen], msgLen, err

	//for {
	//	n, err := c.conn.Read(tmp)
	//	if err != nil {
	//		if err == io.EOF {
	//			break
	//		}
	//		return nil, err
	//	}
	//
	//	buffer.Write(tmp[:n])
	//
	//	// 假设每条消息以某个特定的结束符标志结束，例如换行符 '\n'
	//	if bytes.Contains(buffer.Bytes(), []byte{'"'}) {
	//		break
	//	}
	//}

	//return buffer.Bytes(), nil
}

// Close 关闭 TCP 连接
func (c *TCPClient) Close() error {
	return c.conn.Close()
}

type Message struct {
	Text string `json:"text"`
}

func GbToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.HZGB2312.NewDecoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
