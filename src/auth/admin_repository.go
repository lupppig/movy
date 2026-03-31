package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (a *AuthDep) UpdateUserRole(ctx context.Context, userID, role string) error {
	db := a.DB
	query := `
		UPDATE users
		SET role = $1
		WHERE id = $2
	`

	result, err := db.ExecContext(ctx, query, role, userID)
	if err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found with id: %s", userID)
	}

	return nil
}

func (a *AuthDep) FindUserByID(ctx context.Context, userID string) (*UserRow, error) {
	db := a.DB
	query := `
		SELECT id, name, email, password, role
		FROM users
		WHERE id = $1
	`

	var user UserRow
	err := db.QueryRowContext(ctx, query, userID).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	if err := uuid.Validate(user.ID); err != nil {
		return nil, fmt.Errorf("returned improper uuid format: %w", err)
	}

	return &user, nil
}

func (a *AuthDep) FindAllUsers(ctx context.Context) ([]UserRow, error) {
	db := a.DB
	query := `
		SELECT id, name, email, password, role
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []UserRow
	for rows.Next() {
		var user UserRow
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil
}
