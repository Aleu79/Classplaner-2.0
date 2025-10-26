package repository

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

// Repository centraliza el acceso a los distintos repositorios
type Repository struct {
	db *sql.DB
	sb sq.StatementBuilderType

	UserStorage       UserRepository
	AddressStorage    AddressRepository
	TasksStorage      TasksRepository
	ClassesStorage    ClassesRepository
	SubmissionStorage SubmissionRepository
	CommentStorage    CommentRepository
	CalendarStorage   CalendarRepositoy
}

// New crea una instancia de Repository con todas las dependencias inicializadas
func New(db *sql.DB) *Repository {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar) // ideal para Postgres

	return &Repository{
		db: db,
		sb: builder,

		UserStorage:       &userSQL{db: db, sb: builder},
		AddressStorage:    &AddressSQL{db: db, sb: builder},
		TasksStorage:      &TasksSQL{db: db, sb: builder},
		ClassesStorage:    &ClassesSQL{db: db, sb: builder},
		SubmissionStorage: &SubmissionSQL{db: db, sb: builder},
		CommentStorage:    &CommentSQL{db: db, sb: builder},
		CalendarStorage:   &CalendarSQL{db: db, sb: builder},
	}
}
