package kafka

import (
	"fmt"
	sarama "gopkg.in/Shopify/sarama.v1"
	"os"
	"strings"
)

var producers []producerWithTopic

type producerWithTopic struct {
	Producer sarama.AsyncProducer
	Topic    string
}

// LoadProducers reads the provided environment variables and connects to the
// specified Kafka brokers, preparing asynchronous producers. If we fail to
// connect to any of the Kafka brokers, we return an error.
func LoadProducers() error {
	configs := readConfigs()
	producers = make([]producerWithTopic, len(configs))

	for i, config := range configs {
		producer, err := newAsyncKafkaProducer(config.Brokers)
		if err != nil {
			return err
		}
		producers[i] = producerWithTopic{
			Producer: producer,
			Topic:    config.Topic,
		}
	}

	return nil
}

// Send produces an event on each provided Kafka producer. You should
// call LoadProducers before calling this, or else it won't have anywhere to
// send the message.
func Send(raw string) {
	for _, p := range producers {
		p.Producer.Input() <- &sarama.ProducerMessage{
			Topic: p.Topic,
			Value: sarama.StringEncoder([]byte(raw)),
		}
	}
}

func newAsyncKafkaProducer(brokers []string) (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal

	return sarama.NewAsyncProducer(brokers, config)
}

// readConfigs looks at environment variables to parse Kafka brokers and
// topics. Each producer requires both a broker and topic, although you can
// provide multiple brokers in a comma-delimited list. You can provide an
// arbitrary number of producers, with incrementing integers starting at 1.
//
// Example:
//     KAFKA_BROKERS_1="192.168.100.50:9092"
//     KAFKA_TOPIC_1="t1"
//     KAFKA_BROKERS_2="192.168.100.51:9092,192.168.100.52:9092"
//     KAFKA_TOPIC_2="t2"
func readConfigs() []producerConfig {
	configs := make([]producerConfig, 0)

	for i := 1; ; i += 1 {
		brokers := os.Getenv(fmt.Sprintf("KAFKA_BROKERS_%v", i))
		topic := os.Getenv(fmt.Sprintf("KAFKA_TOPIC_%v", i))

		if brokers == "" || topic == "" {
			break
		}

		configs = append(configs, producerConfig{
			Brokers: strings.Split(brokers, ","),
			Topic:   topic,
		})
	}

	return configs
}

type producerConfig struct {
	Brokers []string
	Topic   string
}
