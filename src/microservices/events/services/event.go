package services

import (
	"context"
	"strings"

	"github.com/cinemaabyss/src/microservices/events/models"
	"github.com/segmentio/kafka-go"
)

const (
	consumerGroup = "consumer-group"
)

// EventService ...
type EventService struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

// NewEventService ...
func NewEventService(kafkaBrokers string) *EventService {
	writer := kafka.Writer{
		Addr:     kafka.TCP(kafkaBrokers),
		Balancer: &kafka.LeastBytes{},
		Async:    true,
	}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBrokers},
		GroupID: consumerGroup,
	})
	return &EventService{writer: &writer, reader: reader}
}

// Write ...
func (s *EventService) Write(ctx context.Context, event *models.Event) error {
	return s.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(event.Key),
			Value: []byte(event.Value),
			Topic: mapTypeToTopic(event.Type),
		},
	)
}

// Read ...
func (s *EventService) Read(ctx context.Context) (*models.Event, error) {
	msg, err := s.reader.ReadMessage(context.Background())
	if err != nil {
		return nil, err
	}
	return &models.Event{
		Key:   string(msg.Key),
		Value: string(msg.Value),
		Type:  mapTopicToType(msg.Topic)}, nil

}

// Close ...
func (s *EventService) Close() {
	if s.writer != nil {
		s.writer.Close()
	}
	if s.reader != nil {
		s.reader.Close()
	}
}

func mapTypeToTopic(eventType string) string {
	switch strings.ToLower(eventType) {
	case "user":
		return "user"
	case "payment":
		return "payment"
	case "movie":
		return "movie"
	default:
		return "Unknown"
	}
}

func mapTopicToType(topicName string) string {
	switch topicName {
	case "user":
		return "User"
	case "payment":
		return "Payment"
	case "movie":
		return "Movie"
	default:
		return ""
	}
}
