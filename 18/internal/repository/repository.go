package repository

import (
	"calendar/internal/models"
	"errors"
	"sync"
	"time"
)

var NotFoundEventErr = errors.New("event not found")

type Repository struct {
	mu        sync.RWMutex
	eventSets map[string]map[string][]models.Event
	nextID    int
}

func NewRepo() *Repository {
	return &Repository{
		mu:        sync.RWMutex{},
		eventSets: make(map[string]map[string][]models.Event),
		nextID:    1,
	}
}

func (r *Repository) CreateEvent(event models.Event) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	event.ID = r.nextID
	r.nextID++

	_, ok := r.eventSets[event.UserID]
	if !ok {
		r.eventSets[event.UserID] = make(map[string][]models.Event)
	}

	r.eventSets[event.UserID][event.Date] = append(r.eventSets[event.UserID][event.Date], event)
	return event.ID, nil
}

func (r *Repository) UpdateEvent(event models.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	userEvents, ok := r.eventSets[event.UserID]
	if !ok {
		return NotFoundEventErr
	}

	events, ok := userEvents[event.Date]
	if !ok {
		return NotFoundEventErr
	}

	for i := range events {
		if events[i].ID == event.ID {
			events[i] = event

			return nil
		}
	}
	return NotFoundEventErr
}

func (r *Repository) GetEventByDay(userID, date string) ([]models.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	userEvents, ok := r.eventSets[userID]
	if !ok {
		return nil, NotFoundEventErr
	}

	events, ok := userEvents[date]
	if !ok {
		return nil, NotFoundEventErr
	}

	return events, nil
}

func (r *Repository) GetEventByWeek(userID, dateStr string) ([]models.Event, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}

	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := date.AddDate(0, 0, -weekday+1)
	r.mu.RLock()
	defer r.mu.RUnlock()

	userEvents, ok := r.eventSets[userID]
	if !ok {
		return nil, NotFoundEventErr
	}

	var res []models.Event
	for i := 0; i < 7; i++ {
		currDay := monday.AddDate(0, 0, i).Format("2006-01-02")
		if events, ok := userEvents[currDay]; ok {
			res = append(res, events...)
		}
	}

	return res, nil
}

func (r *Repository) GetEventByMonth(userID, dateStr string) ([]models.Event, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}

	year, month, _ := date.Date()
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, date.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	r.mu.RLock()
	defer r.mu.RUnlock()

	userEvents, ok := r.eventSets[userID]
	if !ok {
		return nil, NotFoundEventErr
	}

	var res []models.Event
	for d := startOfMonth; !d.After(endOfMonth); d = d.AddDate(0, 0, 1) {
		currDay := d.Format("2006-01-02")
		if events, ok := userEvents[currDay]; ok {
			res = append(res, events...)
		}
	}

	return res, nil
}

func (r *Repository) DeleteEvent(userID, date string, ID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	userEvents, ok := r.eventSets[userID]
	if !ok {
		return NotFoundEventErr
	}

	events, ok := userEvents[date]
	if !ok {
		return NotFoundEventErr
	}

	for i := range events {
		if events[i].ID == ID {
			events = append(events[:i], events[i+1:]...)
			userEvents[date] = events

			if len(events) == 0 {
				delete(userEvents, date)
				if len(userEvents) == 0 {
					delete(r.eventSets, userID)
				}
			}

			return nil
		}
	}

	return NotFoundEventErr
}
