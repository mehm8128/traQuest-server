package router

import (
	"net/http"
	"traQuest/model"

	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func getSignin(c echo.Context) error {
	return nil
}

func getMe(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "something wrong in getting session")
	}
	userID := sess.Values["userID"]
	if userID == nil {
		return echo.NewHTTPError(http.StatusForbidden, "not logged in")
	}
	ctx := c.Request().Context()
	ID, err := uuid.Parse(userID.(string))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	user, err := model.GetUser(ctx, ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return echo.NewHTTPError(http.StatusOK, &user)
}
