package router

import (
	"net/http"
	"os"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/srinathgs/mysqlstore"
)

func SetRouting(store *mysqlstore.MySQLStore) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	api := e.Group("/api")
	{
		apiUsers := api.Group("/users")
		{
			apiUsers.GET("/signin", getSignin)
			apiUsers.GET("", getUsers)
			apiUsers.GET("/me", getMe)
		}
		apiQuests := api.Group("/quests")
		{
			apiQuests.GET("", getQuests)
			apiQuests.POST("", postQuest)
			apiQuests.GET("/:id", getQuest)
			apiQuests.PUT("/:id", putQuest)
			apiQuests.GET("/user/:id", getQuestsByUser)
		}
		apiRanking := api.Group("/ranking")
		{
			apiRanking.GET("", getRanking)
		}
	}

	e.Logger.Fatal(e.Start(":" + port))
}
