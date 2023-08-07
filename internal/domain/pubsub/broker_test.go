package pubsub

import (
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "pubsub-assignment/internal/domain"
    "pubsub-assignment/internal/domain/mocks"
    "testing"
)

func TestInMemBroker_Subscribe(t *testing.T) {
	tests := []struct {
		name       string
		setSubMock func(t *testing.T) domain.Sub
		topic      string
	}{
		{
			name: "successfully subscribes a sub",
			setSubMock: func(t *testing.T) domain.Sub {
				s := mocks.NewSub(t)
				s.On("AddTopic", mock.Anything).Once()
				s.On("GetID").Return("id-1")
				return s
			},
			topic: "test-topic",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBroker()
			s := tt.setSubMock(t)
			b.Subscribe(s, tt.topic)

			subs, ok := b.subs[tt.topic]
			assert.True(t, ok)
			assert.Len(t, subs, 1)
			assert.Equal(t, subs[s.GetID()], s)
		})
	}
}
