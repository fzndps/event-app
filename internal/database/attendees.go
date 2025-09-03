package database

import (
	"context"
	"database/sql"
)

type AtendeeModel struct {
	DB *sql.DB
}

type Attendee struct {
	ID      int `json:"id"`
	UserId  int `json:"user_id"`
	EventId int `json:"event_id"`
}

func (m *AtendeeModel) Insert(ctx context.Context, attendee *Attendee) (*Attendee, error) {
	query := "INSERT INTO attendees (event_id, user_id) VALUES (?, ?)"
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, attendee.EventId, attendee.UserId)

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	attendee.ID = int(id)

	return attendee, nil
}

func (m *AtendeeModel) GetByEventAndAttendee(ctx context.Context, eventId, userId int) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := "SELECT id, event_id, user_id FROM attendees WHERE event_id = ? AND user_id = ?"
	var attendee Attendee
	err := m.DB.QueryRowContext(ctx, query, eventId, userId).Scan(&attendee.ID, &attendee.EventId, &attendee.UserId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &attendee, nil
}

func (m *AtendeeModel) GetAttendeesByEvent(ctx context.Context, eventid int) ([]*User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
		SELECT u.id, u.name, u.email
		FROM users u
		JOIN attendees a ON u.id = a.user_id
		WHERE a.event_id = ?
	`
	rows, err := m.DB.QueryContext(ctx, query, eventid)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (m *AtendeeModel) Delete(ctx context.Context, userId, eventId int) error {
	query := "DELETE FROM attendees WHERE user_id = ? AND event_id = ?"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, eventId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (m *AtendeeModel) GetEventByAttendee(ctx context.Context, attenddeeId int) ([]*Event, error) {
	query := `
		SELECT e.id, e.owner_id, e.name, e.description, e.date, e.location
		FROM events e
		JOIN attendees a ON e.id = a.event_id
		WHERE a.user_id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, attenddeeId)
	if err != nil {
		return nil, err
	}

	var events []*Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location)
		if err != nil {
			return nil, err
		}

		events = append(events, &event)
	}

	return events, nil
}
