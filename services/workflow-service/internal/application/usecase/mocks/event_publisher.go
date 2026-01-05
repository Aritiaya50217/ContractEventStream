package mocks

import "github.com/stretchr/testify/mock"

type EventPublisher struct {
	mock.Mock
}

func (m *EventPublisher) Publish(topic string, payload []byte) error {
	args := m.Called(topic, payload)
	return args.Error(0)
}
