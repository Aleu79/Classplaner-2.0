package repository

import "database/sql"

// Repository centraliza el acceso a los distintos repositorios
type Repository struct {
	db                *sql.DB
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
	return &Repository{
		db:              db,
		UserStorage:     &userSQL{db: db},
		TasksStorage:    &TasksSQL{db: db},
		ClassesStorage:  &ClassesSQL{db: db},
		CommentStorage:  &CommentSQL{db: db},
		CalendarStorage: &CalendarSQL{db: db},
	}
}
