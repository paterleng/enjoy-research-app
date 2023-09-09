package model

import "github.com/streadway/amqp"

type SendRabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
}

//接收端结构体
type ReceiveRabbitMQ struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	q       amqp.Queue
	msgchan <-chan amqp.Delivery
}
