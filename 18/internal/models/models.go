package models

type Event struct {
	ID     int    `json:"id"`
	Date   string `json:"date"`
	UserID string `json:"user_id"`
	Text   string `json:"text"`
}
