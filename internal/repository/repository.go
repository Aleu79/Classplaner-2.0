package repository

import "database/sql"

// Repository centraliza el acceso a los distintos repositorios
type Repository struct {
	db             *sql.DB
	UserStorage    UserRepository
	AddressStorage AddressRepository
}

// New crea una instancia de Repository con todas las dependencias inicializadas
func New(db *sql.DB) *Repository {
	return &Repository{
		db:             db,
		UserStorage:    &userSQL{db: db},
		AddressStorage: &AdressSQL{db: db},
	}
}
