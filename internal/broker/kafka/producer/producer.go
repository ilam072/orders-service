package producer

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	w *kafka.Writer
}

func New(topic string, addr []string) *Producer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(addr...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{w: w}
}

func (p *Producer) Produce(ctx context.Context, message []byte) error {
	err := p.w.WriteMessages(ctx, kafka.Message{
		Value: message,
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *Producer) Close() error {
	return p.w.Close()
}
