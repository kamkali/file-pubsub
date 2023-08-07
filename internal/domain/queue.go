package domain

import (
	"golang.org/x/net/context"
)

const FileTopic = "lines-topic"

type Envelope struct {
	Topic   string
	Payload any
}

type Sub interface {
	AddTopic(topic string)
	Signal(ctx context.Context, e *Envelope)
	Close()
	Consume(ctx context.Context, handler func(e *Envelope) error) error
	GetID() string
	IsActive() bool
}

//go:generate mockery --name=Sub

type Broker interface {
	Subscribe(s Sub, topic string)
	Publish(ctx context.Context, e *Envelope) error
}

//go:generate mockery --name=Broker

type File struct {
	Name  string
	Lines []string
}

type FileQueueService interface {
	ReadFile(ctx context.Context) (f File, err error)
	WriteFile(ctx context.Context, f File) error
}

//go:generate mockery --name=FileQueueService
