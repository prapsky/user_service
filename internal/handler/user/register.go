package handler_user

import (
	goerrors "errors"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/prapsky/user_service/common/errors"
	"github.com/prapsky/user_service/common/response"
	service_user_register "github.com/prapsky/user_service/service/user/register"
)

type RegisterUserHandler struct {
	service *service_user_register.RegisterUserService
}

func NewRegisterUserHandler(service *service_user_register.RegisterUserService) *RegisterUserHandler {
	return &RegisterUserHandler{
		service: service,
	}
}

type RegisterUserRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

func (r RegisterUserRequest) registerInput() service_user_register.RegisterUserInput {
	return service_user_register.RegisterUserInput{
		Name:        r.Name,
		PhoneNumber: r.PhoneNumber,
		Username:    r.Username,
		Password:    r.Password,
	}
}

func (h *RegisterUserHandler) Register(ctx echo.Context) error {
	request := new(RegisterUserRequest)
	if rerr := ctx.Bind(request); rerr != nil {
		return ctx.JSON(http.StatusBadRequest, response.NewError(errors.ErrInvalidRequest))
	}

	registerInput := request.registerInput()
	token, err := h.service.Do(ctx.Request().Context(), registerInput)
	if err != nil {
		res := response.NewError(err)
		if goerrors.Is(err, errors.ErrUsernameAlreadyExists) {
			return ctx.JSON(http.StatusConflict, res)
		}

		res = response.NewError(errors.ErrInternalServerError)
		return ctx.JSON(http.StatusInternalServerError, res)
	}

	response := SuccessResponse(token)
	return ctx.JSON(http.StatusCreated, response)
}

type RegisterResponse struct {
	Token string `json:"token"`
}

func SuccessResponse(token string) RegisterResponse {
	return RegisterResponse{
		Token: token,
	}
}
