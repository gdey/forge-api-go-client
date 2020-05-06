package threelegged

import (
	"github.com/gdey/forge-api-go-client/oauth"
)

type TokenRefresher interface {
	Bearer() *oauth.Bearer
	RefreshTokenIfRequired(auth Auth) error
}

type AuthRefresher interface {
	RefreshToken(refresh_token string) (*oauth.Bearer, error)
}
