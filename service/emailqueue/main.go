package emailqueue

import (
	"allopopot-interconnect-service/config"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Attachments struct {
	Filename string `json:"filename"`
	MimeType string `json:"mimetype"`
	Payload  string `json:"payload"`
}

type EmailPayload struct {
	To          []string      `json:"to"`
	Subject     string        `json:"subject"`
	Body        string        `json:"body"`
	Attachments []Attachments `json:"attachments"`
}

var channel *amqp.Channel
var queueEnabled bool = true

func (at *Attachments) SetPayload(payload string) {
	at.Payload = base64.RawStdEncoding.EncodeToString([]byte(payload))
}

func QueueConnect() error {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", config.AMQP_USERNAME, config.AMQP_PASSWORD, config.AMQP_HOST, config.AMQP_PORT))
	if err != nil {
		log.Println("Cannot connect to queue.")
		disableQueueOnConnectionError()
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Println("Failed to open a channel")
		disableQueueOnConnectionError()
		return err
	}
	channel = ch

	err = ch.ExchangeDeclare(config.AMQP_EXCHANGE_NAME, "fanout", true, false, false, false, nil)
	if err != nil {
		log.Println("Failed to declare an exchange")
		disableQueueOnConnectionError()
		return err
	}
	log.Println("Email Dispatch Queue Connected Successfully", "=", config.AMQP_EXCHANGE_NAME)

	return nil
}

func disableQueueOnConnectionError() {
	queueEnabled = false
	log.Println("Email Dispatch Queue has been disabled")
}

func DispatchEmail(ep EmailPayload) error {
	if !queueEnabled {
		return nil
	}
	amqpMessage := new(amqp.Publishing)
	amqpMessage.ContentType = "text/plain"
	epBytes, err := json.Marshal(ep)
	if err != nil {
		log.Fatalln("Failed to parse EmailPayload")
		return err
	}
	amqpMessage.Body = epBytes

	err = channel.Publish(config.AMQP_EXCHANGE_NAME, "", false, false, *amqpMessage)
	if err != nil {
		log.Println("Could not dispatch email", err)
	}
	return nil

}
