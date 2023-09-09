package rabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"web_app/settings"
)

//初始化RabbitMQ
var MQ *amqp.Connection

func Init(rabbitMQ *settings.RabbitMQConfig) (err error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitMQ.UserName, rabbitMQ.PassWord, rabbitMQ.Host, rabbitMQ.Post)
	conn, err := amqp.Dial(url)
	if err != nil {
		return
	} else {
		MQ = conn
	}
	return
}
