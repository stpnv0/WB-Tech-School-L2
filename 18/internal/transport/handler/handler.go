package handler

import (
	"calendar/internal/models"
	"calendar/internal/repository"
	"calendar/internal/service"
	"calendar/internal/transport"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(event models.Event) (int, error)
	Update(event models.Event) error
	GetByDay(userID, date string) ([]models.Event, error)
	GetByWeek(userID, dateStr string) ([]models.Event, error)
	GetByMonth(userID, dateStr string) ([]models.Event, error)
	Delete(userID, date string, ID int) error
}

type Handler struct {
	serv Service
	log  *slog.Logger
}

func NewHandler(serv Service, log *slog.Logger) *Handler {
	return &Handler{
		serv: serv,
		log:  log,
	}
}

func (h *Handler) Create(c *gin.Context) {
	var req transport.CreateEventRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := models.Event{
		Date:   req.Date,
		UserID: req.UserID,
		Text:   req.Event,
	}
	id, err := h.serv.Create(event)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": id})
}

func (h *Handler) Update(c *gin.Context) {
	var req transport.UpdateEventRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := models.Event{
		ID:     req.ID,
		UserID: req.UserID,
		Date:   req.Date,
		Text:   req.Event,
	}
	if err := h.serv.Update(event); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "updated"})
}

func (h *Handler) Delete(c *gin.Context) {
	var req transport.DeleteEventRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.serv.Delete(req.UserID, req.Date, req.ID); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "deleted"})
}

func (h *Handler) GetByDay(c *gin.Context) {
	h.getEvents(c, h.serv.GetByDay)
}

func (h *Handler) GetByWeek(c *gin.Context) {
	h.getEvents(c, h.serv.GetByWeek)
}
func (h *Handler) GetByMonth(c *gin.Context) {
	h.getEvents(c, h.serv.GetByMonth)
}

func (h *Handler) getEvents(c *gin.Context, fn func(string, string) ([]models.Event, error)) {
	var req transport.GetEventsRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	events, err := fn(req.UserID, req.Date)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": toEventResponses(events),
	})
}

func handleError(c *gin.Context, err error) {
	switch err {
	case service.IncorrectDateErr:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case repository.NotFoundEventErr:
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}
}

func toEventResponses(events []models.Event) []transport.EventResponse {
	res := make([]transport.EventResponse, 0, len(events))

	for i := range events {
		res = append(res, transport.EventResponse{
			ID:   events[i].ID,
			Date: events[i].Date,
			Text: events[i].Text,
		})
	}

	return res
}
