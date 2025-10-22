package repository

import (
	"context"
	"database/sql"
	"fmt"

	"classplanner/internal/model"

	sq "github.com/Masterminds/squirrel"
)

type AddressRepository interface {
	GetByUserID(ctx context.Context, userID int) ([]*model.Address, error)
	CreateAddress(ctx context.Context, address *model.Address) (*model.Address, error)
	UpdateAddress(ctx context.Context, id int, address *model.Address) (*model.Address, error)
	DeleteAddress(ctx context.Context, id int) error
}

type AddressSQL struct {
	db *sql.DB
	sb sq.StatementBuilderType
}

func NewAddressRepository(db *sql.DB, sb sq.StatementBuilderType) AddressRepository {
	return &AddressSQL{db: db, sb: sb}
}

// Obtener todas las direcciones de un usuario
func (r *AddressSQL) GetByUserID(ctx context.Context, userID int) ([]*model.Address, error) {
	query, args, err := r.sb.
		Select(
			"id", "user_id", "name", "is_primary", "city_id", "province_id",
			"address1", "address2", "phone", "email", "post_code",
			"created_at", "updated_at", "deleted_at",
		).
		From("addresses").
		Where(sq.Eq{"user_id": userID}).
		Where("deleted_at IS NULL").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetByUserID build query error: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("GetByUserID query error: %w", err)
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
			return nil, fmt.Errorf("GetByUserID scan error: %w", err)
		}
		addresses = append(addresses, a)
	}

	return addresses, nil
}

// Crear una nueva dirección
func (r *AddressSQL) CreateAddress(ctx context.Context, address *model.Address) (*model.Address, error) {
	query, args, err := r.sb.
		Insert("addresses").
		Columns(
			"user_id", "name", "is_primary", "city_id", "province_id",
			"address1", "address2", "phone", "email", "post_code", "created_at", "updated_at",
		).
		Values(
			address.UserID, address.Name, address.IsPrimary, address.CityID, address.ProvinceID,
			address.Address1, address.Address2, address.Phone, address.Email, address.PostCode,
			sq.Expr("NOW()"), sq.Expr("NOW()"),
		).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("CreateAddress build query error: %w", err)
	}

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&address.ID)
	if err != nil {
		return nil, fmt.Errorf("CreateAddress query error: %w", err)
	}

	return address, nil
}

// Actualizar una dirección existente
func (r *AddressSQL) UpdateAddress(ctx context.Context, id int, address *model.Address) (*model.Address, error) {
	query, args, err := r.sb.
		Update("addresses").
		SetMap(map[string]interface{}{
			"name":        address.Name,
			"is_primary":  address.IsPrimary,
			"city_id":     address.CityID,
			"province_id": address.ProvinceID,
			"address1":    address.Address1,
			"address2":    address.Address2,
			"phone":       address.Phone,
			"email":       address.Email,
			"post_code":   address.PostCode,
			"updated_at":  sq.Expr("NOW()"),
		}).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("UpdateAddress build query error: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("UpdateAddress exec error: %w", err)
	}

	address.ID = id
	return address, nil
}

// Eliminar (lógicamente o fisicamente) una dirección
func (r *AddressSQL) DeleteAddress(ctx context.Context, id int) error {
	query, args, err := r.sb.
		Update("addresses").
		Set("deleted_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("DeleteAddress build query error: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("DeleteAddress exec error: %w", err)
	}

	return nil
}
