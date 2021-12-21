package rabbitMng

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type Producer struct {
	*RabbitMQ
}

// NewProducer 构建一个生产者
func NewProducer(name, key, exchange string) (producer  *Producer,err error) {
	rabbitM := &RabbitMQ{
		Conn:     conn,
		QueueName:     name,
		Key:      key,
		ExchangeName: exchange,
	}
	_,err = rabbitM.GetQueue()
	if err != nil {
		return
	}
	producer =  &Producer{
		RabbitMQ:rabbitM,
	}
	return
}

// Publish 发布任务
func (producer *Producer) Publish(exchange, exchangeType, routingKey, body, expiration string, reliable bool) error {

	log.Printf("got Channel, declaring %q Exchange (%q)", exchangeType, exchange)
	if err := producer.Channel.ExchangeDeclare(
		exchange,     // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	// Reliable publisher confirms require confirm.select support from the connection.
	if reliable {
		log.Printf("enabling publishing confirms.")
		if err := producer.Channel.Confirm(false); err != nil {
			return fmt.Errorf("Channel could not be put into confirm mode: %s", err)
		}

		confirms := producer.Channel.NotifyPublish(make(chan amqp.Confirmation, 1))
		defer producer.confirmOne(confirms)
	}

	log.Printf("declared Exchange, publishing %dB body (%q)", len(body), body)
	if err := producer.Channel.Publish(
		exchange,   // publish to an exchange
		routingKey, // routing to 0 or more queues
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			Expiration:      expiration,     // 设置2小时7200000  测试五秒
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return fmt.Errorf("Exchange Publish: %s", err)
	}

	return nil
}

// One would typically keep a channel of publishings, a sequence number, and a
// set of unacknowledged sequence numbers and loop until the publishing channel
// is closed.
func (producer *Producer) confirmOne(confirms <-chan amqp.Confirmation) {

	log.Printf("waiting for confirmation of one publishing")

	if confirmed := <-confirms; confirmed.Ack {
		log.Printf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		log.Printf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}
