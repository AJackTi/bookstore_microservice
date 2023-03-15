package users

import (
	"net/http"
	"strconv"

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
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("user id should be a number"))
		return
	}

	user, restErr := services.GetUser(userID)
	if restErr != nil {
		c.JSON(restErr.Status, restErr)
		return
	}

	c.JSON(http.StatusOK, *user)
}

func FindUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "implement me!")
}
