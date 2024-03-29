package router

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetRouting() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	e := echo.New()
	e.Use(middleware.Logger())

  var allowOrigins []string
  if os.Getenv("FRONTEND_ORIGIN") != "" {
    allowOrigins = append(allowOrigins, os.Getenv("FRONTEND_ORIGIN"))
  } 
  allowOrigins = append(allowOrigins, "http://localhost:3000")

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins:     allowOrigins,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "X-Traq-User"},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	api := e.Group("/api")
	{
		apiUsers := api.Group("/users")
		{
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
