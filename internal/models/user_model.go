package models

import "github.com/jackc/pgx/v5/pgtype"

type User struct {
	ID        pgtype.Numeric
	Email     pgtype.Text
	Password  pgtype.Text
	Active    pgtype.Bool
	CreatedAt pgtype.Date
	UpdatedAt pgtype.Date
}

type UserConfirmationToken struct {
	ID        pgtype.Numeric
	UserID    pgtype.Numeric
	Token     pgtype.Text
	Confirmed pgtype.Bool
	CreatedAt pgtype.Date
	UpdatedAt pgtype.Date
}
