package db

import (
	"github.com/AJackTi/bookstore_oauth-api/src/domain/access_token"
	"github.com/AJackTi/bookstore_oauth-api/src/utils/errors"
)

func NewRepository() DbRepository {
	return &dbRepository{}
}

type DbRepository interface {
	GetByID(string) (*access_token.AccessToken, *errors.RestErr)
}

type dbRepository struct {
}

func (r *dbRepository) GetByID(id string) (*access_token.AccessToken, *errors.RestErr) {
	return nil, errors.NewInternalServerError("database connection not implemented yet")
}
