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
	//Queue   *amqp.Queue      // 队列，仅是针对接收方（consumer）的，由接收方根据需求创建的。只有队列创建了，交换机才会将新接受到的消息送到队列中，交换机是不会在队列创建之前的消息放进来的。 即在建立队列之前，发出的所有消息都被丢弃了。
	// QueueName    string
	BindingKey   string       // 绑定（Binding）Exchange与Queue的同时，一般会指定一个binding key ; binding key 并不是在所有情况下都生效，它依赖于Exchange Type，比如fanout类型的Exchange就会无视binding key，而是将消息路由到所有绑定到该Exchange的Queue
	RoutingKey   string       // 消费者将消息发送给Exchange时，一般会指定一个routing key，当binding key与routing key相匹配时，消息将会被路由到对应的Queue中
	ExchangeName string       // 交换机名称
	ExchangeType ExchangeType // 交换器类型，常用的Exchange Type有 Fanout、 Direct、 Topic、 Headers 这四种。
}

type ExchangeType string

const Fanout ExchangeType = "fanout" // 【fanout】类型的Exchange路由会把所有发送到该Exchange的消息路由到所有与它绑定的Queue中。
const Direct ExchangeType = "direct" // 【direct】类型的Exchange路由会把消息路由到那些binding key与routing key完全匹配的queue中。
const Topic ExchangeType = "topic"   // 【topic】类型的Exchange路由会把消息路由到binding key与routing key相匹配的Queue中。

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


func NewRabbitMQ(exchangeName string,exchangeType ExchangeType,bindingKey,routingKey string) (mng *RabbitMQ,err error){
	mng = &RabbitMQ{
		Conn:         conn,
		BindingKey:   bindingKey,
		RoutingKey:   routingKey,
		ExchangeName: exchangeName,
		ExchangeType: exchangeType,
	}
	err = mng.SetChannel()
	return
}

// SetChannel 获取信道
func (mng *RabbitMQ) SetChannel() (err error) {
	mng.Channel, err = mng.Conn.Channel()
	return
}

// SetExchange 申明交换机
func (mng *RabbitMQ) SetExchange(arguments *amqp.Table) (err error) {
	err = mng.Channel.ExchangeDeclare(
		mng.ExchangeName,         // name of the exchange
		string(mng.ExchangeType), // type
		true,                     // durable 持久化
		false,                    // delete when complete 完成后是否删除
		false,                    // internal
		false,                    // noWait
		*arguments,               // arguments
	)
	return
}

// BindQueue 申明并绑定队列
func (mng *RabbitMQ) BindQueue(queueName string) (queue amqp.Queue, err error) {

	//【2】申明交换机
	err = mng.SetExchange(nil)
	if err != nil {
		return
	}

	//【3】申明队列
	queue, err = mng.Channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		true,
		nil)

	//【4】队列绑定至交换机
	err = mng.Channel.QueueBind(
		queueName,
		mng.BindingKey, // Producer
		mng.ExchangeName,
		true,
		nil)

	return
}

// GetDelayQueue 获取延迟队列
// 声明延时队列队列，该队列中消息如果过期，就将消息发送到交换器上，交换器就分发消息到普通队列
func (mng *RabbitMQ) GetDelayQueue(queueName string) (queue amqp.Queue, err error) {

	//【2】申明交换机
	err = mng.SetExchange(&amqp.Table{
		"x-dead-letter-exchange": "logs", //当消息过期时把消息发送到logs这个交换器
		//"x-dead-letter-routing-key": "hah.t",
	})
	if err != nil {
		return
	}

	//【3】申明队列
	queue, err = mng.Channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		true,
		nil)

	//【4】队列绑定至交换机
	err = mng.Channel.QueueBind(
		queueName,
		mng.BindingKey, // Producer
		mng.ExchangeName,
		true,
		nil)

	return
}

//
///**
// * @func: UserCouponDelay 用户的优惠券延迟
// * @author Hank
// * @date   2019-06-24
// */
//func UserCouponDelay(user_coupon_id, delay_time string) {
//
//	conn, err := getMqConn()
//	defer conn.Close()
//
//	ch, err := conn.Channel()
//	defer ch.Close()
//
//	// 将消息发送到延时队列上
//	err = ch.Publish(
//		"",             // exchange 这里为空则不选择 exchange
//		"coupon_delay", // routing key
//		false,          // mandatory
//		false,          // immediate
//		amqp.Publishing{
//			ContentType: "text/plain",
//			Body:        []byte(user_coupon_id),
//			Expiration:  delay_time, // 设置2小时7200000  测试五秒
//		})
//
//	failOnError(err, "Failed to publish a message")
//
//}
