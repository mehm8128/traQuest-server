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
type UserIDRequest struct {
	UserID uuid.UUID `json:"userId"`
}

const adminUserID = "c714a848-2886-4c10-a313-de9bc61cb2bb"

func getQuests(c echo.Context) error {
	ctx := c.Request().Context()

	userId, err := GetMeTraq(c)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	quests, err := model.GetQuests(ctx, userId)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return echo.NewHTTPError(http.StatusOK, &quests)
}

func getUnapprovedQuests(c echo.Context) error {
	ctx := c.Request().Context()
	userId, err := GetMeTraq(c)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if userId != adminUserID {
		return echo.NewHTTPError(http.StatusForbidden, "you are not admin")
	}
	quests, err := model.GetUnapprovedQuests(ctx)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return echo.NewHTTPError(http.StatusOK, &quests)
}

func getQuest(c echo.Context) error {
	ctx := c.Request().Context()
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	userId, err := GetMeTraq(c)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	quest, err := model.GetQuest(ctx, ID, userId)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return echo.NewHTTPError(http.StatusOK, &quest)
}

func completeQuest(c echo.Context) error {
	ctx := c.Request().Context()
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	var userID UserIDRequest
	err = c.Bind(&userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	err = model.CompleteQuest(ctx, ID, userID.UserID)
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
	return echo.NewHTTPError(http.StatusCreated, &res)
}

func rejectQuest(c echo.Context) error {
	ctx := c.Request().Context()
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	var userID UserIDRequest
	err = c.Bind(&userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if userID.UserID.String() != adminUserID {
		return echo.NewHTTPError(http.StatusForbidden, "you are not admin")
	}
	err = model.RejectQuest(ctx, ID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return echo.NewHTTPError(http.StatusOK, "ok")
}

func approveQuest(c echo.Context) error {
	ctx := c.Request().Context()
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	var userID UserIDRequest
	err = c.Bind(&userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if userID.UserID.String() != adminUserID {
		return echo.NewHTTPError(http.StatusForbidden, "you are not admin")
	}
	res, err := model.ApproveQuest(ctx, ID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return echo.NewHTTPError(http.StatusCreated, &res)
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
	return echo.NewHTTPError(http.StatusOK, &res)
}
