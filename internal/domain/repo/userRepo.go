package repo

import (
	"context"
	"database/sql"
	"resapi/internal/domain/models"
)

type UserRepoInterfase interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, ID int) error
}

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (ur *UserRepo) CreateUser(ctx context.Context, user *models.User) error {

	query := `INSERT INTO users (name, password, email, created_at) VALUES ($1, $2, $3, $4) RETURNING id`

	return ur.db.QueryRowContext(ctx, query, user.Name, user.Password, user.Email, user.CreatedAt).Scan(&user.ID)

}

func (ur *UserRepo) GetUserByID(ctx context.Context, id int) (*models.User, error) {

	query := `SELECT id, name, password, email, created_at FROM users WHERE id = $1`

	user := &models.User{}

	err := ur.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Password,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepo) UpdateUser(ctx context.Context, user *models.User) error {
	query := `UPDATE users SET name=$1, email=$2, password=$3, created_at=$4 WHERE id=$5`

	res, err := ur.db.ExecContext(ctx, query,
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.ID,
	)
	if err != nil {
		return err
	}

	rowAff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowAff == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (ur *UserRepo) DeleteUser(ctx context.Context, ID int) error {
	query := `DELETE FROM users WHERE id=$1`

	result, err := ur.db.ExecContext(ctx, query, ID)

	if err != nil {
		return err
	}

	rowAff, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowAff == 0 {
		return sql.ErrNoRows
	}
	return nil
}
