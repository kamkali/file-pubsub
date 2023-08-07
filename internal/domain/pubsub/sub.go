package pubsub

import (
	"fmt"
	"golang.org/x/net/context"
	"pubsub-assignment/internal/domain"
	"sync"
)

type Subscriber struct {
	mu sync.Mutex

	Id           string
	Active       bool
	messages     chan *domain.Envelope
	subscribedTo map[string]struct{}
}

func NewSubscriber(id string) *Subscriber {
	return &Subscriber{
		Id:           id,
		Active:       true,
		messages:     make(chan *domain.Envelope),
		subscribedTo: make(map[string]struct{}),
	}
}

func (s *Subscriber) Signal(ctx context.Context, e *domain.Envelope) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Think about defer recover here
	if s.Active {
		s.messages <- e
	}
}

func (s *Subscriber) Consume(ctx context.Context, handler func(e *domain.Envelope) error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case msg, ok := <-s.messages:
		if !ok {
			return fmt.Errorf("done")
		}
		if err := handler(msg); err != nil {
			return fmt.Errorf("consume error: %w", err)
		}
		return nil
	}
}

func (s *Subscriber) AddTopic(topic string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.subscribedTo[topic] = struct{}{}
}

func (s *Subscriber) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Active = false
	close(s.messages)
}

func (s *Subscriber) GetID() string {
	return s.Id
}

func (s *Subscriber) IsActive() bool {
	return s.Active
}
