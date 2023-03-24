package access_token

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/AJackTi/bookstore_oauth-api/src/utils/cryptos"
	"github.com/AJackTi/bookstore_utils-go/rest_errors"
)

const (
	//number of hours a at is valid.
	accesstokenTTL             = 24
	grantTypePassword          = "password"
	grandTypeClientCredentials = "client_credentials"
)

// ACCESS TOKEN
type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserID      int64  `json:"user_id"`
	ClientID    int64  `json:"client_id"`
	Expires     int64  `json:"expires"`
}

func (AT *AccessToken) IsExpired() bool {
	return time.Unix(AT.Expires, 0).Before(time.Now().UTC())
}

func GetNewAccessToken(userID int64) AccessToken {
	return AccessToken{
		UserID:  userID,
		Expires: time.Now().UTC().Add(accesstokenTTL * time.Hour).Unix(),
	}
}

func (at *AccessToken) Validate() *rest_errors.RestErr {
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if at.AccessToken == "" {
		return rest_errors.New(http.StatusInternalServerError, errors.New("invalid access token"))
	}
	if at.UserID <= 0 {
		return rest_errors.New(http.StatusBadRequest, errors.New("invalid user"))
	}
	if at.ClientID <= 0 {
		return rest_errors.New(http.StatusBadRequest, errors.New("invalid client id"))
	}
	if at.Expires <= 0 {
		return rest_errors.New(http.StatusBadRequest, errors.New("invalid Expired time"))
	}

	return nil
}

// ACCESS TOKEN REQUEST
type AccessTokenRequest struct {
	GrantType string `json:"grant_type"`
	Scope     string `json:"scope"`

	//Used for password grant type
	UserName string `json:"username"`
	Password string `json:"password"`

	//Used for client credentials
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (atr *AccessTokenRequest) Validate() *rest_errors.RestErr {

	atr.GrantType = strings.TrimSpace(atr.GrantType)

	switch atr.GrantType {
	case grantTypePassword:
		atr.UserName = strings.TrimSpace(atr.UserName)
		atr.Password = strings.TrimSpace(atr.Password)
		if atr.UserName == "" {
			return rest_errors.New(http.StatusBadRequest, errors.New("invalid username or password"))
		}
		if atr.Password == "" {
			return rest_errors.New(http.StatusBadRequest, errors.New("invalid username or password"))
		}
		break

	case grandTypeClientCredentials:
		if atr.ClientId == "" {
			return rest_errors.New(http.StatusBadRequest, errors.New("invalid client id"))
		}
		if atr.ClientSecret == "" {
			return rest_errors.New(http.StatusBadRequest, errors.New("invalid client secret"))
		}
		break

	default:
		return rest_errors.New(http.StatusBadRequest, errors.New("invalid grant type"))
	}

	return nil
}

func (at *AccessToken) Generate() {
	at.AccessToken = cryptos.GetMd5(fmt.Sprintf("at-%d-%d-ran", at.UserID, at.Expires))
}
