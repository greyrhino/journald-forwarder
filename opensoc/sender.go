package opensoc

import (
	"encoding/json"
	"fmt"
	sarama "gopkg.in/Shopify/sarama.v1"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func newAsyncKafkaProducer(broker string) (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal

	return sarama.NewAsyncProducer(brokers, config)
}

func SendEvent(entry string, broker string, topic string) {
	producer, err := newAsyncKafkaProducer(broker)
	if err != nil {
		log.Println("oh snap there was an error")
		return 1
	}
	producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder([]byte(entry)),
	}
}
