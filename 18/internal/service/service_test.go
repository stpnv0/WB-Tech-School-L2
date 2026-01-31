package service

import (
	"calendar/internal/models"
	"errors"
	"log/slog"
	"os"
	"testing"
)

type MockRepo struct {
	CreateFunc     func(event models.Event) (int, error)
	UpdateFunc     func(event models.Event) error
	GetByDayFunc   func(userID, date string) ([]models.Event, error)
	GetByWeekFunc  func(userID, dateStr string) ([]models.Event, error)
	GetByMonthFunc func(userID, dateStr string) ([]models.Event, error)
	DeleteFunc     func(userID, date string, ID int) error
}

func (m *MockRepo) CreateEvent(event models.Event) (int, error) { return m.CreateFunc(event) }
func (m *MockRepo) UpdateEvent(event models.Event) error        { return m.UpdateFunc(event) }
func (m *MockRepo) GetEventByDay(userID, date string) ([]models.Event, error) {
	return m.GetByDayFunc(userID, date)
}
func (m *MockRepo) GetEventByWeek(userID, dateStr string) ([]models.Event, error) {
	return m.GetByWeekFunc(userID, dateStr)
}
func (m *MockRepo) GetEventByMonth(userID, dateStr string) ([]models.Event, error) {
	return m.GetByMonthFunc(userID, dateStr)
}
func (m *MockRepo) DeleteEvent(userID, date string, ID int) error {
	return m.DeleteFunc(userID, date, ID)
}

func TestService_Create(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mockRepo := &MockRepo{
		CreateFunc: func(event models.Event) (int, error) {
			return 1, nil
		},
	}
	svc := NewService(mockRepo, logger)

	// Invalid Date
	_, err := svc.Create(models.Event{Date: "rubbish"})
	if !errors.Is(err, IncorrectDateErr) {
		t.Errorf("expected IncorrectDateErr, got %v", err)
	}

	// Valid Date
	id, err := svc.Create(models.Event{Date: "2023-10-27"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 1 {
		t.Errorf("expected 1, got %d", id)
	}
}

func TestService_GetByDay(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mockRepo := &MockRepo{
		GetByDayFunc: func(userID, date string) ([]models.Event, error) {
			return []models.Event{}, nil
		},
	}
	svc := NewService(mockRepo, logger)

	// Invalid Date
	_, err := svc.GetByDay("1", "invalid")
	if !errors.Is(err, IncorrectDateErr) {
		t.Errorf("expected IncorrectDateErr, got %v", err)
	}

	// Valid
	_, err = svc.GetByDay("1", "2023-10-27")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
