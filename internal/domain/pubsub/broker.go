package pubsub

import (
    "context"
    "fmt"
    "pubsub-assignment/internal/domain"
    "sync"
)

type Subscribers map[string]domain.Sub

type InMemBroker struct {
	mu   sync.Mutex
	subs map[string]Subscribers
}

func NewBroker() *InMemBroker {
	return &InMemBroker{
		subs: make(map[string]Subscribers),
	}
}

func (b *InMemBroker) Subscribe(s domain.Sub, topic string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.subs[topic] == nil {
		b.subs[topic] = Subscribers{}
	}
	s.AddTopic(topic)

	b.subs[topic][s.GetID()] = s
}

func (b *InMemBroker) Publish(ctx context.Context, e *domain.Envelope) error {
	if e == nil {
		return fmt.Errorf("nil envelope")
	}

	b.mu.Lock()
	subs := b.subs[e.Topic]
	b.mu.Unlock()

	for _, s := range subs {
		if !s.IsActive() {
			fmt.Printf("domain.Sub: %s not active, deleting from subs", s.GetID())
			delete(b.subs[e.Topic], s.GetID())
			return nil
		}
		go func(s domain.Sub) {
			s.Signal(ctx, e)
		}(s)
	}
	return nil
}
