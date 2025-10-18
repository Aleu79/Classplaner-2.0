package repository

import (
	"database/sql"

	"classplanner-2.0/internal/model"
)

type AddressRepository interface {
	GetByUserID(userID int) ([]*model.Address, error)
	CreateAddress(address *model.Address) (*model.Address, error)
	UpdateAddress(id int, address *model.Address) (*model.Address, error)
	DeleteAddress(id int) error
}

type AdressSQL struct {
	db *sql.DB
}

func NewAddressRepository(db *sql.DB) AddressRepository {
	return &AdressSQL{db: db}
}

// Obtener todas las direcciones de un usuario
func (r *AdressSQL) GetByUserID(userID int) ([]*model.Address, error) {
	q := `
		SELECT id, user_id, name, is_primary, city_id, province_id,
		       address1, address2, phone, email, post_code,
		       created_at, updated_at, deleted_at
		FROM addresses
		WHERE user_id = $1 AND deleted_at IS NULL
	`

	rows, err := r.db.Query(q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []*model.Address
	for rows.Next() {
		a := &model.Address{}
		if err := rows.Scan(
			&a.ID, &a.UserID, &a.Name, &a.IsPrimary, &a.CityID, &a.ProvinceID,
			&a.Address1, &a.Address2, &a.Phone, &a.Email, &a.PostCode,
			&a.CreatedAt, &a.UpdatedAt, &a.DeletedAt,
		); err != nil {
			return nil, err
		}
		addresses = append(addresses, a)
	}

	return addresses, nil
}

// Crear una nueva dirección
func (r *AdressSQL) CreateAddress(address *model.Address) (*model.Address, error) {
	q := `
		INSERT INTO addresses (
			user_id, name, is_primary, city_id, province_id,
			address1, address2, phone, email, post_code,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
		RETURNING id
	`

	err := r.db.QueryRow(q,
		address.UserID, address.Name, address.IsPrimary, address.CityID, address.ProvinceID,
		address.Address1, address.Address2, address.Phone, address.Email, address.PostCode,
	).Scan(&address.ID)

	if err != nil {
		return nil, err
	}

	return address, nil
}

// Actualizar una dirección existente
func (r *AdressSQL) UpdateAddress(id int, address *model.Address) (*model.Address, error) {
	q := `
		UPDATE addresses
		SET name = $1, is_primary = $2, city_id = $3, province_id = $4,
		    address1 = $5, address2 = $6, phone = $7, email = $8, post_code = $9,
		    updated_at = NOW()
		WHERE id = $10
	`

	_, err := r.db.Exec(q,
		address.Name, address.IsPrimary, address.CityID, address.ProvinceID,
		address.Address1, address.Address2, address.Phone, address.Email, address.PostCode,
		id,
	)
	if err != nil {
		return nil, err
	}

	address.ID = id
	return address, nil
}

// Eliminar (lógicamente o físicamente) una dirección
func (r *AdressSQL) DeleteAddress(id int) error {
	q := `
		UPDATE addresses
		SET deleted_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(q, id)
	return err
}
