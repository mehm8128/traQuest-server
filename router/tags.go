package router

import (
	"net/http"
	"traQuest/model"

	"github.com/labstack/echo/v4"
)

type TagRequest struct {
	Name string `json:"name"`
}

func postTag(c echo.Context) error {
	ctx := c.Request().Context()
	var tag TagRequest
	err := c.Bind(&tag)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	res, err := model.CreateTag(ctx, tag.Name)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return echo.NewHTTPError(http.StatusCreated, res)
}
