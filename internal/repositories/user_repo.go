package repositories

import (
	"context"
	"errors"
	"quicknotes/internal/models"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrDuplicateEmail = newRepositoryError(errors.New("email duplicado"))

type UserRepository interface {
	Create(ctx context.Context, email, password string) (*models.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository função de estanciação
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) Create(ctx context.Context, email, password string) (*models.User, error) {
	var user models.User
	user.Email = pgtype.Text{String: email, Valid: true}
	user.Password = pgtype.Text{String: password, Valid: true}
	query := `INSERT INTO users(email, password)
			  VALUES($1, $2)
			  RETURNING id, created_at;`
	row := ur.db.QueryRow(ctx, query, user.Email, user.Password)
	if err := row.Scan(&user.ID, &user.CreatedAt); err != nil {
		if strings.Contains(err.Error(), "violates unique constraint") {
			return &user, ErrDuplicateEmail
		}
		return &user, newRepositoryError(err)
	}
	return &user, nil
}
