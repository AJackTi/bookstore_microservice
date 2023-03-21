package access_token

import (
	"fmt"
	"strings"
	"time"

	"github.com/AJackTi/bookstore_oauth-api/src/utils/cryptos"
	"github.com/AJackTi/bookstore_oauth-api/src/utils/errors"
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

func (at *AccessToken) Validate() *errors.RestErr {
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if at.AccessToken == "" {
		return errors.New(500)
	}
	if at.UserID <= 0 {
		return errors.New(400, "invalid user")
	}
	if at.ClientID <= 0 {
		return errors.New(400, "invalid client id")
	}
	if at.Expires <= 0 {
		return errors.New(400, "invalid Expired time")
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

func (atr *AccessTokenRequest) Validate() *errors.RestErr {

	atr.GrantType = strings.TrimSpace(atr.GrantType)

	switch atr.GrantType {
	case grantTypePassword:
		atr.UserName = strings.TrimSpace(atr.UserName)
		atr.Password = strings.TrimSpace(atr.Password)
		if atr.UserName == "" {
			return errors.New(400)
		}
		if atr.Password == "" {
			return errors.New(400)
		}
		break

	case grandTypeClientCredentials:
		if atr.ClientId == "" {
			return errors.New(400)
		}
		if atr.ClientSecret == "" {
			return errors.New(400)
		}
		break

	default:
		return errors.New(400, "invalid grantype")
	}

	return nil
}

func (at *AccessToken) Generate() {
	at.AccessToken = cryptos.GetMd5(fmt.Sprintf("at-%d-%d-ran", at.UserID, at.Expires))
}
