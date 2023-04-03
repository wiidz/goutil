package amqpMng

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/wiidz/goutil/structs/configStruct"
	"log"
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

	//Conn         *amqp.Connection // 连接，是指物理的连接，一个client与一个server之间有一个连接；
	Channel *amqp.Channel // 频道，一个连接上可以建立多个channel，可以理解为逻辑上的连接。一般应用的情况下
	Queue   *amqp.Queue   // 队列，仅是针对接收方（consumer）的，由接收方根据需求创建的。只有队列创建了，交换机才会将新接受到的消息送到队列中，交换机是不会在队列创建之前的消息放进来的。 即在建立队列之前，发出的所有消息都被丢弃了。
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
	ExchangeName        string       // 交换机名称
	ExchangeType        ExchangeType // 交换机类型
	ExchangeDeclareArgs amqp.Table   // 交换机补充数据

	QueueName        string     // 队列名称
	QueueTTL         int32      // 秒数
	QueueDeclareArgs amqp.Table // 队列补充数据
	QueueBindArgs    amqp.Table // 队列补充数据

	// 我们默认这俩一样所以在设置的时候也搞一样就好了
	BindingKey string // 绑定（Binding）Exchange与Queue的同时，一般会指定一个binding key ; binding key 并不是在所有情况下都生效，它依赖于Exchange Type，比如fanout类型的Exchange就会无视binding key，而是将消息路由到所有绑定到该Exchange的Queue
	RoutingKey string // 消费者将消息发送给Exchange时，一般会指定一个routing key，当binding key与routing key相匹配时，消息将会被路由到对应的Queue中

	// 只有delay类型有这两个值
	TargetExchangeName string // 在delayExchange中过期后，到哪个exchange去
	TargetRoutingKey   string // 同上，告诉exchange推入哪个队列（当然也需要queue的bindingKey相同才能推进去）

	// 以下是消费者的配置
	//ConsumerTag string
}

func NewRabbitMQ(config *Config) (mng *RabbitMQ, err error) {

	//【1】构建管理器
	mng = &RabbitMQ{
		Config: config,
	}

	//【2】打开信道
	err = mng.SetChannel()
	if err != nil {
		return
	}

	//【3】声明交换机
	err = mng.SetExchange()
	if err != nil {
		return
	}

	//【4】声明并绑定队列
	if config.ExchangeType == XDelayedMessage {
		mng.Queue, err = mng.BindDelayQueue()
	} else {
		mng.Queue, err = mng.BindQueue()
	}
	if err != nil {
		return
	}

	return
}

// SetChannel 获取信道
func (mng *RabbitMQ) SetChannel() (err error) {

	mng.Channel, err = conn.Channel()

	if err != nil {
		err = fmt.Errorf("RabbitMQ SetChannel: %s", err)
	}

	return
}

// SetExchange 申明交换机
func (mng *RabbitMQ) SetExchange() (err error) {
	err = mng.Channel.ExchangeDeclare(
		mng.Config.ExchangeName,         // name of the exchange
		string(mng.Config.ExchangeType), // type
		true,                            // durable 持久化
		false,                           // delete when complete 完成后是否删除
		false,                           // internal
		false,                           // noWait
		mng.Config.ExchangeDeclareArgs,  // arguments
	)
	if err != nil {
		err = fmt.Errorf("RabbitMQ channel SetExchange: %s", err)
	}
	return
}

// BindQueue 申明并绑定队列到当前channel和exchange上 ttl 是毫秒,-1表示不设置
func (mng *RabbitMQ) BindQueue() (queue *amqp.Queue, err error) {

	//【3】申明队列
	var args = mng.Config.QueueDeclareArgs
	if mng.Config.QueueTTL != -1 {
		args["x-message-ttl"] = mng.Config.QueueTTL
	}

	log.Printf("mng.Channel", mng.Channel)

	temp, err := mng.Channel.QueueDeclare(
		mng.Config.QueueName,
		true,
		false,
		false,
		true,
		args)

	if err != nil {
		err = fmt.Errorf("RabbitMQ channel QueueDeclare: %s", err)
		return
	}

	queue = &temp

	//【4】队列绑定至交换机
	// 发布延时任务，注意千万不能绑定队列，不然会直接推到队列里去？
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

// BindDelayQueue 申明并绑定延迟队列到当前的channel信道的exchange上
func (mng *RabbitMQ) BindDelayQueue() (queue *amqp.Queue, err error) {

	//【1】验一下参数
	if mng.Config.TargetExchangeName == "" || mng.Config.TargetRoutingKey == "" {
		err = errors.New("绑定队列失败，参数不齐全")
	}

	var args = amqp.Table{
		"x-dead-letter-exchange":    mng.Config.TargetExchangeName, // 将过期消息发送到执行的exchange中
		"x-dead-letter-routing-key": mng.Config.TargetRoutingKey,   // 将过期消息发送到指定的路由中
	}
	if mng.Config.QueueTTL != -1 {
		args["x-message-ttl"] = mng.Config.QueueTTL
	}

	//【1】申明延迟队列
	*queue, err = mng.Channel.QueueDeclare(
		mng.Config.QueueName, // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		args,                 // arguments
	)

	if err != nil {
		err = fmt.Errorf("RabbitMQ channel QueueDeclare: %s", err)
		return
	}

	//【2】常规队列绑定至交换机
	err = mng.Channel.QueueBind(
		queue.Name,
		mng.Config.BindingKey,
		mng.Config.ExchangeName,
		false,
		nil)

	if err != nil {
		err = fmt.Errorf("RabbitMQ channel QueueBind: %s", err)
	}

	return
}

// Publish 推入task
func (mng *RabbitMQ) Publish(body string, expiration int, reliable bool) (err error) {

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

	log.Printf("declared Exchange, publishing %dB body (%q)", len(body), body)

	//【2】区分一下
	var publishing amqp.Publishing
	if mng.Config.ExchangeType == XDelayedMessage {
		publishing = amqp.Publishing{
			Headers: amqp.Table{
				"x-delay": expiration,
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
			DeliveryMode:    amqp.Persistent,    // 1=non-persistent, 2=persistent
			Priority:        0,                  // 0-9
			Expiration:      string(expiration), // 设置2小时7200000  测试五秒
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

// Start 开始消费（可以开多个消费者）
func (mng *RabbitMQ) Start(consumerTag string, handleFunc func(delivery amqp.Delivery) error) (err error) {

	//【1】开始消费
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
