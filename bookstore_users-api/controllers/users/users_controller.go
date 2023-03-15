package users

import (
	"net/http"

	"github.com/AJackTi/bookstore_users-api/domain/users"
	"github.com/AJackTi/bookstore_users-api/services"
	"github.com/AJackTi/bookstore_users-api/utils/errors"
	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid json body"))
		return
	}

	result, restErr := services.CreateUser(&user)
	if restErr != nil {
		c.JSON(restErr.Status, *restErr)
		return
	}

	c.JSON(http.StatusCreated, *result)
}

func GetUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "implement me!")
}

func FindUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "implement me!")
}
