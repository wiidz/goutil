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
func NewProducer(exchangeName string, exchangeType ExchangeType) (producer *Producer, err error) {

	var rabbitM *RabbitMQ
	rabbitM, err = NewRabbitMQ(exchangeName, exchangeType)
	if err != nil {
		return
	}

	producer = &Producer{
		RabbitMQ: rabbitM,
	}
	return
}

// Publish 发布任务
func (producer *Producer) Publish(routingKey string,body, expiration string, reliable bool) error {

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
		producer.ExchangeName, // publish to an exchange
		routingKey,   // routing to 0 or more queues
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Persistent, // 1=non-persistent, 2=persistent
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

// PublishDelay 发布延时任务
func (producer *Producer) PublishDelay(routingKey,body, expiration string, reliable bool) error {

	log.Printf("declared Exchange, publishing %dB body (%q)", len(body), body)

	// Reliable publisher confirms require confirm.select support from the connection.
	if reliable {
		log.Printf("enabling publishing confirms.")
		if err := producer.Channel.Confirm(false); err != nil {
			return fmt.Errorf("Channel could not be put into confirm mode: %s", err)
		}

		confirms := producer.Channel.NotifyPublish(make(chan amqp.Confirmation, 1))
		defer producer.confirmOne(confirms)
	}

	err := producer.Channel.Publish(
		//"",                  // exchange 这里为空则不选择 exchange
		producer.ExchangeName,
		routingKey, // routing to 0 or more queues
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Persistent, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			Expiration:      expiration,     // 设置2小时7200000  测试五秒
			// a bunch of application/implementation-specific fields
		},
	)
	if err != nil {
		return fmt.Errorf("Exchange Publish: %s", err)
	}

	return nil
}
