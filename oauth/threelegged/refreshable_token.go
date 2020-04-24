package threelegged

import (
	"errors"
	"sync"
	"time"

	"github.com/gdey/forge-api-go-client/oauth"
)

type RefreshableToken struct {
	bearer          *oauth.Bearer
	TokenExpireTime time.Time
	readMutex       sync.Mutex
	writeMutex      sync.Mutex
}

func NewRefreshableToken(bearer *oauth.Bearer, expiryTime time.Time) *RefreshableToken {
	return &RefreshableToken{
		bearer:          bearer,
		TokenExpireTime: expiryTime,
	}
}

func (t *RefreshableToken) RefreshTokenIfRequired(auth Auth) error {
	if t == nil {
		return errors.New("Invalid Token")
	}

	// Check if token has expired
	now := time.Now()
	expiryTime := t.TokenExpireTime
	if now.Before(expiryTime) {
		return nil
	}

	refreshedBearer, err := auth.RefreshToken(t.bearer.RefreshToken)
	if err != nil {
		return err
	}

	// Refresh "now" and add new token expiration time to API struct along with new credentials
	now = time.Now()
	newExpiryTime := now.Add(time.Second * time.Duration(refreshedBearer.ExpiresIn))

	t.writeMutex.Lock()
	t.TokenExpireTime = newExpiryTime

	t.bearer.AccessToken = refreshedBearer.AccessToken
	t.bearer.ExpiresIn = refreshedBearer.ExpiresIn
	t.bearer.RefreshToken = refreshedBearer.RefreshToken
	t.bearer.TokenType = refreshedBearer.TokenType
	t.writeMutex.Unlock()

	return nil
}

func (t *RefreshableToken) Bearer() *oauth.Bearer {
	t.readMutex.Lock()
	defer t.readMutex.Unlock()
	return t.bearer
}
