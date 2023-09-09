package controller

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"strconv"
	"web_app/dao/rabbitMQ"
)

func NewSendRabbitMQ(reception uint) *SendRabbitMQ {
	mysend := new(SendRabbitMQ)
	var err error
	mysend.conn = rabbitMQ.MQ
	mysend.ch, err = mysend.conn.Channel()
	queuename := "->" + strconv.Itoa(int(reception))
	mysend.q, err = mysend.ch.QueueDeclare(
		queuename, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	//cu.FailOnError(err, "Failed to declare a queue")
	fmt.Println(err)
	return mysend
}

type SendRabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
}

func SendData(Data string, sr *SendRabbitMQ) {
	var err error
	err = sr.ch.Publish(
		"",        // exchange
		sr.q.Name, // routing key
		false,     // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(Data),
		})
	fmt.Println(err)
	log.Printf(" [x] Sent %s", Data)
}

func (sr *SendRabbitMQ) Close() {
	sr.ch.Close()
	sr.conn.Close()
}

type ReceiveRabbitMQ struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	q       amqp.Queue
	msgchan <-chan amqp.Delivery
}

func NewReceiveRabbitMQ(queuename string) *ReceiveRabbitMQ {
	myreceive := new(ReceiveRabbitMQ)
	var err error
	myreceive.conn = rabbitMQ.MQ
	myreceive.ch, err = myreceive.conn.Channel()
	myreceive.q, err = myreceive.ch.QueueDeclare(
		queuename, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	err = myreceive.ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	myreceive.msgchan, err = myreceive.ch.Consume(
		myreceive.q.Name, // queue
		"",               // consumer
		false,            // auto-ack
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // args
	)
	fmt.Println(err)
	return myreceive
}

func ReceiveData(rr *ReceiveRabbitMQ, msg string) {
	forever := make(chan bool)
	go func() {
		for d := range rr.msgchan {
			fmt.Println(msg)
			d.Ack(true)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (rr *ReceiveRabbitMQ) Close() {
	rr.ch.Close()
	rr.conn.Close()
}
