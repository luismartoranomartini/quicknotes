package repositories

import (
	"context"
	"errors"
	"log/slog"
	"quicknotes/internal/models"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrDuplicateEmail = newRepositoryError(errors.New("duplicated email"))
var ErrEmailNotFound = newRepositoryError(errors.New("email not found"))
var ErrInvalidTokenOrUserAlreadyConfirmed = newRepositoryError(errors.New("invalid token or user already confirmed"))

type UserRepository interface {
	Create(ctx context.Context, email, password, token string) (*models.User, string, error)
	ConfirmUserByToken(ctx context.Context, token string) error
	CreateResetPasswordToken(ctx context.Context, email, hashToken string) (string, error)
	GetUserConfirmationByToken(ctx context.Context, token string) (*models.UserConfirmationToken, error)
	UpdatePasswordByToken(ctx context.Context, newPassword, token string) (string, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository função de estanciação
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) GetUserConfirmationByToken(ctx context.Context, token string) (*models.UserConfirmationToken, error) {
	var userToken models.UserConfirmationToken
	query := `select id, user_id, token, confirmed, created_at, updated_at
	from users_confirmation_tokens
	where token = $1`

	row := ur.db.QueryRow(ctx, query, token)
	if err := row.Scan(&userToken.ID, &userToken.UserID,
		&userToken.Token, &userToken.Confirmed, &userToken.CreatedAt,
		&userToken.UpdatedAt); err != nil {
		return nil, newRepositoryError(err)
	}

	return &userToken, nil

}

func (ur *userRepository) CreateResetPasswordToken(ctx context.Context, email, hashToken string) (string, error) {
	user, err := ur.FindByEmail(ctx, email)
	if err != nil || !user.Active.Bool {
		return "", ErrEmailNotFound
	}
	tx, err := ur.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", newRepositoryError(err)
	}
	userToken, err := ur.createConfirmationToken(tx, ctx, user, hashToken)
	if err != nil {
		return "", ErrEmailNotFound
	}
	if err = tx.Commit(ctx); err != nil {
		return "", newRepositoryError(err)
	}
	return userToken.Token.String, nil
}

func (ur *userRepository) UpdatePasswordByToken(ctx context.Context, newPassword, token string) (string, error) {
	//atualizar o token
	query := `SELECT u.id u_id, u.email, t.id t_id FROM  users u INNER JOIN users_confirmation_tokens t
	ON u.id = t.user_id
	WHERE t.confirmed = false 
	AND t.token = $1`

	row := ur.db.QueryRow(ctx, query, token)
	var userID, tokenID pgtype.Numeric
	var email pgtype.Text
	err := row.Scan(&userID, &email, &tokenID)
	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Error("token não encontrado ou já confirmado: " + err.Error())
			return "", ErrInvalidTokenOrUserAlreadyConfirmed
		}
		slog.Error(err.Error())
		return "", newRepositoryError(err)
	}

	//atualiza o confirmation token
	fail := func(err error) (string, error) {
		slog.Error(err.Error())
		return "", newRepositoryError(err)
	}
	tx, err := ur.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback(ctx)

	query = "UPDATE users_confirmation_tokens SET confirmed = true, updated_at = now() WHERE id = $1"
	_, err = tx.Exec(ctx, query, tokenID)
	if err != nil {
		slog.Error(err.Error())
		return fail(err)
	}

	//atualiza a senha do usuário
	query = "UPDATE users SET password = $1, updated_at = now() WHERE id = $2"
	_, err = tx.Exec(ctx, query, newPassword, userID)
	if err != nil {
		slog.Error(err.Error())
		return fail(err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fail(err)
	}

	return email.String, nil
}

func (ur *userRepository) Create(ctx context.Context, email, password, hashToken string) (*models.User, string, error) {
	var user models.User
	user.Email = pgtype.Text{String: email, Valid: true}
	user.Password = pgtype.Text{String: password, Valid: true}

	fail := func(err error) (*models.User, string, error) {
		slog.Error(err.Error())
		return &user, "", newRepositoryError(err)
	}

	tx, err := ur.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO users(email, password)
			  VALUES($1, $2)
			  RETURNING id, created_at;`
	row := tx.QueryRow(ctx, query, user.Email, user.Password)
	if err := row.Scan(&user.ID, &user.CreatedAt); err != nil {
		if strings.Contains(err.Error(), "violates unique constraint") {
			return fail(ErrDuplicateEmail)
		}
		return fail(err)
	}
	userToken, err := ur.createConfirmationToken(tx, ctx, &user, hashToken)
	if err != nil {
		return fail(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return fail(err)
	}
	return &user, userToken.Token.String, nil
}

func (ur *userRepository) createConfirmationToken(tx pgx.Tx, ctx context.Context, user *models.User, token string) (*models.UserConfirmationToken, error) {
	var userToken models.UserConfirmationToken
	userToken.Token = pgtype.Text{String: token, Valid: true}
	userToken.UserID = user.ID
	query := `INSERT INTO users_confirmation_tokens (user_id, token)
	VALUES($1, $2)
	RETURNING id, created_at`

	row := tx.QueryRow(ctx, query, userToken.UserID, userToken.Token)
	if err := row.Scan(&userToken.ID, &userToken.CreatedAt); err != nil {
		tx.Rollback(ctx)
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

	// escopo de transação

	fail := func(err error) error {
		slog.Error(err.Error())
		return newRepositoryError(err)
	}
	tx, err := ur.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback(ctx)

	// tornar o usuário active
	queryUpdateUser := "UPDATE users SET active = true, updated_at = now() WHERE id = $1"
	_, err = tx.Exec(ctx, queryUpdateUser, userID)
	if err != nil {
		return fail(err)
	}

	queryUpdateToken := `UPDATE users_confirmation_tokens 
	SET confirmed = true, updated_at = now()
	WHERE id = $1`
	_, err = tx.Exec(ctx, queryUpdateToken, tokenID)
	if err != nil {
		return newRepositoryError(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return fail(err)
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
