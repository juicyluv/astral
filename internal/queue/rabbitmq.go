package queue

import (
	"fmt"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Queue struct {
	logger *zap.SugaredLogger
	cfg    *Config
	conn   *amqp.Connection
	ch     *amqp.Channel
}

func NewQueue(logger *zap.SugaredLogger, cfg *Config) (*Queue, error) {
	q := Queue{
		logger: logger,
		cfg:    cfg,
	}

	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.User, cfg.Password, cfg.Host, cfg.Port)

	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, err
	}
	q.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	q.ch = ch

	_, err = q.ch.QueueDeclare(
		cfg.Name,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &q, nil
}

func (q *Queue) Dispatch(msg []byte) error {
	return q.ch.Publish(
		"",
		q.cfg.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		},
	)
}

func (q *Queue) Close() error {
	err := q.ch.Close()
	if err != nil {
		return err
	}
	return q.conn.Close()
}
