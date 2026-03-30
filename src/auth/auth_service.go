package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/lupppig/movy/internal/openapi"
	"github.com/lupppig/movy/internal/utils"
	openapi_types "github.com/oapi-codegen/runtime/types"
)



type UserError struct {
	Code string
	Message string
}


type UserReq struct {
        Name     string `valid:"required,length(2|50)"`
        Email    openapi_types.Email `valid:"required,email"`
        Password string `valid:"required,length(8|255)"`
}


func (u UserError) Error() string {
	return u.Message
}

func (a *AuthDep) CreateUserService(req openapi.SignupRequest) (int, *openapi.UserResponse, *UserError) {
    v := UserReq{
        Name:     req.Name,
        Email:    req.Email,
        Password: req.Password,
    }

    _, err := govalidator.ValidateStruct(v)
    if err != nil {
		a.Logger.Error().Err(err).Msg("validation failed for user")
        return http.StatusBadRequest, nil, &UserError{
            Code:    openapi.CodeInvalidInput,
            Message: err.Error(),
        }
    }

	v.Password, err = utils.HashPassword(v.Password)
	if err != nil {
		a.Logger.Error().Err(err).Msg("failed to hash user password")
		return http.StatusInternalServerError,  nil, &UserError{
			Code: openapi.CodeInternalError, 
			Message: "an internal server error occured",
		}
	}

	// add data to table
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	emailExist, err := a.CheckEmailExists(ctx, string(v.Email))
	if err != nil {
		a.Logger.Error().Err(err).Str("email", string(v.Email)).Msg("failed to check if user with this email exists")
		return http.StatusInternalServerError, nil, &UserError{
			Code: openapi.CodeInternalError,
			Message: "an internal server error occured",
		}
	}

	if emailExist {
		a.Logger.Warn().Str("email", string(v.Email)).Msg("user with this email already exists")
		return http.StatusConflict, nil, &UserError{
			Code: openapi.CodeEmailTaken,
			Message: "user with this email already exists",
		}
	}
	
	
	id, err := a.CreateUser(ctx, v)
	if err != nil {
		a.Logger.Error().Str("id", id).Str("email", string(v.Email)).Msg("user couldn't be created")
		return http.StatusInternalServerError, nil, &UserError{
			Code: openapi.CodeInternalError,
			Message: "an internal server error occured",
		}
	}

	uId, _ := uuid.Parse(id)
    return  http.StatusCreated, &openapi.UserResponse{
		Email: v.Email,
		Id: openapi_types.UUID(uId),
		Name: v.Name,
	}, nil
}
