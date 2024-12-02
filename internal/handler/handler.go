package handler

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func GetToken(ctx echo.Context) string {
	key := "Bearer "
	reqToken := ctx.Request().Header.Get("Authorization")
	if !strings.HasPrefix(reqToken, key) {
		return ""
	}

	token := strings.TrimPrefix(reqToken, key)
	return token
}
