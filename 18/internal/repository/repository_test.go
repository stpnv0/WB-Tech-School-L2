package repository

import (
	"calendar/internal/models"
	"testing"
)

func TestRepository_CreateEvent(t *testing.T) {
	repo := NewRepo()

	event := models.Event{
		UserID: "1",
		Date:   "2023-10-27",
		Text:   "Test Event",
	}

	id, err := repo.CreateEvent(event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if id != 1 {
		t.Errorf("expected ID 1, got %d", id)
	}

	res, err := repo.GetEventByDay("1", "2023-10-27")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res) != 1 {
		t.Fatalf("expected 1 event, got %d", len(res))
	}

	if res[0].Text != "Test Event" {
		t.Errorf("expected 'Test Event', got '%s'", res[0].Text)
	}
}

func TestRepository_UpdateEvent(t *testing.T) {
	repo := NewRepo()

	event := models.Event{
		UserID: "1",
		Date:   "2023-10-27",
		Text:   "Original",
	}
	id, _ := repo.CreateEvent(event)
	event.ID = id

	// Update
	updatedEvent := models.Event{
		UserID: "1",
		Date:   "2023-10-27",
		ID:     id,
		Text:   "Updated",
	}

	err := repo.UpdateEvent(updatedEvent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	res, _ := repo.GetEventByDay("1", "2023-10-27")
	if res[0].Text != "Updated" {
		t.Errorf("expected 'Updated', got '%s'", res[0].Text)
	}

	updatedEvent.ID = 999
	err = repo.UpdateEvent(updatedEvent)
	if err != NotFoundEventErr {
		t.Errorf("expected NotFoundEventErr, got %v", err)
	}
}

func TestRepository_DeleteEvent(t *testing.T) {
	repo := NewRepo()
	event := models.Event{UserID: "1", Date: "2023-10-27", Text: "Delete Me"}
	id, _ := repo.CreateEvent(event)

	err := repo.DeleteEvent("1", "2023-10-27", id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	res, err := repo.GetEventByDay("1", "2023-10-27")
	if err != nil {
		if err != NotFoundEventErr {
			t.Errorf("expected NotFoundEventErr, got %v", err)
		}
	} else if len(res) != 0 {
		t.Errorf("expected 0 events, got %d", len(res))
	}

	// Delete non-existing
	err = repo.DeleteEvent("1", "2023-10-27", 999)
	if err != NotFoundEventErr {
		t.Errorf("expected NotFoundEventErr, got %v", err)
	}
}

func TestRepository_GetEventByWeek(t *testing.T) {
	repo := NewRepo()

	repo.CreateEvent(models.Event{UserID: "1", Date: "2023-10-23", Text: "Monday"})
	repo.CreateEvent(models.Event{UserID: "1", Date: "2023-10-25", Text: "Wednesday"})
	repo.CreateEvent(models.Event{UserID: "1", Date: "2023-10-29", Text: "Sunday"})
	repo.CreateEvent(models.Event{UserID: "1", Date: "2023-10-30", Text: "Next Monday"}) // Should not show

	// Query for a date in the week (e.g., Wednesday)
	res, err := repo.GetEventByWeek("1", "2023-10-25")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res) != 3 {
		t.Errorf("expected 3 events, got %d", len(res))
	}
}

func TestRepository_GetEventByMonth(t *testing.T) {
	repo := NewRepo()
	repo.CreateEvent(models.Event{UserID: "1", Date: "2023-10-01", Text: "Start"})
	repo.CreateEvent(models.Event{UserID: "1", Date: "2023-10-31", Text: "End"})
	repo.CreateEvent(models.Event{UserID: "1", Date: "2023-11-01", Text: "Next Month"})

	res, err := repo.GetEventByMonth("1", "2023-10-15")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res) != 2 {
		t.Errorf("expected 2 events, got %d", len(res))
	}
}
