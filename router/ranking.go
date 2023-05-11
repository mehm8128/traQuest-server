package router

import (
	"traQuest/model"

	"github.com/labstack/echo/v4"
)

func getRanking(c echo.Context) error {
	ranking, err := model.GetRanking()
	if err != nil {
		return err
	}
	return c.JSON(200, ranking)
}
