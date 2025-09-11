package main

import (
	"fmt"
	"net/http"
	"rest-api-event-app/internal/database"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateEvent creates a new event
//
//	@Summary		Creates a new event
//	@Description	Creates a new event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			event	body		database.Event	true	"Event"
//	@Success		201		{object}	database.Event
//	@Router			/events [post]
//	@Security		BearerAuth
func (app *application) createEvent(c *gin.Context) {
	var event database.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := app.GetUserFromContext(c)
	event.OwnerId = user.ID

	result, err := app.models.Events.InsertEvent(c, &event)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetAllEvent return all events
//
//	@Summary		Return all events
//	@Description	Return all events
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	database.Event
//	@Router			/events [get]
func (app *application) getAllEvent(c *gin.Context) {
	events, err := app.models.Events.GetAllEvent()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// GetEvent return a single event
//
//	@Summary		Return a single event
//	@Description	Return a single event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			eventId	path	int	true	"Event ID"
//	@Success		200		{object}	database.Event
//	@Router			/events/{eventId} [get]
func (app *application) getEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("eventId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := app.models.Events.GetEventById(id)

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retreive event"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// UpdateEvent updates an existing event
//
//	@Summary		Updates an existing event
//	@Description	Updates an existing event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			eventId	path		int				true	"Event ID"
//	@Param			event	body		database.Event	true	"Event"
//	@Success		200		{object}	database.Event
//	@Router			/events/{eventId} [put]
//	@Security		BearerAuth
func (app *application) updateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("eventId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	user := app.GetUserFromContext(c)
	existingEvent, err := app.models.Events.GetEventById(id)

	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
	}

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retreive event"})
		return
	}

	if existingEvent.OwnerId != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this event"})
		return
	}

	updatedEvent := &database.Event{}

	if err := c.ShouldBindJSON(updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedEvent.ID = id

	if err := app.models.Events.UpdateEvent(updatedEvent); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}

	c.JSON(http.StatusOK, updatedEvent)
}

// DeleteEvent deletes an existing event
//
//	@Summary		Deletes an existing event
//	@Description	Deletes an existing event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			eventId	path	int	true	"Event ID"
//	@Success		204		{string}	string	"No Content"
//	@Router			/events/{eventId} [delete]
//	@Security		BearerAuth
func (app *application) deleteEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("eventId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := app.GetUserFromContext(c)
	existingEvent, err := app.models.Events.GetEventById(eventId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if existingEvent.OwnerId != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this event"})
		return
	}

	if err := app.models.Events.DeleteEvent(eventId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete"})
	}

	c.JSON(http.StatusNoContent, nil)
}

// AddAttendeeToEvent adds an attendee to an event
//
//	@Summary		Adds an attendee to an event
//	@Description	Adds an attendee to an event
//	@Tags			attendees
//	@Accept			json
//	@Produce		json
//	@Param			eventId	path		int	true	"Event ID"
//	@Param			userId	path		int	true	"User ID"
//	@Success		201		{object}	database.Attendee
//	@Router			/events/{eventId}/attendees/{userId} [post]
//	@Security		BearerAuth
func (app *application) addAttendeeToEvent(c *gin.Context) {
	// Ubah parameter eventId menjadi integer
	eventId, err := strconv.Atoi(c.Param("eventId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "Invalid event id"})
		return
	}

	// Ubah parameter userId menjadi integer
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "Invalid event id"})
		return
	}

	event, err := app.models.Events.GetEventById(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	userToAdd, err := app.models.Users.GetUserById(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	if userToAdd == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user := app.GetUserFromContext(c)

	if event.OwnerId != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to add an attendee"})
		return
	}

	existingAttendee, err := app.models.Attendees.GetByEventAndAttendee(c, event.ID, userToAdd.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendee"})
		return
	}

	if existingAttendee != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Attendee is already exists"})
		return
	}

	attendee := database.Attendee{
		EventId: event.ID,
		UserId:  userToAdd.ID,
	}

	_, err = app.models.Attendees.Insert(c, &attendee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add attendee"})
		return
	}

	c.JSON(http.StatusCreated, attendee)
}

// GetAttendeesForEvent returns all attendees for a given event
//
//	@Summary		Returns all attendees for a given event
//	@Description	Returns all attendees for a given event
//	@Tags			attendees
//	@Accept			json
//	@Produce		json
//	@Param			eventId	path		int	true	"Event ID"
//	@Success		200		{array}		database.User
//	@Router			/events/{eventId}/attendees [get]
func (app *application) getAttendeesForEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("eventId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event id"})
		return
	}

	users, err := app.models.Attendees.GetAttendeesByEvent(c, eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendees for event"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// DeleteAttendeeFromEvent deletes an attendee from an event
//
//	@Summary		Deletes an attendee from an event
//	@Description	Deletes an attendee from an event
//	@Tags			attendees
//	@Accept			json
//	@Produce		json
//	@Param			eventId	path	int	true	"Event ID"
//	@Param			userId	path	int	true	"User ID"
//	@Success		204		{string}	string	"No Content"
//	@Router			/events/{eventId}/attendees/{userId} [delete]
//	@Security		BearerAuth
func (app *application) deleteAttendeeFromEvent(c *gin.Context) {
	// Ubah parameter eventId menjadi integer
	eventId, err := strconv.Atoi(c.Param("eventId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "Invalid event id"})
		return
	}

	// Ubah parameter userId menjadi integer
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "Invalid event id"})
		return
	}

	event, err := app.models.Events.GetEventById(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error:": "Something went wrong"})
		return
	}

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error:": "Event not found"})
		return
	}

	user := app.GetUserFromContext(c)

	if event.OwnerId != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete an attendee from event"})
		return
	}

	err = app.models.Attendees.Delete(c, userId, eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attendees"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetEventsByAttendee returns all events for a given attendee
//
//	@Summary		Returns all events for a given attendee
//	@Description	Returns all events for a given attendee
//	@Tags			attendees
//	@Accept			json
//	@Produce		json
//	@Param			attendeeId	path	int		true	"Attendee ID"
//	@Success		200			{array}	database.Event
//	@Router			/attendees/{attendeeId}/events [get]
func (app *application) getEventByAttendee(c *gin.Context) {
	attendeeId, err := strconv.Atoi(c.Param("attendeeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendee id"})
		return
	}

	events, err := app.models.Attendees.GetEventByAttendee(c, attendeeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get events"})
		return
	}

	c.JSON(http.StatusOK, events)
}
