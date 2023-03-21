package db

import (
	"net/http"

	"github.com/AJackTi/bookstore_oauth-api/src/clients/cassandra"
	"github.com/AJackTi/bookstore_oauth-api/src/domain/access_token"
	"github.com/AJackTi/bookstore_oauth-api/src/utils/errors"
	"github.com/gocql/gocql"
)

const (
	queryGetAccessToken    = "SELECT access_token, user_id, client_id, expires FROM access_tokens WHERE access_token = ?"
	queryCreateAccessToken = "INSERT INTO access_tokens(access_token, user_id, client_id, expires) VALUES(?, ?, ?, ?)"
	queryUpdateExpires     = "UPDATE access_token SET expires = ? WHERE access_token = ?"
)

func NewRepository() DbRepository {
	return &dbRepository{}
}

type DbRepository interface {
	GetByID(string) (*access_token.AccessToken, *errors.RestErr)
	Create(*access_token.AccessToken) *errors.RestErr
	UpdateExpirationTime(*access_token.AccessToken) *errors.RestErr
}

type dbRepository struct {
}

func (r *dbRepository) GetByID(id string) (*access_token.AccessToken, *errors.RestErr) {
	var result access_token.AccessToken
	if err := cassandra.GetSession().Query(queryGetAccessToken, id).Scan(
		&result.AccessToken,
		&result.UserID,
		&result.ClientID,
		&result.Expires); err != nil {
		if err == gocql.ErrNotFound {
			return nil, errors.New(http.StatusNotFound, "no access token found with given id")
		}
		return nil, errors.New(http.StatusInternalServerError, err.Error())
	}

	return &result, nil
}

func (r *dbRepository) Create(at *access_token.AccessToken) *errors.RestErr {
	if err := cassandra.GetSession().
		Query(queryCreateAccessToken, at.AccessToken, at.UserID, at.ClientID, at.Expires).
		Exec(); err != nil {
		return errors.New(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (r *dbRepository) UpdateExpirationTime(at *access_token.AccessToken) *errors.RestErr {
	if err := cassandra.GetSession().
		Query(queryUpdateExpires, at.Expires, at.AccessToken).
		Exec(); err != nil {
		return errors.New(http.StatusInternalServerError, err.Error())
	}

	return nil
}
