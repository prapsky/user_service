package handler_user

import (
	goerrors "errors"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/prapsky/user_service/common/errors"
	"github.com/prapsky/user_service/common/response"
	"github.com/prapsky/user_service/entity"
	"github.com/prapsky/user_service/internal/handler"
	service_user_detail "github.com/prapsky/user_service/service/user/detail"
)

type UserDetailHandler struct {
	service *service_user_detail.UserDetailService
}

func NewUserDetailHandler(service *service_user_detail.UserDetailService) *UserDetailHandler {
	return &UserDetailHandler{
		service: service,
	}
}

func (h *UserDetailHandler) Detail(ctx echo.Context) error {
	token := handler.GetToken(ctx)
	if token == "" {
		return ctx.JSON(http.StatusUnauthorized, response.NewError(errors.ErrUnauthorizedUser))
	}

	input := service_user_detail.UserDetailInput{
		Token: token,
	}
	result, err := h.service.Do(ctx.Request().Context(), input)
	if err != nil {
		res := response.NewError(err)
		if goerrors.Is(err, errors.ErrInvalidToken) {
			return ctx.JSON(http.StatusBadRequest, res)
		}
		if goerrors.Is(err, errors.ErrAccountNotRegistered) {
			return ctx.JSON(http.StatusBadRequest, res)
		}
		res = response.NewError(errors.ErrInternalServerError)
		return ctx.JSON(http.StatusInternalServerError, res)
	}

	response := SuccessUserDetailResponse(result)
	return ctx.JSON(http.StatusOK, response)
}

type DetailResponse struct {
	Data UserDetailResponse `json:"data"`
}

type UserDetailResponse struct {
	Name        string `json:"name"`
	Username    string `json:"username"`
	Phonenumber string `json:"phone_number"`
}

func SuccessUserDetailResponse(user *entity.User) DetailResponse {
	return DetailResponse{
		Data: UserDetailResponse{
			Name:        user.Name,
			Username:    user.Username,
			Phonenumber: user.PhoneNumber,
		},
	}
}
