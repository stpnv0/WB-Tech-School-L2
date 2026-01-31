package router

import (
	"github.com/gin-gonic/gin"
)

type handler interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByDay(c *gin.Context)
	GetByWeek(c *gin.Context)
	GetByMonth(c *gin.Context)
}

func InitRouter(h handler) *gin.Engine {
	r := gin.Default()

	r.POST("/create_event", h.Create)
	r.POST("/update_event", h.Update)
	r.POST("/delete_event", h.Delete)

	r.GET("/events_for_day", h.GetByDay)
	r.GET("/events_for_week", h.GetByWeek)
	r.GET("/events_for_month", h.GetByMonth)

	return r
}
