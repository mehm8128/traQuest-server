package router

import (
	"fmt"
	"net/http"
	"os"
	"traQuest/model"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func authorizeHandler(c echo.Context) error {
	err := godotenv.Load()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://q.trap.jp/api/v3/oauth2/authorize?response_type=code&client_id=%s", os.Getenv("CLIENT_ID")), nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.Redirect(http.StatusFound, req.URL.String())
}

func callbackHandler(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "code is empty")
	}
	_, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://q.trap.jp/api/v3/oauth2/token?grant_type=authorization_code&client_id=%s&code=%s", os.Getenv("CLIENT_ID"), code), nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	// usernameを取得してDBに保存

	return echo.NewHTTPError(http.StatusOK, "ok")
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
