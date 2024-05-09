package amqpMng

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/wiidz/goutil/structs/configStruct"
	"log"
	"strconv"
)

var conn *amqp.Connection // 连接

// 备忘：
// queue的名称没啥关系，和怎么消费怎么生产没有任何关系
// 声明的exchange和queue只能被一个channel占有，另外新起的话没有用（单线程）
// 所以我们这次修改amqp的代码，将queue保存回来
// 也就是说，conn是共用一个，但是每个mnq对应一个exchange和一个queue（简单模式就不管他了）

// RabbitMQ rabbit队列管理器
type RabbitMQ struct {
	Config *Config

	Conn    *amqp.Connection // 连接，是指物理的连接，一个client与一个server之间有一个连接；
	Channel *amqp.Channel    // 频道，一个连接上可以建立多个channel，可以理解为逻辑上的连接。一般应用的情况下
	Queue   *amqp.Queue      // 队列，仅是针对接收方（consumer）的，由接收方根据需求创建的。只有队列创建了，交换机才会将新接受到的消息送到队列中，交换机是不会在队列创建之前的消息放进来的。 即在建立队列之前，发出的所有消息都被丢弃了。
}

type ExchangeType string

const Fanout ExchangeType = "fanout"                     // 【fanout】类型的Exchange路由会把所有发送到该Exchange的消息路由到所有与它绑定的Queue中。
const Direct ExchangeType = "direct"                     // 【direct】类型的Exchange路由会把消息路由到那些binding key与routing key完全匹配的queue中。
const Topic ExchangeType = "topic"                       // 【topic】类型的Exchange路由会把消息路由到binding key与routing key相匹配的Queue中。
const XDelayedMessage ExchangeType = "x-delayed-message" // 【XDelayedMessage】延迟插件rabbitMq-dsn

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

// Config 启动的参数
type Config struct {
	IsDurable bool // 是否是持久化

	ExchangeName        string       // 交换机名称
	ExchangeType        ExchangeType // 交换机类型
	ExchangeDeclareArgs amqp.Table   // 交换机补充数据

	QueueName        string     // 队列名称
	QueueTTL         int32      // 秒数
	QueueDeclareArgs amqp.Table // 队列补充数据
	QueueBindArgs    amqp.Table // 队列补充数据

	// 我们默认这俩一样所以在设置的时候也搞一样就好了
	BindingKey string // 绑定（Binding）Exchange与Queue的同时，一般会指定一个binding key ; binding key 并不是在所有情况下都生效，它依赖于Exchange Type，比如fanout类型的Exchange就会无视binding key，而是将消息路由到所有绑定到该Exchange的Queue
	RoutingKey string // 生产者将消息发送给Exchange时，一般会指定一个routing key，当binding key与routing key相匹配时，消息将会被路由到对应的Queue中

	// 只有delay类型有这两个值
	TargetExchangeName string // 在delayExchange中过期后，到哪个exchange去
	TargetRoutingKey   string // 同上，告诉exchange推入哪个队列（当然也需要queue的bindingKey相同才能推进去）

	// 以下是消费者的配置
	//ConsumerTag string
}

func NewRabbitMQ(config *Config) (mng *RabbitMQ, err error) {

	//【1】构建管理器
	mng = &RabbitMQ{
		Conn:   conn,
		Config: config,
	}

	return
}

// SetChannel 获取信道
func (mng *RabbitMQ) SetChannel() (err error) {

	var channel *amqp.Channel
	channel, err = mng.Conn.Channel()

	if err != nil {
		err = fmt.Errorf("RabbitMQ SetChannel: %s", err)
	}

	mng.Channel = channel

	return
}

// SetExchange 申明交换机
func (mng *RabbitMQ) SetExchange() (err error) {

	args := mng.Config.ExchangeDeclareArgs
	if mng.Config.ExchangeType == XDelayedMessage {
		if args == nil {
			args = amqp.Table{}
		}
		args["durable"] = mng.Config.IsDurable
		args["x-delayed-type"] = string(Direct)
	}

	err = mng.Channel.ExchangeDeclare(
		mng.Config.ExchangeName,         // name of the exchange
		string(mng.Config.ExchangeType), // type
		mng.Config.IsDurable,            // durable 持久化
		false,                           // delete when complete 完成后是否删除
		false,                           // internal
		false,                           // noWait
		args,                            // arguments
	)
	if err != nil {
		err = fmt.Errorf("RabbitMQ channel SetExchange: %s", err)
	}
	return
}

// DeclareQueue 声明队列（可以简单理解为预创建一个队列，确定他在不在，能不能投递进去）
// 生产者需要先声明，不需要绑定
func (mng *RabbitMQ) DeclareQueue() (queue *amqp.Queue, err error) {

	//【1】初始化参数
	var args = mng.Config.QueueDeclareArgs
	if args == nil {
		args = amqp.Table{}
	}

	//【2】设置ttl
	if mng.Config.QueueTTL != -1 {
		args["x-message-ttl"] = mng.Config.QueueTTL
	}

	//【3】设置其他的（主要是delay）
	if mng.Config.ExchangeType == XDelayedMessage {
		if mng.Config.TargetExchangeName == "" || mng.Config.TargetRoutingKey == "" {
			err = errors.New("绑定队列失败，参数不齐全")
		}
		args["x-dead-letter-exchange"] = mng.Config.TargetExchangeName  // 将过期消息发送到执行的exchange中
		args["x-dead-letter-routing-key"] = mng.Config.TargetRoutingKey // 将过期消息发送到指定的路由中
	}

	temp, err := mng.Channel.QueueDeclare(
		mng.Config.QueueName,
		mng.Config.IsDurable,
		false,
		false,
		true,
		args)

	if err != nil {
		err = fmt.Errorf("RabbitMQ channel QueueDeclare: %s", err)
	}
	queue = &temp

	return
}

