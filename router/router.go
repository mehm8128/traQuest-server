package router

import (
	"encoding/gob"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/srinathgs/mysqlstore"
	"golang.org/x/oauth2"
)

func SetRouting(store *mysqlstore.MySQLStore) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))

	clientOrigin := os.Getenv("FRONTEND_ORIGIN")
	if clientOrigin == "" {
		clientOrigin = "http://localhost:3000"
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{clientOrigin},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	defer store.Close()
	defer store.StopCleanup(store.Cleanup(time.Minute * 5))

	gob.Register("sessions")
	gob.Register(&oauth2.Token{})

	api := e.Group("/api")
	{
		apiUsers := api.Group("/users")
		{
			apiUsers.GET("/authorize", authorizeHandler)
			apiUsers.GET("/callback", callbackHandler)
			apiUsers.GET("/me", getMe)
		}
		apiQuests := api.Group("/quests")
		{
			apiQuests.GET("", getQuests)
			apiQuests.GET("/unapproved", getUnapprovedQuests)
			apiQuests.GET("/:id", getQuest)
			apiQuests.POST("/:id/complete", completeQuest)
			apiQuests.POST("", postQuest)
			apiQuests.POST("/:id/reject", rejectQuest)
			apiQuests.POST("/:id/approve", approveQuest)
			apiQuests.PUT("/:id", putQuest)
		}
		apiTags := api.Group("/tags")
		{
			apiTags.GET("", getTags)
			apiTags.POST("", postTag)
		}
		apiRanking := api.Group("/ranking")
		{
			apiRanking.GET("", getRanking)
		}
	}

	e.Logger.Fatal(e.Start(":" + port))
}
