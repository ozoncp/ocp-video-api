package producer

import (
	"encoding/json"
	"github.com/Shopify/sarama"
)

type Sender interface {
	Init() error
	Close() error
	Send(e Event) error
}

func NewSaramaSender(broker string, topic string) Sender {
	return &sender{
		broker: broker,
		topic:  topic,
	}
}

type sender struct {
	broker string
	topic  string
	prod   sarama.SyncProducer
}

func (s *sender) Init() error {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	prod, err := sarama.NewSyncProducer([]string{s.broker}, config)
	if err != nil {
		return err
	}

	s.prod = prod

	return nil
}

func (s *sender) Close() error {
	return s.prod.Close()
}

func (s *sender) Send(e Event) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	message := sarama.ProducerMessage{
		Topic:     s.topic,
		Partition: -1,
		Value:     sarama.StringEncoder(data),
	}
	_, _, err = s.prod.SendMessage(&message)

	return err
}
