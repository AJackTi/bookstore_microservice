package http

import (
	"errors"
	"net/http"

	"github.com/AJackTi/bookstore_oauth-api/src/domain/access_token"
	"github.com/AJackTi/bookstore_oauth-api/src/services"
	"github.com/AJackTi/bookstore_utils-go/rest_errors"
	"github.com/gin-gonic/gin"
)

type AccessTokenHandler interface {
	GetByID(*gin.Context)
	Create(*gin.Context)
}

type accessTokenHandler struct {
	service services.Service
}

func NewHandler(service services.Service) AccessTokenHandler {
	return &accessTokenHandler{
		service: service,
	}
}

func (handler *accessTokenHandler) GetByID(c *gin.Context) {
	accessTokenID := c.Param("access_token_id")
	accessToken, err := handler.service.GetByID(accessTokenID)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, accessToken)
}

func (handler *accessTokenHandler) Create(c *gin.Context) {
	var at access_token.AccessTokenRequest
	if err := c.ShouldBindJSON(&at); err != nil {
		c.JSON(http.StatusBadRequest, rest_errors.New(http.StatusBadRequest, errors.New("invalid json body")))
		return
	}

	if _, err := handler.service.Create(at); err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusCreated, at)
}
