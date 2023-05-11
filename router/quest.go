package router

import (
	"net/http"
	"traQuest/model"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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

func getQuest(c echo.Context) error {
	ctx := c.Request().Context()
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	quest, err := model.GetQuest(ctx, ID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return echo.NewHTTPError(http.StatusOK, quest)
}
