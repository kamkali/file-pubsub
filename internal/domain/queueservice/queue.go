package queueservice

import (
    "fmt"
    "golang.org/x/net/context"
    "pubsub-assignment/internal/domain"
)

type QueueService struct {
	broker domain.Broker
	sub    domain.Sub
}

func New(broker domain.Broker, sub domain.Sub) *QueueService {
	return &QueueService{
		broker: broker,
		sub:    sub,
	}
}

func (q *QueueService) ReadFile(ctx context.Context) (f domain.File, err error) {
	var file domain.File
	err = q.sub.Consume(ctx, func(e *domain.Envelope) error {
		f, ok := e.Payload.(domain.File)
		if !ok {
			return fmt.Errorf("cannot retrieve file from the queue")
		}
		file = f
		return nil
	})
	if err != nil {
		return domain.File{}, err
	}

	return file, err
}

func (q *QueueService) WriteFile(ctx context.Context, f domain.File) error {
	e := &domain.Envelope{
		Topic:   domain.FileTopic,
		Payload: f,
	}

	return q.broker.Publish(ctx, e)
}
