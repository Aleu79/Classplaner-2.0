package repository

import "database/sql"

type CalendarRepositoy interface {
}

type CalendarSQL struct {
	db *sql.DB
}

func NewCalendarRepository(db *sql.DB) CalendarRepositoy {
	return CalendarSQL{db: db}
}
