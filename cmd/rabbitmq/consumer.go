package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/juicyluv/astral/configs"
	"github.com/juicyluv/astral/internal/mail"
	"github.com/juicyluv/astral/internal/queue"
	"github.com/streadway/amqp"
)

func main() {
	if err := configs.LoadConfigs("configs/dev.yml"); err != nil {
		panic(err)
	}

	cfg := queue.NewConfig()

	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.User, cfg.Password, cfg.Host, cfg.Port)

	conn, err := amqp.Dial(dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	// Start receiving messages
	messages, err := ch.Consume(
		cfg.Name, // queue name
		"",       // consumer name
		true,     // autoAck
		false,    // exclusive
		false,    // noLocal
		false,    // noWait
		nil,      // args
	)
	if err != nil {
		panic(err)
	}

	forever := make(chan struct{})

	// Read and parse every message
	go func() {
		for message := range messages {
			msg, err := DeserializeMessage(message.Body)
			if err != nil {
				fmt.Println("could not deserialize message: ", err.Error())
			}

			err = mail.SendEmail(msg.EmailTo, msg.Subject, msg.Mime, string(msg.Message))
			if err != nil {
				fmt.Printf("Could not send message to %s. Error: %s\n", msg.EmailTo, err.Error())
				continue
			}
			fmt.Printf("Message has been sent to %s\n", msg.EmailTo)
		}
	}()

	fmt.Println("Connected to RabbitMQ, listening for messages.")
	<-forever
}

func DeserializeMessage(b []byte) (mail.Message, error) {
	var msg mail.Message
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&msg)
	return msg, err
}
