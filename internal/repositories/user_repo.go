package repositories

import (
	"context"
	"errors"
	"quicknotes/internal/models"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrDuplicateEmail = newRepositoryError(errors.New("email duplicado"))
var ErrInvalidTokenOrUserAlreadyConfirmed = newRepositoryError(errors.New("invalid token or user already confirmed"))

type UserRepository interface {
	Create(ctx context.Context, email, password, token string) (*models.User, string, error)
	ConfirmUserByToken(ctx context.Context, token string) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository função de estanciação
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) Create(ctx context.Context, email, password, hashKey string) (*models.User, string, error) {
	var user models.User
	user.Email = pgtype.Text{String: email, Valid: true}
	user.Password = pgtype.Text{String: password, Valid: true}
	query := `INSERT INTO users(email, password)
			  VALUES($1, $2)
			  RETURNING id, created_at;`
	row := ur.db.QueryRow(ctx, query, user.Email, user.Password)
	if err := row.Scan(&user.ID, &user.CreatedAt); err != nil {
		if strings.Contains(err.Error(), "violates unique constraint") {
			return &user, "", ErrDuplicateEmail
		}
		return &user, "", newRepositoryError(err)
	}
	userToken, err := ur.createConfirmationToken(ctx, &user, hashKey)
	if err != nil {
		return nil, "", err
	}
	return &user, userToken.Token.String, nil
}

func (ur *userRepository) createConfirmationToken(ctx context.Context, user *models.User, token string) (*models.UserConfirmationToken, error) {
	var userToken models.UserConfirmationToken
	userToken.Token = pgtype.Text{String: token, Valid: true}
	userToken.UserID = user.ID
	query := `INSERT INTO users_confirmation_tokens (user_id, token)
	VALUES($1, $2)
	RETURNING id, created_at`

	row := ur.db.QueryRow(ctx, query, userToken.UserID, userToken.Token)
	if err := row.Scan(&userToken.ID, &userToken.CreatedAt); err != nil {
		return nil, err
	}

	return &userToken, nil
}

func (ur *userRepository) ConfirmUserByToken(ctx context.Context, token string) error {
	query := `SELECT u.id u_id, t.id t_id FROM  users u INNER JOIN  users_confirmation_tokens t
	ON u.id = t.user_id
	WHERE u.active = false
	AND t.confirmed = false 
	AND t.token = $1`

	row := ur.db.QueryRow(ctx, query, token)
	var userID, tokenID pgtype.Numeric
	err := row.Scan(&userID, &tokenID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrInvalidTokenOrUserAlreadyConfirmed
		}
		return newRepositoryError(err)
	}
	queryUpdateUser := "UPDATE users SET active = true, updated_at = now() WHERE id = $1"
	_, err = ur.db.Exec(ctx, queryUpdateUser, userID)
	if err != nil {
		return newRepositoryError(err)
	}

	queryUpdateToken := `UPDATE users_confirmation_tokens 
	SET confirmed = true, updated_at = now()
	WHERE id = $1`
	_, err = ur.db.Exec(ctx, queryUpdateToken, tokenID)
	if err != nil {
		return newRepositoryError(err)
	}
	return nil
}

func (ur *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := "SELECT id, email, password, active FROM users WHERE email = $1"
	row := ur.db.QueryRow(ctx, query, email)
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Active); err != nil {
		return nil, newRepositoryError(err)
	}
	return &user, nil
}
