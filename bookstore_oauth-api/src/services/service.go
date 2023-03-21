package services

import (
	"net/http"
	"strings"

	"github.com/AJackTi/bookstore_oauth-api/src/domain/access_token"
	"github.com/AJackTi/bookstore_oauth-api/src/repository/db"
	"github.com/AJackTi/bookstore_oauth-api/src/repository/rest"
	"github.com/AJackTi/bookstore_oauth-api/src/utils/errors"
)

type Service interface {
	GetByID(string) (*access_token.AccessToken, *errors.RestErr)
	Create(access_token.AccessTokenRequest) (*access_token.AccessToken, *errors.RestErr)
	UpdateExpirationTime(*access_token.AccessToken) *errors.RestErr
}

type service struct {
	repository          db.DbRepository
	restUsersRepository rest.RestUsersRepository
}

func NewService(repository db.DbRepository, restUsersRepository rest.RestUsersRepository) Service {
	return &service{
		repository,
		restUsersRepository,
	}
}

func (s *service) GetByID(accessTokenID string) (*access_token.AccessToken, *errors.RestErr) {
	accessTokenID = strings.TrimSpace(accessTokenID)
	if len(accessTokenID) == 0 {
		return nil, errors.New(http.StatusBadRequest, "invalid access token id")
	}
	accessToken, err := s.repository.GetByID(accessTokenID)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

func (s *service) Create(request access_token.AccessTokenRequest) (*access_token.AccessToken, *errors.RestErr) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	//Authenicate User against the user API with email / password
	user, err := s.restUsersRepository.LoginUser(request.UserName, request.Password)
	if err != nil {
		return nil, err
	}
	//Genreate a new Access token
	at := access_token.GetNewAccessToken(user.ID)
	at.Generate()

	//save the new accesstoken in the cassandra db
	if err := s.repository.Create(&at); err != nil {
		return nil, err
	}

	return &at, nil
}

func (s *service) UpdateExpirationTime(at *access_token.AccessToken) *errors.RestErr {
	if err := at.Validate(); err != nil {
		return err
	}

	return s.repository.UpdateExpirationTime(at)
}
