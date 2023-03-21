package http

import (
	"net/http"

	"github.com/AJackTi/bookstore_oauth-api/src/domain/access_token"
	"github.com/AJackTi/bookstore_oauth-api/src/utils/errors"
	"github.com/gin-gonic/gin"
)

type AccessTokenHandler interface {
	GetByID(*gin.Context)
	Create(*gin.Context)
}

type accessTokenHandler struct {
	service access_token.Service
}

func NewHandler(service access_token.Service) AccessTokenHandler {
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

	c.JSON(http.StatusNotImplemented, accessToken)
}

func (handler *accessTokenHandler) Create(c *gin.Context) {
	var at access_token.AccessToken
	if err := c.ShouldBindJSON(&at); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid json body"))
		return
	}

	if err := handler.service.Create(&at); err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusCreated, at)
}
