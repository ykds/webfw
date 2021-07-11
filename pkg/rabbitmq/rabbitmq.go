package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"strconv"
	"time"
)

var (
	DEAD_LETTER_QUEUE    = "DEAD_LETTER_QUEUE"
	DEAD_LETTER_KEY      = "DEAD_LETTER_KEY"
	DEAD_LETTER_EXCHANGE = "DEAD_LETTER_EXCHANGE"

	DELAY_LETTER_QUEUE    = "DELAY_LETTER_QUEUE"
	DELAY_LETTER_KEY      = "DELAY_LETTER_KEY"
	DELAY_LETTER_EXCHANGE = "DELAY_LETTER_EXCHANGE"
)

var conn *amqp.Connection

func NewRabbitMq(user, pass, host, vhost string) {
	url := fmt.Sprintf(
		"amqp://%s:%s@%s%s",
		user,
		pass,
		host,
		vhost,
	)

	var err error
	conn, err = amqp.Dial(url)
	if err != nil {
		log.Fatalf("init rabbitmq failed. err: %v", err)
	}

	DeadLetterQueueDeclare()
	DelayQueueDeclare()
}

func getChannel() (*amqp.Channel, error) {
	if conn == nil {
		return nil, fmt.Errorf("conn is nil")
	}

	return conn.Channel()
}

func DeclareQueue(queueName string, args map[string]interface{}) (amqp.Queue, error) {
	channel, err := getChannel()
	if err != nil {
		return amqp.Queue{}, fmt.Errorf("declare queue fail, err: %v", err)
	}

	queue, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		args,
	)
	return queue, nil
}

func DeclareExchange(exchangeName, kind string) error {
	channel, err := getChannel()
	if err != nil {
		return err
	}

	return channel.ExchangeDeclare(
		exchangeName,
		kind,
		true,
		false,
		false,
		false,
		nil,
	)
}

func Binding(queueName, key, exchangeName string) error {
	channel, err := getChannel()
	if err != nil {
		return err
	}

	return channel.QueueBind(
		queueName,
		key,
		exchangeName,
		false,
		nil,
	)
}

func DeadLetterQueueDeclare() {
	DeclareQueue(DEAD_LETTER_QUEUE, nil)
	DeclareExchange(DEAD_LETTER_EXCHANGE, amqp.ExchangeDirect)
	Binding(DEAD_LETTER_QUEUE, DEAD_LETTER_KEY, DEAD_LETTER_EXCHANGE)
}

func DelayQueueDeclare() {
	DeclareQueue(DELAY_LETTER_QUEUE, map[string]interface{}{
		"x-dead-letter-exchange":    DEAD_LETTER_EXCHANGE,
		"x-dead_letter-routing-key": DEAD_LETTER_KEY,
	})
	DeclareExchange(DELAY_LETTER_EXCHANGE, amqp.ExchangeDirect)
	Binding(DELAY_LETTER_QUEUE, DELAY_LETTER_KEY, DELAY_LETTER_EXCHANGE)
}

func SendMsg(exchange, key, msg string) error {
	channel, err := getChannel()
	if err != nil {
		return err
	}

	return channel.Publish(
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			Body:         []byte(msg),
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			MessageId:    GenMessageId(),
		})
}

func SendDelayMsg(exchange, key, msg string, delay time.Duration) error {
	channel, err := getChannel()
	if err != nil {
		return err
	}

	return channel.Publish(
		DELAY_LETTER_EXCHANGE,
		DELAY_LETTER_KEY,
		false,
		false,
		amqp.Publishing{
			Headers: amqp.Table{
				"return-exchange":    exchange,
				"return-routing-key": key,
			},
			Body:         []byte(msg),
			ContentType:  "application/json",
			DeliveryMode: amqp.Transient,
			MessageId:    GenMessageId(),
			Expiration:   strconv.FormatInt(delay.Milliseconds(), 10),
		},
	)
}

func GenMessageId() string {
	seed := []rune("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM")
	b := make([]rune, 11)
	for i := range b {
		b[i] = seed[rand.Intn(len(seed))]
	}
	return string(b)
}
