package database

import (
	"database/sql"
	"time"
)

type Models struct {
	Users     UserModel
	Events    EventModel
	Attendees AtendeeModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:     UserModel{DB: db},
		Events:    EventModel{DB: db},
		Attendees: AtendeeModel{DB: db},
	}
}

var QueryTimeoutDuration = time.Second * 3
