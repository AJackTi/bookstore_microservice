package items

import (
	"errors"
	"net/http"

	"github.com/AJackTi/bookstore_items-api/clients/elasticsearch"
	"github.com/AJackTi/bookstore_utils-go/rest_errors"
)

const (
	indexItems = "items"
)

func (i *Item) Save() *rest_errors.RestErr {
	result, err := elasticsearch.Client.Index(indexItems, i)
	if err != nil {
		return rest_errors.New(http.StatusInternalServerError, errors.New("error when trying to save item"))
	}
	i.ID = result.Id
	return nil
}
