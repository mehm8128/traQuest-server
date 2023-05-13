package router

import (
	"net/http"
	"traQuest/model"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type QuestRequest struct {
	Title       string      `json:"title" db:"title"`
	Description string      `json:"description" db:"description"`
	Level       int         `json:"level" db:"level"`
	Tags        []uuid.UUID `json:"tags"`
}

func getQuests(c echo.Context) error {
	ctx := c.Request().Context()
	quests, err := model.GetQuests(ctx, uuid.Nil) //todo: userID
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return echo.NewHTTPError(http.StatusOK, quests)
}

func getUnapprovedQuests(c echo.Context) error {
	// todo: admin check
	ctx := c.Request().Context()
	quests, err := model.GetUnapprovedQuests(ctx)
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

func completeQuest(c echo.Context) error {
	ctx := c.Request().Context()
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	err = model.CompleteQuest(ctx, ID, uuid.Nil) //todo: userID
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return echo.NewHTTPError(http.StatusOK, "ok")
}

func postQuest(c echo.Context) error {
	ctx := c.Request().Context()
	var quest QuestRequest
	err := c.Bind(&quest)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	res, err := model.CreateQuest(ctx, quest.Title, quest.Description, quest.Level, quest.Tags)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return echo.NewHTTPError(http.StatusCreated, res)
}

func approveQuest(c echo.Context) error {
	ctx := c.Request().Context()
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	res, err := model.ApproveQuest(ctx, ID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return echo.NewHTTPError(http.StatusCreated, res)
}

func putQuest(c echo.Context) error {
	ctx := c.Request().Context()
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	var quest QuestRequest
	err = c.Bind(&quest)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	res, err := model.UpdateQuest(ctx, ID, quest.Title, quest.Description, quest.Level, quest.Tags)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return echo.NewHTTPError(http.StatusOK, res)
}
