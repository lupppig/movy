package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lupppig/movy/internal/openapi"
	"github.com/lupppig/movy/internal/role"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type UserResponse struct {
	ID    openapi_types.UUID `json:"id"`
	Name  string             `json:"name"`
	Email string             `json:"email"`
	Role  string             `json:"role"`
}

func (a *AuthDep) PromoteUserService(userID string) (int, *UserResponse, *UserError) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	user, err := a.FindUserByID(ctx, userID)
	if err != nil {
		a.Logger.Error().Err(err).Str("user_id", userID).Msg("failed to find user")
		return http.StatusNotFound, nil, &UserError{
			Code:    openapi.CodeNotFound,
			Message: "user not found",
		}
	}

	if user.Role == role.Admin {
		return http.StatusConflict, nil, &UserError{
			Code:    openapi.CodeInvalidInput,
			Message: "user is already an admin",
		}
	}

	if err := a.UpdateUserRole(ctx, userID, role.Admin); err != nil {
		a.Logger.Error().Err(err).Str("user_id", userID).Msg("failed to update user role")
		return http.StatusInternalServerError, nil, &UserError{
			Code:    openapi.CodeInternalError,
			Message: "an internal server error occured",
		}
	}

	uID, _ := uuid.Parse(user.ID)
	return http.StatusOK, &UserResponse{
		ID:    openapi_types.UUID(uID),
		Name:  user.Name,
		Email: user.Email,
		Role:  role.Admin,
	}, nil
}

func (a *AuthDep) DemoteUserService(userID string) (int, *UserResponse, *UserError) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	user, err := a.FindUserByID(ctx, userID)
	if err != nil {
		a.Logger.Error().Err(err).Str("user_id", userID).Msg("failed to find user")
		return http.StatusNotFound, nil, &UserError{
			Code:    openapi.CodeNotFound,
			Message: "user not found",
		}
	}

	if user.Role == role.User {
		return http.StatusConflict, nil, &UserError{
			Code:    openapi.CodeInvalidInput,
			Message: "user is already a regular user",
		}
	}

	if err := a.UpdateUserRole(ctx, userID, role.User); err != nil {
		a.Logger.Error().Err(err).Str("user_id", userID).Msg("failed to update user role")
		return http.StatusInternalServerError, nil, &UserError{
			Code:    openapi.CodeInternalError,
			Message: "an internal server error occured",
		}
	}

	uID, _ := uuid.Parse(user.ID)
	return http.StatusOK, &UserResponse{
		ID:    openapi_types.UUID(uID),
		Name:  user.Name,
		Email: user.Email,
		Role:  role.User,
	}, nil
}

func (a *AuthDep) GetUsersService() (int, []UserResponse, *UserError) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	users, err := a.FindAllUsers(ctx)
	if err != nil {
		a.Logger.Error().Err(err).Msg("failed to fetch users")
		return http.StatusInternalServerError, nil, &UserError{
			Code:    openapi.CodeInternalError,
			Message: "an internal server error occured",
		}
	}

	var responses []UserResponse
	for _, user := range users {
		uID, _ := uuid.Parse(user.ID)
		responses = append(responses, UserResponse{
			ID:    openapi_types.UUID(uID),
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		})
	}

	return http.StatusOK, responses, nil
}

func (a *AuthDep) GetUserService(userID string) (int, *UserResponse, *UserError) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	user, err := a.FindUserByID(ctx, userID)
	if err != nil {
		a.Logger.Error().Err(err).Str("user_id", userID).Msg("failed to find user")
		return http.StatusNotFound, nil, &UserError{
			Code:    openapi.CodeNotFound,
			Message: "user not found",
		}
	}

	uID, _ := uuid.Parse(user.ID)
	return http.StatusOK, &UserResponse{
		ID:    openapi_types.UUID(uID),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}
