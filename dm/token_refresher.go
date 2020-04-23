package dm

import "github.com/gdey/forge-api-go-client/oauth"

type TokenRefresher interface {
	Bearer() *oauth.Bearer
	RefreshTokenIfRequired(auth oauth.ThreeLeggedAuth) error
}
