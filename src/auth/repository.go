package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (a *AuthDep) CreateUser(ctx context.Context, userReq UserReq) (string, error) {
	db := a.DB
	var Id string

	query := `
        INSERT INTO users (name, email, password, created_at)
        VALUES ($1, $2, $3, NOW()) RETURNING id
    `

	err := db.QueryRowContext(ctx, query, userReq.Name, string(userReq.Email), userReq.Password).Scan(&Id)

	if err != nil {
		return "", fmt.Errorf("failed to insert user and return ID: %w", err)
	}

	if err := uuid.Validate(Id); err != nil {
		return "", fmt.Errorf("returned improper uuid format %w", err)
	}
	return Id, err
}

func (a *AuthDep) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	db := a.DB
	var exists bool
	query := `
	SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
	`
	err := db.QueryRowContext(ctx, query, email).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

type UserRow struct {
	ID       string
	Name     string
	Email    string
	Password string
}

func (a *AuthDep) FindUserByEmail(ctx context.Context, email string) (*UserRow, error) {
	db := a.DB
	query := `
		SELECT id, name, email, password
		FROM users
		WHERE email = $1
	`

	var user UserRow
	err := db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	if err := uuid.Validate(user.ID); err != nil {
		return nil, fmt.Errorf("returned improper uuid format: %w", err)
	}

	return &user, nil
}
