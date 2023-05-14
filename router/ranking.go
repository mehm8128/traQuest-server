package router

import (
	"net/http"
	"traQuest/model"

	"github.com/labstack/echo/v4"
)

func getRanking(c echo.Context) error {
	ctx := c.Request().Context()
	ranking, err := model.GetRanking(ctx)
	if err != nil {
		return err
	}
	return echo.NewHTTPError(http.StatusOK, &ranking)
}
