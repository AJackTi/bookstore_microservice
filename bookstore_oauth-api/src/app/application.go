package app

import (
	"github.com/AJackTi/bookstore_oauth-api/src/http"
	"github.com/AJackTi/bookstore_oauth-api/src/repository/db"
	"github.com/AJackTi/bookstore_oauth-api/src/repository/rest"
	"github.com/AJackTi/bookstore_oauth-api/src/services"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	dbRepository := db.NewRepository()
	restRepository := rest.NewRepository()
	atService := services.NewService(dbRepository, restRepository)
	atHandler := http.NewHandler(atService)

	router.GET("/oauth/access_token/:access_token_id", atHandler.GetByID)
	router.POST("/oauth/access_token", atHandler.Create)

	router.Run(":8081")
}
