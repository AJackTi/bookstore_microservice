package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/AJackTi/bookstore_oauth-api/src/domain/users"
	"github.com/AJackTi/bookstore_utils-go/rest_errors"
	"github.com/mercadolibre/golang-restclient/rest"
)

var (
	usersRestClient = rest.RequestBuilder{
		BaseURL: "http://localhost:8080",
		Timeout: 100 * time.Millisecond,
	}
)

type RestUsersRepository interface {
	LoginUser(string, string) (*users.User, *rest_errors.RestErr)
}

type usersRepository struct {
}

func NewRepository() RestUsersRepository {
	return &usersRepository{}
}

func (r *usersRepository) LoginUser(email string, password string) (*users.User, *rest_errors.RestErr) {
	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}
	response := usersRestClient.Post("/users/login", request)
	if response == nil || response.Response == nil {
		return nil, rest_errors.New(http.StatusInternalServerError, errors.New("invalid restclient response when trying to login user"))
	}

	if response.StatusCode > 299 {
		var restErr rest_errors.RestErr
		err := json.Unmarshal(response.Bytes(), &restErr)
		if err != nil {
			return nil, rest_errors.New(http.StatusInternalServerError, errors.New("invalid error interface when trying to login user"))
		}

		return nil, &restErr
	}

	var user users.User
	if err := json.Unmarshal(response.Bytes(), &user); err != nil {
		return nil, rest_errors.New(http.StatusInternalServerError, errors.New("error when trying to unmarshal users response"))
	}

	return &user, nil
}
