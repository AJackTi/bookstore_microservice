package users

import (
	"net/http"
	"strconv"

	"github.com/AJackTi/bookstore_oauth-go/oauth"
	"github.com/AJackTi/bookstore_users-api/domain/users"
	"github.com/AJackTi/bookstore_users-api/services"
	"github.com/AJackTi/bookstore_users-api/utils/errors"
	"github.com/gin-gonic/gin"
)

func getUserID(userIDParam string) (int64, *errors.RestErr) {
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		return 0, errors.NewBadRequestError("user id should be a number")
	}

	return userID, nil
}

func Create(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid json body"))
		return
	}

	result, restErr := services.UsersService.CreateUser(&user)
	if restErr != nil {
		c.JSON(restErr.Status, *restErr)
		return
	}

	c.JSON(http.StatusCreated, *result)
}

func Get(c *gin.Context) {
	if err := oauth.AuthenticateRequest(c.Request); err != nil {
		c.JSON(err.Status, err)
		return
	}

	userID, err := getUserID(c.Param("user_id"))
	if err != nil {
		c.JSON(err.Status, errors.NewBadRequestError("user id should be a number"))
		return
	}

	user, restErr := services.UsersService.GetUser(userID)
	if restErr != nil {
		c.JSON(restErr.Status, restErr)
		return
	}

	if oauth.GetCallerId(c.Request) == user.ID {
		c.JSON(http.StatusOK, user.Marshal(false))
		return
	}

	c.JSON(http.StatusOK, (*user).Marshal(c.GetHeader("X-Public") == "true"))
}

func Update(c *gin.Context) {
	userID, err := getUserID(c.Param("user_id"))
	if err != nil {
		c.JSON(err.Status, errors.NewBadRequestError("user id should be a number"))
		return
	}

	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid json body"))
		return
	}

	user.ID = userID

	isPartial := c.Request.Method == http.MethodPatch

	result, restErr := services.UsersService.UpdateUser(isPartial, &user)
	if restErr != nil {
		c.JSON(restErr.Status, restErr)
		return
	}

	c.JSON(http.StatusOK, (*result).Marshal(c.GetHeader("X-Public") == "true"))
}

func Delete(c *gin.Context) {
	userID, err := getUserID(c.Param("user_id"))
	if err != nil {
		c.JSON(err.Status, errors.NewBadRequestError("user id should be a number"))
		return
	}

	if err := services.UsersService.DeleteUser(userID); err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func Search(c *gin.Context) {
	status := c.Query("status")

	users, err := services.UsersService.Search(status)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, users.Marshal(c.GetHeader("X-Public") == "true"))
}

func Login(c *gin.Context) {
	var request users.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid json body"))
		return
	}

	user, err := services.UsersService.LoginUser(&request)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, (*user).Marshal(c.GetHeader("X-Public") == "true"))
}
