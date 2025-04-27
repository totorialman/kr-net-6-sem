package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/totorialman/kr-net-6-sem/internal/utils"
	"github.com/totorialman/kr-net-6-sem/internal/storage"
)

const (
	KafkaAddr  = "kafka:29092"
	KafkaTopic = "segments"
)

func ReadFromKafka() error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// создание consumer-а
	consumer, err := sarama.NewConsumer([]string{KafkaAddr}, config)
	if err != nil {
		return fmt.Errorf("error creating consumer: %w", err)
	}
	defer consumer.Close()

	// подключение consumer-а к топика
	partitionConsumer, err := consumer.ConsumePartition(KafkaTopic, 0, sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("error opening topic: %w", err)
	}
	defer partitionConsumer.Close()

	// бесконечный цикл чтения
	for {
		select {
		case message := <-partitionConsumer.Messages():
			segment := utils.Segment{}
			if err := json.Unmarshal(message.Value, &segment); err != nil {
				fmt.Printf("Error reading from kafka: %v", err)
			}
			fmt.Printf("%+v\n", segment) // выводим в консоль прочитанный сегмент
			storage.AddSegment(segment)
		case err := <-partitionConsumer.Errors():
			fmt.Printf("Error: %s\n", err.Error())
		}
	}
}

func WriteToKafka(segment utils.Segment) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	// создание producer-а
	producer, err := sarama.NewSyncProducer([]string{KafkaAddr}, config)
	if err != nil {
		return fmt.Errorf("error creating producer: %w", err)
	}
	defer producer.Close()

	// превращение segment в сообщение для Kafka
	segmentString, _ := json.Marshal(segment)
	message := &sarama.ProducerMessage{
		Topic: KafkaTopic,
		Value: sarama.StringEncoder(segmentString),
	}

	// отправка сообщения
	_, _, err = producer.SendMessage(message)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
