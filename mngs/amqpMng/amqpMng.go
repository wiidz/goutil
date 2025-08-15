package amqpMng

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/wiidz/goutil/structs/configStruct"
	"log"
	"strconv"
)

var conn *amqp.Connection // 全局 Connection

// Init 初始化 Connection
func Init(config *configStruct.RabbitMQConfig) (err error) {

	//【1】构建DSN
	dsn := "amqp://" + config.Username + ":" + config.Password + "@" + config.Host + "/"

	//【2】构建DB对象
	log.Println("【rabbitMq-dsn】", dsn)
	conn, err = amqp.Dial(dsn)
	if err != nil {
		log.Println("【rabbit-mq-init-err】", err)
		return
	}
	return
}

// NewRabbitMQ 新建管理对象
func NewRabbitMQ(cfg *configStruct.RabbitMQConfig) (*RabbitMQ, error) {
	if conn == nil {
		return nil, errors.New("[RabbitMQ] conn not initialized, call Init first")
	}
	return &RabbitMQ{
		Config: cfg,
		Conn:   conn,
	}, nil
}

// SetExchange 声明
func (mng *RabbitMQ) SetExchange(channel *amqp.Channel) error {
	args := mng.Config.ExchangeDeclareArgs
	if args == nil {
		args = amqp.Table{}
	}
	if mng.Config.ExchangeType == configStruct.XDelayedMessage {
		args["x-delayed-type"] = "direct" // 或者用 mng.Config.ExchangeType
	}
	return channel.ExchangeDeclare(
		mng.Config.ExchangeName,
		string(mng.Config.ExchangeType),
		mng.Config.IsDurable,
		false, // auto-delete
		false, // internal
		false, // noWait
		args,
	)
}

// DeclareQueue 声明
func (mng *RabbitMQ) DeclareQueue(channel *amqp.Channel) (*amqp.Queue, error) {
	args := mng.Config.QueueDeclareArgs
	if args == nil {
		args = amqp.Table{}
	}

	switch mng.Config.ExchangeType {
	case configStruct.XDelayedMessage:
		// x-delayed-message机制，不用设置 TTL 和死信
	case configStruct.DeadLetterDelay:
		if mng.Config.QueueTTL <= 0 {
			return nil, errors.New("dead letter delay QueueTTL 必须大于0（单位ms）")
		}
		args["x-message-ttl"] = mng.Config.QueueTTL
		if mng.Config.TargetExchangeName == "" || mng.Config.TargetRoutingKey == "" {
			return nil, errors.New("dead letter delay 必须指定 target exchange 和 routing key")
		}
		args["x-dead-letter-exchange"] = mng.Config.TargetExchangeName
		args["x-dead-letter-routing-key"] = mng.Config.TargetRoutingKey
	default:
		if mng.Config.QueueTTL > 0 {
			args["x-message-ttl"] = mng.Config.QueueTTL
		}
	}

	q, err := channel.QueueDeclare(
		mng.Config.QueueName,
		mng.Config.IsDurable,
		false, // auto-delete
		false, // exclusive
		false, // noWait
		args,
	)
	if err != nil {
		return nil, fmt.Errorf("RabbitMQ QueueDeclare: %w", err)
	}
	return &q, nil
}

// BindQueue 队列绑定
func (mng *RabbitMQ) BindQueue(channel *amqp.Channel) error {
	return channel.QueueBind(
		mng.Config.QueueName,
		mng.Config.BindingKey,
		mng.Config.ExchangeName,
		false, // noWait
		mng.Config.QueueBindArgs,
	)
}

// Publish 发布消息（参数 expiration 单位毫秒。reliable 表示用 Publisher Confirm）
func (mng *RabbitMQ) Publish(body string, expiration int, reliable bool) error {
	channel, err := mng.Conn.Channel()
	if err != nil {
		return fmt.Errorf("SetChannel: %w", err)
	}
	defer channel.Close()

	// 声明 Exchange
	if err := mng.SetExchange(channel); err != nil {
		return fmt.Errorf("SetExchange: %w", err)
	}
	// 声明 Queue （可选，通常生产前先确保队列存在）
	if _, err := mng.DeclareQueue(channel); err != nil {
		return fmt.Errorf("DeclareQueue: %w", err)
	}
	// 绑定
	if err := mng.BindQueue(channel); err != nil {
		return fmt.Errorf("BindQueue: %w", err)
	}

	if reliable {
		if err := channel.Confirm(false); err != nil {
			return fmt.Errorf("channel confirm mode: %w", err)
		}
		confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1))
		defer confirmOne(confirms)
	}

	var pub amqp.Publishing
	bodyBytes := []byte(body)
	switch mng.Config.ExchangeType {
	case configStruct.XDelayedMessage:
		pub = amqp.Publishing{
			Headers:      amqp.Table{"x-delay": expiration},
			ContentType:  "text/plain",
			Body:         bodyBytes,
			DeliveryMode: amqp.Persistent,
		}
	default: // 普通/死信/自定义
		pub = amqp.Publishing{
			ContentType:  "text/plain",
			Body:         bodyBytes,
			DeliveryMode: amqp.Persistent,
		}
		if expiration > 0 {
			pub.Expiration = strconv.Itoa(expiration) // 单位 ms，字符串类型
		}
	}

	return channel.Publish(
		mng.Config.ExchangeName,
		mng.Config.RoutingKey,
		false,
		false,
		pub,
	)
}

func confirmOne(confirms <-chan amqp.Confirmation) {
	if confirmed := <-confirms; confirmed.Ack {
		log.Printf("[RabbitMQ] confirmed delivery: %d", confirmed.DeliveryTag)
	} else {
		log.Printf("[RabbitMQ] failed delivery: %d", confirmed.DeliveryTag)
	}
}

// Consume 消费队列（handleFunc 业务处理，每次消费一个消息）
func (mng *RabbitMQ) Consume(consumerTag string, handleFunc func(d amqp.Delivery) error) error {
	channel, err := mng.Conn.Channel()
	if err != nil {
		return fmt.Errorf("SetChannel: %w", err)
	}
	defer channel.Close()

	if err := mng.SetExchange(channel); err != nil {
		return fmt.Errorf("SetExchange: %w", err)
	}
	if _, err := mng.DeclareQueue(channel); err != nil {
		return fmt.Errorf("DeclareQueue: %w", err)
	}
	if err := mng.BindQueue(channel); err != nil {
		return fmt.Errorf("BindQueue: %w", err)
	}

	deliveries, err := channel.Consume(
		mng.Config.QueueName,
		consumerTag,
		false, // autoAck  手动确认
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,
	)
	if err != nil {
		return fmt.Errorf("Consume: %w", err)
	}

	// 启动消费循环
	for d := range deliveries {
		if err := handleFunc(d); err == nil {
			_ = d.Ack(false)
		} else {
			_ = d.Nack(false, true)
			log.Printf("[RabbitMQ][Consume] handleFunc error: %v", err)
		}
	}
	return nil
}
