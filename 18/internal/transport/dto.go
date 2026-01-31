package transport

type CreateEventRequest struct {
	UserID string `json:"user_id" form:"user_id" binding:"required"`
	Date   string `json:"date" form:"date" binding:"required"`
	Event  string `json:"event" form:"event" binding:"required"`
}

type UpdateEventRequest struct {
	ID     int    `json:"id" form:"id" binding:"required"`
	UserID string `json:"user_id" form:"user_id" binding:"required"`
	Date   string `json:"date" form:"date" binding:"required"`
	Event  string `json:"event" form:"event" binding:"required"`
}

type DeleteEventRequest struct {
	ID     int    `json:"id" form:"id" binding:"required"`
	UserID string `json:"user_id" form:"user_id" binding:"required"`
	Date   string `json:"date" form:"date" binding:"required"`
}

type GetEventsRequest struct {
	UserID string `form:"user_id" binding:"required"`
	Date   string `form:"date" binding:"required"`
}

type EventResponse struct {
	ID   int    `json:"id"`
	Date string `json:"date"`
	Text string `json:"event"`
}
