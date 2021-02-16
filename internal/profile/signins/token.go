// +build server

package signins

import (
	"context"
	"github.com/GontikR99/chillmodeinfo/internal/dao"
	"github.com/GontikR99/chillmodeinfo/internal/comms/httputil"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"net/http"
	"strings"
)

func validateGoogleIdToken(ctx context.Context, idToken string) (string, error) {
	oauth2Service, err := oauth2.NewService(ctx, option.WithHTTPClient(&http.Client{}))
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return "", err
	}
	userId := IdTypeGoogle + tokenInfo.UserId
	return userId, dao.RegisterUser(userId, tokenInfo.Email)
}

func ValidateToken(ctx context.Context, idToken string) (string, error) {
	if strings.HasPrefix(idToken, TokenGoogle) {
		return validateGoogleIdToken(ctx, idToken[len(TokenGoogle):])
	} else if strings.HasPrefix(idToken, TokenClientId) {
		gId, present, err := dao.LookupClientId(idToken[len(TokenClientId):])
		if err != nil {
			return "", err
		} else if !present {
			return "", httputil.NewError(http.StatusForbidden, "This client hasn't signed in yet")
		} else {
			return gId, nil
		}
	} else {
		return "", httputil.NewError(http.StatusForbidden, "No identity provided, or unsupported token type")
	}
}
