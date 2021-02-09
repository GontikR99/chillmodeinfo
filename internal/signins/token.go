// +build server

package signins

import (
	"context"
	"errors"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"net/http"
)

func ValidateToken(ctx context.Context, idToken string) (string, error) {
	oauth2Service, err := oauth2.NewService(ctx, option.WithHTTPClient(&http.Client{}))
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return "", err
	}
	return tokenInfo.UserId, nil
}

func ValidateClientId(ctx context.Context, clientId string) (string, error) {
	return "", errors.New("this client hasn't signed in yet")
}
