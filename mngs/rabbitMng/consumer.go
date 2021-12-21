package rabbitMng

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type Consumer struct {
	*RabbitMQ
	BindingKey  string
	ConsumerTag string
	done        chan error
	//HandleFunc  func(deliveries <-chan amqp.Delivery, done chan error)
}

// NewConsumer 获取消费者
func NewConsumer(exchangeName string, exchangeType ExchangeType,isDelay bool) (consumer *Consumer, err error) {

	//【1】创建mq
	var rabbitM *RabbitMQ
	if isDelay {
		rabbitM, err = NewRabbitMQDelay(exchangeName, exchangeType)
	} else {
		rabbitM, err = NewRabbitMQ(exchangeName, exchangeType)
	}
	if err != nil {
		return
	}

	//【3】构建消费者
	consumer = &Consumer{
		RabbitMQ:   rabbitM,
		done:       make(chan error),
		//HandleFunc:  handleFunc,
	}

	go func() {
		// 通知信道关闭
		fmt.Printf("closing: %s", <-consumer.Conn.NotifyClose(make(chan *amqp.Error)))
	}()

	return consumer, nil
}

// Start 开始消费
// Tips：记得在外部先绑定队列
func (consumer *Consumer) Start(queueName, consumerTag string, handleFunc func(delivery amqp.Delivery) error ) (err error) {

	//【1】开始消费
	log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", consumerTag)
	consumer.ConsumerTag = consumerTag

	var deliveries <-chan amqp.Delivery
	deliveries, err = consumer.Channel.Consume(
		queueName,            // name
		consumer.ConsumerTag, // consumerTag,
		false,                // noAck
		false,                // exclusive
		false,                // noLocal
		false,                // noWait
		nil,                  // arguments
	)
	if err != nil {
		return
	}

	forever := make(chan bool)
	var delivery amqp.Delivery
	go func() {
		for {
			select {
			case delivery = <-deliveries:
				err := handleFunc(delivery)
				if err == nil {
					_ = delivery.Ack(false)
				}
			}
		}
	}()

	<-forever

	return
}

// Shutdown 关闭conn的回调
func (consumer *Consumer) Shutdown() error {

	// will close() the deliveries channel
	if err := consumer.Channel.Cancel(consumer.ConsumerTag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := consumer.Conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-consumer.done
}
