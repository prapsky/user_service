package handler_user

import (
	goerrors "errors"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/prapsky/user_service/common/errors"
	"github.com/prapsky/user_service/common/response"
	service_user_login "github.com/prapsky/user_service/service/user/login"
)

type LoginUserHandler struct {
	service *service_user_login.LoginUserService
}

func NewLoginUserHandler(service *service_user_login.LoginUserService) *LoginUserHandler {
	return &LoginUserHandler{
		service: service,
	}
}

type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r LoginUserRequest) LoginInput() service_user_login.LoginUserInput {
	return service_user_login.LoginUserInput{
		Username: r.Username,
		Password: r.Password,
	}
}

func (h *LoginUserHandler) Login(ctx echo.Context) error {
	request := new(LoginUserRequest)
	if rerr := ctx.Bind(request); rerr != nil {
		return ctx.JSON(http.StatusBadRequest, response.NewError(errors.ErrInvalidRequest))
	}

	loginInput := request.LoginInput()
	token, err := h.service.Do(ctx.Request().Context(), loginInput)
	if err != nil {
		res := response.NewError(err)
		if goerrors.Is(err, errors.ErrAccountNotRegistered) {
			return ctx.JSON(http.StatusBadRequest, res)
		}

		if goerrors.Is(err, errors.ErrIncorrectPassword) {
			return ctx.JSON(http.StatusBadRequest, res)
		}

		res = response.NewError(errors.ErrInternalServerError)
		return ctx.JSON(http.StatusInternalServerError, res)
	}

	response := SuccessLoginResponse(token)
	return ctx.JSON(http.StatusOK, response)
}

type LoginResponse struct {
	Token string `json:"token"`
}

func SuccessLoginResponse(token string) LoginResponse {
	return LoginResponse{
		Token: token,
	}
}
