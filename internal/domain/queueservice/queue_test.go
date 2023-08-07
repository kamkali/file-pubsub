package queueservice

import (
    "errors"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "golang.org/x/net/context"
    "pubsub-assignment/internal/domain"
    "pubsub-assignment/internal/domain/mocks"
    "reflect"
    "testing"
)

func TestQueueService_ReadFile(t *testing.T) {
	tests := []struct {
		name       string
		setSubMock func(t *testing.T, sub *mocks.Sub)
		wantF      domain.File
		wantErr    bool
	}{
		{
			name: "successfully consumes from queue",
			setSubMock: func(t *testing.T, sub *mocks.Sub) {
				sub.On("Consume", mock.Anything, mock.Anything).
					Return(nil)
			},
			wantF:   domain.File{},
			wantErr: false,
		},
		{
			name: "consume error queue",
			setSubMock: func(t *testing.T, sub *mocks.Sub) {
				sub.On("Consume", mock.Anything, mock.Anything).
					Return(errors.New("err"))
			},
			wantF:   domain.File{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subMock := mocks.NewSub(t)
			if tt.setSubMock != nil {
				tt.setSubMock(t, subMock)
			}
			q := New(nil, subMock)
			gotF, err := q.ReadFile(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotF, tt.wantF) {
				t.Errorf("ReadFile() gotF = %v, want %v", gotF, tt.wantF)
			}
		})
	}
}

func TestQueueService_WriteFile(t *testing.T) {
	ctx := context.Background()
	f := domain.File{
		Name:  "test",
		Lines: []string{"one", "two"},
	}
	e := &domain.Envelope{
		Topic:   domain.FileTopic,
		Payload: f,
	}
	brokerMock := mocks.NewBroker(t)
	q := New(brokerMock, nil)

	t.Run("happy path", func(t *testing.T) {
		brokerMock.On("Publish", ctx, e).
			Return(nil).Once()
		err := q.WriteFile(ctx, f)
		assert.NoError(t, err)
	})
	t.Run("broker error", func(t *testing.T) {
		brokerMock.On("Publish", ctx, e).
			Return(errors.New("err")).Once()
		err := q.WriteFile(ctx, f)
		assert.Error(t, err)
	})
}
