package items

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/AJackTi/bookstore_items-api/clients/elasticsearch"
	"github.com/AJackTi/bookstore_items-api/domain/queries"
	"github.com/AJackTi/bookstore_utils-go/rest_errors"
)

const (
	indexItems = "items"
	typeItem   = "_doc"
)

func (i *Item) Save() *rest_errors.RestErr {
	result, err := elasticsearch.Client.Index(indexItems, typeItem, i)
	if err != nil {
		return rest_errors.New(http.StatusInternalServerError, errors.New("error when trying to save item"))
	}
	i.ID = result.Id
	return nil
}

func (i *Item) Get() *rest_errors.RestErr {
	itemID := i.ID
	result, err := elasticsearch.Client.Get(indexItems, typeItem, i.ID)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return rest_errors.New(http.StatusNotFound, fmt.Errorf("no item found with id %s", i.ID))
		}
		return rest_errors.New(http.StatusInternalServerError, fmt.Errorf("error when trying to get id %s", i.ID))
	}

	bytes, err := result.Source.MarshalJSON()
	if err != nil {
		return rest_errors.New(http.StatusInternalServerError, fmt.Errorf("error when trying to marshal response json"))
	}
	if err := json.Unmarshal(bytes, &i); err != nil {
		return rest_errors.New(http.StatusInternalServerError, fmt.Errorf("error when trying to parse database response"))
	}
	fmt.Println(string(bytes))
	i.ID = itemID
	return nil
}

func (i *Item) Search(query queries.EsQuery) ([]Item, *rest_errors.RestErr) {
	result, err := elasticsearch.Client.Search(indexItems, query.Build())
	if err != nil {
		return nil, rest_errors.New(http.StatusInternalServerError, fmt.Errorf("error when trying to search documents"))
	}

	log.Println(result)

	items := make([]Item, result.TotalHits())
	for index, hit := range result.Hits.Hits {
		bytes, _ := hit.Source.MarshalJSON()
		var item Item
		if err := json.Unmarshal(bytes, &item); err != nil {
			return nil, rest_errors.New(http.StatusInternalServerError, fmt.Errorf("error when trying to parse response"))
		}
		item.ID = hit.Id
		items[index] = item
	}

	if len(items) == 0 {
		return nil, rest_errors.New(http.StatusNotFound, fmt.Errorf("no items found matching given criteria"))
	}

	return items, nil
}
