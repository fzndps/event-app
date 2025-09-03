package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	v1 := g.Group("/api/v1")
	{
		v1.GET("/events", app.getAllEvent)
		v1.GET("/events/:eventId", app.getEvent)
		v1.GET("/events/:eventId/attendees", app.getAttendeesForEvent)
		v1.GET("/attendees/:attendeeId/event", app.getEventByAttendee)

		v1.POST("/auth/register", app.registerUser)
		v1.POST("/auth/login", app.login)
	}

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	{
		authGroup.POST("/events", app.createEvent)
		authGroup.PUT("/events/:eventId", app.updateEvent)
		authGroup.DELETE("/events/:eventId", app.deleteEvent)
		authGroup.POST("/events/:eventId/attendees/:userId", app.addAttendeeToEvent)
		authGroup.DELETE("/events/:eventId/attendees/:userId", app.deleteAttendeeFromEvent)
	}

	return g
}
