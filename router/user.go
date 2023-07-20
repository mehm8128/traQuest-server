package router

import (
	"fmt"
	"net/http"
	"traQuest/model"

	"github.com/labstack/echo/v4"
)

type User struct {
	ID string `json:"id"`
}

const (
	SHOWCASE_USER_KEY = "X-Forwarded-User"
)

func getMe(c echo.Context) error {
	userId, err := GetMeTraq(c)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	ctx := c.Request().Context()
	isUserExist, err := model.IsUserExist(ctx, userId)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if !isUserExist {
		_, err := model.CreateUser(ctx, userId)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	me, err := model.GetUser(ctx, userId)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return echo.NewHTTPError(http.StatusOK, &me)
}

func GetMeTraq(c echo.Context) (string, error) {
	userId := c.Request().Header.Get(SHOWCASE_USER_KEY)
	fmt.Printf("%s %#v", userId, c.Request().Header)
	if userId == "" {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}
	return userId, nil
}
