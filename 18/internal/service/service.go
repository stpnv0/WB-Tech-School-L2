package service

import (
	"calendar/internal/models"
	"errors"
	"log/slog"
	"time"
)

var IncorrectDateErr = errors.New("Date have incorrect format, expected YYYY-MM-DD")

type Repository interface {
	CreateEvent(event models.Event) (int, error)
	UpdateEvent(event models.Event) error
	GetEventByDay(userID, date string) ([]models.Event, error)
	GetEventByWeek(userID, dateStr string) ([]models.Event, error)
	GetEventByMonth(userID, dateStr string) ([]models.Event, error)
	DeleteEvent(userID, date string, ID int) error
}

type Service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) Create(event models.Event) (int, error) {
	if err := validateDate(event.Date); err != nil {
		return 0, err
	}

	return s.repo.CreateEvent(event)
}

func (s *Service) Update(event models.Event) error {
	if err := validateDate(event.Date); err != nil {
		return err
	}

	return s.repo.UpdateEvent(event)
}

func (s *Service) GetByDay(userID, date string) ([]models.Event, error) {
	if err := validateDate(date); err != nil {
		return nil, err
	}

	return s.repo.GetEventByDay(userID, date)
}

func (s *Service) GetByWeek(userID, date string) ([]models.Event, error) {
	if err := validateDate(date); err != nil {
		return nil, err
	}

	return s.repo.GetEventByWeek(userID, date)
}

func (s *Service) GetByMonth(userID, date string) ([]models.Event, error) {
	if err := validateDate(date); err != nil {
		return nil, err
	}

	return s.repo.GetEventByMonth(userID, date)
}

func (s *Service) Delete(userID, date string, ID int) error {
	if err := validateDate(date); err != nil {
		return err
	}

	return s.repo.DeleteEvent(userID, date, ID)
}

func validateDate(date string) error {
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return IncorrectDateErr
	}
	return nil
}
