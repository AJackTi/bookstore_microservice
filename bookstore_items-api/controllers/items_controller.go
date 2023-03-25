package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/AJackTi/bookstore_items-api/domain/items"
	"github.com/AJackTi/bookstore_items-api/domain/queries"
	"github.com/AJackTi/bookstore_items-api/services"
	"github.com/AJackTi/bookstore_items-api/utils/http_utils"
	"github.com/AJackTi/bookstore_oauth-go/oauth"
	"github.com/AJackTi/bookstore_utils-go/rest_errors"
	"github.com/gorilla/mux"
)

var (
	ItemsController itemsControllerInterface = &itemsController{}
)

type itemsControllerInterface interface {
	Create(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
	Search(http.ResponseWriter, *http.Request)
}

type itemsController struct{}

func (c *itemsController) Create(w http.ResponseWriter, r *http.Request) {
	if err := oauth.AuthenticateRequest(r); err != nil {
		http_utils.ResponseError(w, err)
		return
	}
	sellerID := oauth.GetCallerId(r)
	if sellerID == 0 {
		http_utils.ResponseError(w, rest_errors.New(http.StatusUnauthorized, errors.New("unauthorized permission")))
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

	itemRequest.Seller = sellerID
	result, createErr := services.ItemsService.Create(&itemRequest)
	if createErr != nil {
		http_utils.ResponseError(w, createErr)
		return
	}

	http_utils.ResponseJson(w, http.StatusCreated, result)
}

func (c *itemsController) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemID := strings.TrimSpace(vars["id"])

	item, err := services.ItemsService.Get(itemID)
	if err != nil {
		http_utils.ResponseError(w, err)
		return
	}

	http_utils.ResponseJson(w, http.StatusOK, item)
}

func (c *itemsController) Search(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.ResponseError(w, rest_errors.New(http.StatusBadRequest, fmt.Errorf("invalid json body")))
		return
	}
	defer r.Body.Close()

	var query queries.EsQuery
	if err := json.Unmarshal(bytes, &query); err != nil {
		http_utils.ResponseError(w, rest_errors.New(http.StatusBadRequest, fmt.Errorf("invalid json body")))
		return
	}

	items, restErr := services.ItemsService.Search(query)
	if restErr != nil {
		http_utils.ResponseError(w, restErr)
		return
	}

	http_utils.ResponseJson(w, http.StatusOK, items)
}
