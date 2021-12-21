package rabbitMng

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type Consumer struct {
	*RabbitMQ
	ConsumerTag string
	done        chan error
	//HandleFunc  func(deliveries <-chan amqp.Delivery, done chan error)
}

// NewConsumer 获取消费者
func NewConsumer(exchangeName string, exchangeType ExchangeType, queueName, routingKey string) (consumer *Consumer, err error) {
	rabbitM := &RabbitMQ{
		Conn:         conn,
		QueueName:    queueName,
		RoutingKey:   routingKey,
		ExchangeName: exchangeName,
		ExchangeType: exchangeType,
	}
	_, err = rabbitM.GetQueue()
	if err != nil {
		return
	}

	//【1】构建消费者
	consumer = &Consumer{
		RabbitMQ:    rabbitM,
		done:        make(chan error),
		//HandleFunc:  handleFunc,
	}
	go func() {
		// 通知信道关闭
		fmt.Printf("closing: %s", <-consumer.Conn.NotifyClose(make(chan *amqp.Error)))
	}()

	return consumer, nil
}

func (consumer *Consumer) Start(consumerTag string,handleFunc func(deliveries <-chan amqp.Delivery, done chan error)) (err error) {
	//【1】开始消费
	log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", consumerTag)

	consumer.ConsumerTag = consumerTag

	var deliveries <-chan amqp.Delivery
	deliveries, err = consumer.Channel.Consume(
		consumer.QueueName,   // name
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

	go handleFunc(deliveries, consumer.done)
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

// handle 具体处理逻辑
//func handle(deliveries <-chan amqp.Delivery, done chan error) {
//	for d := range deliveries {
//		log.Printf(
//			"got %dB delivery: [%v] %q",
//			len(d.Body),
//			d.DeliveryTag,
//			d.Body,
//		)
//		_ = d.Ack(false)
//	}
//	log.Printf("handle: deliveries channel closed")
//	done <- nil
//}