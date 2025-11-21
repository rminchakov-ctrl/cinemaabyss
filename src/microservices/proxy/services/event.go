package services

import httpclient "github.com/cinemaabyss/src/microservices/proxy/http_client"

// EventService ...
type EventService struct {
	client httpclient.Client
}

// NewEventService ...
func NewEventService(baseURL string) *EventService {
	return &EventService{
		client: *httpclient.NewClient(baseURL),
	}
}
