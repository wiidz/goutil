package rabbitMng

import (
	"github.com/streadway/amqp"
	"github.com/wiidz/goutil/structs/configStruct"
	"log"
)

var conn *amqp.Connection

// RabbitMQ rabbit队列管理器
type RabbitMQ struct {
	Conn    *amqp.Connection // 连接，是指物理的连接，一个client与一个server之间有一个连接；
	Channel *amqp.Channel    // 频道，一个连接上可以建立多个channel，可以理解为逻辑上的连接。一般应用的情况下，有一个channel就够用了，不需要创建更多的channel
	// Queue   *amqp.Queue      // 队列，仅是针对接收方（consumer）的，由接收方根据需求创建的。只有队列创建了，交换机才会将新接受到的消息送到队列中，交换机是不会在队列创建之前的消息放进来的。 即在建立队列之前，发出的所有消息都被丢弃了。
	// QueueName    string
	// BindingKey   string       // 绑定（Binding）Exchange与Queue的同时，一般会指定一个binding key ; binding key 并不是在所有情况下都生效，它依赖于Exchange Type，比如fanout类型的Exchange就会无视binding key，而是将消息路由到所有绑定到该Exchange的Queue
	// RoutingKey   string       // 消费者将消息发送给Exchange时，一般会指定一个routing key，当binding key与routing key相匹配时，消息将会被路由到对应的Queue中
	ExchangeName string       // 交换机名称
	ExchangeType ExchangeType // 交换器类型，常用的Exchange Type有 Fanout、 Direct、 Topic、 Headers 这四种。
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

func NewRabbitMQ(exchangeName string, exchangeType ExchangeType) (mng *RabbitMQ, err error) {
	mng = &RabbitMQ{
		Conn:         conn,
		ExchangeName: exchangeName,
		ExchangeType: exchangeType,
	}
	err = mng.SetChannel()
	if err != nil {
		return
	}

	err = mng.SetExchange(nil)
	return
}

// NewRabbitMQDelay
// RabbitMq对交换器，队列，消息都可以声明持久化属性，交换器和队列持久化属性为durable（其属性值为false代表不持久化，属性值为true代表持久化），
// 消息持久化属性为deliveryMode（其属性值为1代表不持久化，属性值为2代表持久化）。
// 在说明durable之前，需要指出的是：
// 1、RabbitMq实例以broker表示，当broker重启时，所有未申明durable的交换器和队列都会被删除。
// 2、RabbitMq中消息都被保存在队列中，所以如果队列都删除了，消息不管有没有设置deliveryMode=2都不管用了。
// 综上所述，可以得出的结论为，交换器未声明durable属性不会影响队列的持久化（但是发送方的producer会被影响，无法正常发送消息）；
// 只声明队列持久化，重启之后消息会丢失；只声明消息的持久化，重启之后消息随队列一起丢失。单单设置消息持久化而不设置队列的持久化没有任何意义。
// 设置 durable=false，如果消费者长期未启动的话，发布者启用队列的时候会出错。这时候需要改 durable=true.
func NewRabbitMQDelay(exchangeName string, exchangeType ExchangeType, isDurable bool) (mng *RabbitMQ, err error) {
	mng = &RabbitMQ{
		Conn:         conn,
		ExchangeName: exchangeName,
		ExchangeType: exchangeType,
	}
	err = mng.SetChannel()
	if err != nil {
		return
	}

	err = mng.SetExchange(amqp.Table{
		"durable":        isDurable,
		"x-delayed-type": string(Direct),
	})
	return
}

// SetChannel 获取信道
func (mng *RabbitMQ) SetChannel() (err error) {
	mng.Channel, err = mng.Conn.Channel()
	return
}

// SetExchange 申明交换机
func (mng *RabbitMQ) SetExchange(arguments amqp.Table) (err error) {
	err = mng.Channel.ExchangeDeclare(
		mng.ExchangeName,         // name of the exchange
		string(mng.ExchangeType), // type
		true,                     // durable 持久化
		false,                    // delete when complete 完成后是否删除
		false,                    // internal
		false,                    // noWait
		arguments,                // arguments
	)
	return
}

// BindQueue 申明并绑定队列到当前channel和exchange上 ttl 是毫秒,-1表示不设置
func (mng *RabbitMQ) BindQueue(queueName, bindingKey string, ttl int32) (queue amqp.Queue, err error) {

	//【3】申明队列
	var args = amqp.Table{}
	if ttl != -1 {
		args["x-message-ttl"] = ttl
	}

	queue, err = mng.Channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		true,
		args)

	//【4】队列绑定至交换机
	err = mng.Channel.QueueBind(
		queueName,
		bindingKey, // Producer
		mng.ExchangeName,
		true,
		nil)

	return
}

// BindDelayQueue 申明并绑定延迟队列到当前的channel信道的exchange上
// queueName 绑定到当前信道的bindingKey上
// 过期后的信息会被推送到targetExchangeName，路由是routingKey
func (mng *RabbitMQ) BindDelayQueue(queueName, bindingKey, targetExchangeName, targetExchangeRoutingKey string, ttl int32) (queue amqp.Queue, err error) {

	var args = amqp.Table{
		"x-dead-letter-exchange":    targetExchangeName,       // 将过期消息发送到执行的exchange中
		"x-dead-letter-routing-key": targetExchangeRoutingKey, // 将过期消息发送到指定的路由中
	}
	if ttl != -1 {
		args["x-message-ttl"] = ttl
	}

	//【1】申明延迟队列
	queue, err = mng.Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		args,      // arguments
	)
	if err != nil {
		return
	}

	//【2】常规队列绑定至交换机
	err = mng.Channel.QueueBind(
		queue.Name,
		bindingKey,
		mng.ExchangeName,
		false,
		nil)

	return
}
