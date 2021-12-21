package rabbitMng

import (
	"github.com/streadway/amqp"
	"github.com/wiidz/goutil/structs/configStruct"
	"log"
)

var conn *amqp.Connection

// RabbitMQ rabbit队列管理器
type RabbitMQ struct {
	Conn         *amqp.Connection // 连接，是指物理的连接，一个client与一个server之间有一个连接；
	Channel      *amqp.Channel    // 频道，一个连接上可以建立多个channel，可以理解为逻辑上的连接。一般应用的情况下，有一个channel就够用了，不需要创建更多的channel
	Queue        *amqp.Queue      // 队列，仅是针对接收方（consumer）的，由接收方根据需求创建的。只有队列创建了，交换机才会将新接受到的消息送到队列中，交换机是不会在队列创建之前的消息放进来的。 即在建立队列之前，发出的所有消息都被丢弃了。
	QueueName    string
	Key          string
	ExchangeName string // 交换机
}

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

// GetQueue 获取队列
func (mng *RabbitMQ) GetQueue() (queue amqp.Queue, err error) {

	channel, _ := mng.Conn.Channel()

	//【2】申明交换机
	err = mng.Channel.ExchangeDeclare(
		mng.ExchangeName, // name of the exchange
		"direct",         // type
		true,             // durable
		false,            // delete when complete
		false,            // internal
		false,            // noWait
		nil,              // arguments
	)
	if err != nil {
		return
	}

	//【3】申明队列
	queue, err = channel.QueueDeclare(
		mng.QueueName,
		true,
		false,
		false,
		true,
		nil)

	//【4】队列绑定至交换机
	err = channel.QueueBind(
		mng.QueueName,
		mng.Key,
		mng.ExchangeName,
		true,
		nil)

	mng.Channel = channel
	mng.Queue = &queue
	
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
