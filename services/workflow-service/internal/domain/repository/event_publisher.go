package repository

type EventPublisher interface {
	Publish(topic string, payload []byte) error
}
