package router

import (
	"net/http"
	"traQuest/model"

	"github.com/labstack/echo/v4"
)

type TagRequest struct {
	Name string `json:"name"`
}

func getTags(c echo.Context) error {
	ctx := c.Request().Context()
	tags, err := model.GetTags(ctx)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return echo.NewHTTPError(http.StatusOK, &tags)
}

func postTag(c echo.Context) error {
	ctx := c.Request().Context()
	var tags []string
	err := c.Bind(&tags)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	res, err := model.CreateTag(ctx, tags)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return echo.NewHTTPError(http.StatusCreated, &res)
}
