package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	httpclient "github.com/cinemaabyss/src/microservices/proxy/http_client"
	"github.com/cinemaabyss/src/microservices/proxy/models"
)

const (
	moviePath = "/api/movies"
)

// MovieService ...
type MovieService struct {
	momolithClient *httpclient.Client
	movieClient    *httpclient.Client
	moviesPrc      int
}

// NewMoviesService ...
func NewMoviesService(monolithURL string, movieURL string, moviesPrc int) *MovieService {
	return &MovieService{
		momolithClient: httpclient.NewClient(monolithURL),
		movieClient:    httpclient.NewClient(movieURL),
		moviesPrc:      moviesPrc,
	}
}

// GetMovies ...
func (s *MovieService) GetMovies(ctx context.Context) ([]*models.Movie, error) {
	client := s.getClient()

	req, err := client.NewGetRequest(ctx, moviePath)
	if err != nil {
		return nil, fmt.Errorf("error create request: %w", err)
	}

	body, err := client.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching movies data: %w", err)
	}

	var movies []models.Movie
	if err = json.Unmarshal(body, &movies); err != nil {
		return nil, fmt.Errorf("error decoding movies response: %w", err)
	}

	retVal := make([]*models.Movie, 0, len(movies))
	for _, movie := range movies {
		retVal = append(retVal, &movie)
	}

	return retVal, nil
}

// GetMovieByID ...
func (s *MovieService) GetMovieByID(ctx context.Context, ID string) (*models.Movie, error) {

	client := s.getClient()

	req, err := client.NewGetRequest(ctx, fmt.Sprintf("%s/%s", moviePath, ID))
	if err != nil {
		return nil, fmt.Errorf("error create request: %w", err)
	}

	body, err := client.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching movies data: %w", err)
	}

	var movie models.Movie
	if err := json.Unmarshal(body, &movie); err != nil {
		return nil, fmt.Errorf("error decoding movie response: %w", err)
	}

	return &movie, nil
}

func (s *MovieService) getClient() *httpclient.Client {
	if rand.Intn(100) >= s.moviesPrc {
		return s.movieClient
	}
	return s.momolithClient
}

// CreateMovie ...
func (s *MovieService) CreateMovie(ctx context.Context, movie *models.Movie) (*models.Movie, error) {
	if movie == nil {
		return nil, fmt.Errorf("empty movie")
	}
	client := s.getClient()

	req, err := client.NewPostRequest(ctx, moviePath, movie)
	if err != nil {
		return nil, fmt.Errorf("error create request: %w", err)
	}

	body, err := client.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching movies data: %w", err)
	}

	var newMovie models.Movie
	if err := json.Unmarshal(body, &newMovie); err != nil {
		return nil, fmt.Errorf("error decoding movie response: %w", err)
	}

	return &newMovie, nil
}

/*
func (s *MonolithService) UpdateSensor(ctx context.Context, sensor *models.Sensor) (*models.Sensor, error) {
	if sensor == nil {
		return nil, fmt.Errorf("empty sensor")
	}

	req, err := s.client.NewPutRequest(ctx, fmt.Sprintf("%s/%d", path, sensor.ID), sensor)
	if err != nil {
		return nil, fmt.Errorf("error create request: %w", err)
	}

	body, err := s.client.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching sensors data: %w", err)
	}

	var newSensor models.Sensor
	if err := json.Unmarshal(body, &newSensor); err != nil {
		return nil, fmt.Errorf("error decoding sensor response: %w", err)
	}

	return &newSensor, nil
}

func (s *MonolithService) DeleteSensor(ctx context.Context, ID int64) error {
	req, err := s.client.NewDeleteRequest(ctx, fmt.Sprintf("%s/%d", path, ID))
	if err != nil {
		return fmt.Errorf("error create request: %w", err)
	}

	_, err = s.client.DoRequest(req)
	return err
}

func (s *MonolithService) UpdateSensorValue(ctx context.Context, sensor *models.Sensor) error {
	if sensor == nil {
		return fmt.Errorf("empty sensor")
	}

	req, err := s.client.NewPatchRequest(ctx, fmt.Sprintf("%s/%d", path, sensor.ID), sensor)
	if err != nil {
		return fmt.Errorf("error create request: %w", err)
	}

	_, err = s.client.DoRequest(req)
	return err
}
*/
