package main

import (
	"context"
	"fmt"
	"os"
	"wb-l0/internal/broker/kafka/producer"
	"wb-l0/internal/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	p := producer.New(cfg.KafkaConfig.Topic, cfg.KafkaConfig.Brokers)

	if err := sendValidMessage(context.Background(), p); err != nil {
		fmt.Println(err)
	}

}

func sendValidMessage(ctx context.Context, p *producer.Producer) error {
	message, err := os.ReadFile("C:/Users/adami/GolandProjects/wb-l0/assets/valid_order.json")
	if err != nil {
		return err
	}

	fmt.Println(string(message))

	if err = p.Produce(ctx, message); err != nil {
		return err
	}

	return nil
}
