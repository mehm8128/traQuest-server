package router

import (
	"net/http"
	"traQuest/model"

	"github.com/labstack/echo"
)

func getQuests(c echo.Context) error {
	ctx := c.Request().Context()
	quests, err := model.GetQuests(ctx)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return echo.NewHTTPError(http.StatusOK, quests)
}
