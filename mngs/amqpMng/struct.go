package amqpMng

import (
	"github.com/streadway/amqp"
	"github.com/wiidz/goutil/structs/configStruct"
)

// RabbitMQ Exchange 类型
type ExchangeType string

const (
	Fanout          ExchangeType = "fanout"
	Direct          ExchangeType = "direct"
	Topic           ExchangeType = "topic"
	XDelayedMessage ExchangeType = "x-delayed-message"
	DeadLetterDelay ExchangeType = "dead_letter_delay"
)

// Core 配置结构体（可和你的 configStruct 对接）
type Config struct {
	IsDurable bool

	ExchangeName        string
	ExchangeType        ExchangeType
	ExchangeDeclareArgs amqp.Table

	QueueName        string
	QueueTTL         int // 毫秒
	QueueDeclareArgs amqp.Table
	QueueBindArgs    amqp.Table

	BindingKey string
	RoutingKey string

	TargetExchangeName string // 用于死信延迟和普通死信
	TargetRoutingKey   string

	// 可扩展更多参数
}

// 管理器
type RabbitMQ struct {
	Config *configStruct.RabbitMQConfig
	Conn   *amqp.Connection
}
