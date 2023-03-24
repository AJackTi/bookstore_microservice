package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/AJackTi/bookstore_items-api/domain/items"
	"github.com/AJackTi/bookstore_items-api/services"
	"github.com/AJackTi/bookstore_items-api/utils/http_utils"
	"github.com/AJackTi/bookstore_oauth-go/oauth"
	"github.com/AJackTi/bookstore_utils-go/rest_errors"
)

var (
	ItemsController itemsControllerInterface = &itemsController{}
)

type itemsControllerInterface interface {
	Create(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
}

type itemsController struct{}

func (c *itemsController) Create(w http.ResponseWriter, r *http.Request) {
	if err := oauth.AuthenticateRequest(r); err != nil {
		http_utils.ResponseError(w, err)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.ResponseError(w, rest_errors.New(http.StatusBadRequest, errors.New("invalid request body")))
		return
	}
	defer r.Body.Close()

	var itemRequest items.Item
	if err := json.Unmarshal(requestBody, &itemRequest); err != nil {
		http_utils.ResponseError(w, rest_errors.New(http.StatusBadRequest, errors.New("invalid item json body")))
		return
	}

	itemRequest.Seller = oauth.GetClientId(r)
	result, createErr := services.ItemsService.Create(itemRequest)
	if createErr != nil {
		http_utils.ResponseError(w, createErr)
		return
	}

	http_utils.ResponseJson(w, http.StatusCreated, result)
}

func (c *itemsController) Get(w http.ResponseWriter, r *http.Request) {
	return
}