// BindQueue 绑定队列到当前channel和exchange上
func (mng *RabbitMQ) BindQueue() (err error) {

	err = mng.Channel.QueueBind(
		mng.Config.QueueName,
		mng.Config.BindingKey, // Producer
		mng.Config.ExchangeName,
		true,
		mng.Config.QueueBindArgs)

	if err != nil {
		err = fmt.Errorf("RabbitMQ channel QueueBind: %s", err)
	}

	return
}

// Publish 推入task
// Q：在使用RabbitMQ时，一些常见的操作会产生异常，例如连接错误、频道未打开、交换机/队列不存在等等。在您的情况下，当您重复使用同一 channel 发送消息时，可能会发生频道未打开的异常，这是由于 channel 在某些情况下可能会被关闭，例如由于网络问题或由于其他进程关闭了该 channel。
// A：因此，建议您在每次发送消息时都创建一个新的 channel，并在完成后关闭该 channel。这样可以确保 channel 的状态是正确的，并且能够避免意外的异常。
//
//	当然，您也可以在程序中添加异常处理程序，以处理可能出现的异常情况。例如，如果 channel 关闭，则可以在 catch 语句块中重新打开一个新的 channel。但是，这种方法可能会增加代码复杂性，并且不如直接创建新的 channel 安全可靠。
func (mng *RabbitMQ) Publish(body string, expiration int, reliable bool) (err error) {

	//【1】打开信道
	err = mng.SetChannel()
	if err != nil {
		return
	}
	defer mng.Channel.Close()

	//【2】声明交换机
	err = mng.SetExchange()
	if err != nil {
		return
	}

	//【3】声明队列
	_, err = mng.DeclareQueue()
	if err != nil {
		return
	}

	//【4】推入
	// Reliable publisher confirms require confirm.select support from the connection.
	if reliable {
		log.Printf("enabling publishing confirms.")
		err = mng.Channel.Confirm(false)
		if err != nil {
			return fmt.Errorf("channel could not be put into confirm mode: %s", err)
		}

		confirms := mng.Channel.NotifyPublish(make(chan amqp.Confirmation, 1))
		defer mng.confirmOne(confirms)
	}

	log.Printf("[%s] get new push\n", mng.Config.QueueName)
	log.Printf("%dB body (%q)\n", len(body), body)

	//【5】区分一下
	var publishing amqp.Publishing
	exprStr := strconv.Itoa(expiration)
	if mng.Config.ExchangeType == XDelayedMessage {
		publishing = amqp.Publishing{
			Headers: amqp.Table{
				"x-delay": exprStr,
			},
			ContentType: "text/plain",
			Body:        []byte(body),
			//Expiration:      expiration,     // ms，设置2小时7200000  测试五秒
		}
	} else {
		publishing = amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Persistent, // 1=non-persistent, 2=persistent
			Priority:        0,               // 0-9
			Expiration:      exprStr,         // 设置2小时7200000  测试五秒
		}
	}

	err = mng.Channel.Publish(
		mng.Config.ExchangeName, // publish to an exchange
		mng.Config.RoutingKey,   // routing to 0 or more queues
		false,                   // mandatory 没有的话会强制创建queue
		false,                   // immediate
		publishing,
	)

	if err != nil {
		return fmt.Errorf("exchange Publish: %s", err)
	}

	return nil
}

// One would typically keep a channel of publishings, a sequence number, and a
// set of unacknowledged sequence numbers and loop until the publishing channel
// is closed.
func (mng *RabbitMQ) confirmOne(confirms <-chan amqp.Confirmation) {

	log.Printf("waiting for confirmation of one publishing")

	if confirmed := <-confirms; confirmed.Ack {
		log.Printf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		log.Printf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}

// Consume 开始消费（可以开多个消费者）
func (mng *RabbitMQ) Consume(consumerTag string, handleFunc func(delivery amqp.Delivery) error) (err error) {

	//【1】打开信道
	err = mng.SetChannel()
	if err != nil {
		return
	}
	defer mng.Channel.Close()

	//【2】声明交换机
	err = mng.SetExchange()
	if err != nil {
		return
	}

	//【3】声明队列
	_, err = mng.DeclareQueue()
	if err != nil {
		return
	}

	//【4】绑定队列
	err = mng.BindQueue()

	if err != nil {
		return
	}

	//【4】开始消费
	log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", consumerTag)

	var deliveries <-chan amqp.Delivery

	var ch *amqp.Channel
	ch, err = conn.Channel()

	deliveries, err = ch.Consume(
		mng.Config.QueueName, // name
		consumerTag,          // consumerTag,
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
				handledError := handleFunc(delivery)
				if handledError == nil {
					_ = delivery.Ack(false)
				}
			}
		}
	}()

	<-forever

	return
}

//// Shutdown 关闭conn的回调
//func (mng *RabbitMQ) Shutdown() error {
//
//	// will close() the deliveries channel
//	if err := mng.Channel.Cancel(mng.ConsumerTag, true); err != nil {
//		return fmt.Errorf("Consumer cancel failed: %s", err)
//	}
//
//	if err := mng.Conn.Close(); err != nil {
//		return fmt.Errorf("AMQP connection close error: %s", err)
//	}
//
//	defer log.Printf("AMQP shutdown OK")
//
//	// wait for handle() to exit
//	return <-consumer.done
//}
