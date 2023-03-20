package app

import (
	"github.com/AJackTi/bookstore_users-api/logger"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	route()
	logger.Info("Start the application...")
	router.Run(":8080")
}
