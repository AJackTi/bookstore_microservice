package services

import (
	"errors"
	"net/http"

	"github.com/AJackTi/bookstore_items-api/domain/items"
	"github.com/AJackTi/bookstore_utils-go/rest_errors"
)

var (
	ItemsService itemsServiceInterface = &itemService{}
)

type itemsServiceInterface interface {
	Create(*items.Item) (*items.Item, *rest_errors.RestErr)
	Get(string) (*items.Item, *rest_errors.RestErr)
}

type itemService struct{}

func (s *itemService) Create(item *items.Item) (*items.Item, *rest_errors.RestErr) {
	if err := item.Save(); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *itemService) Get(id string) (*items.Item, *rest_errors.RestErr) {
	return nil, rest_errors.New(http.StatusNotImplemented, errors.New("implement me"), "implement me")
}
