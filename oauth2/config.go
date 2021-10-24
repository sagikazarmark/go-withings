// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oauth2

import (
	"context"
	"net/http"
	"net/url"

	"github.com/sagikazarmark/go-withings/oauth2/internal"
	"golang.org/x/oauth2"
)

// Client returns an HTTP client using the provided token.
// The token will auto-refresh as necessary. The underlying
// HTTP transport will be obtained using the provided context.
// The returned client and its Transport should not be modified.
func (c *WithingsConfig) Client(ctx context.Context, t *oauth2.Token) *http.Client {
	return oauth2.NewClient(ctx, c.TokenSource(ctx, t))
}

// TokenSource returns a TokenSource that returns t until t expires,
// automatically refreshing it as necessary using the provided context.
//
// Most users will use Config.Client instead.
func (c *WithingsConfig) TokenSource(ctx context.Context, t *oauth2.Token) oauth2.TokenSource {
	tkr := &tokenRefresher{
		ctx:  ctx,
		conf: c.Config,
	}
	if t != nil {
		tkr.refreshToken = t.RefreshToken
	}
	return oauth2.ReuseTokenSource(t, tkr)
}

// retrieveToken takes a *Config and uses that to retrieve an *internal.Token.
// This token is then mapped from *internal.Token into an *oauth2.Token which is returned along
// with an error..
func retrieveToken(ctx context.Context, c *Config, v url.Values) (*oauth2.Token, error) {
	tk, err := internal.RetrieveToken(ctx, c.ClientID, c.ClientSecret, c.Endpoint.TokenURL, v, internal.AuthStyle(c.Endpoint.AuthStyle))
	if err != nil {
		if rErr, ok := err.(*internal.RetrieveError); ok {
			return nil, (*oauth2.RetrieveError)(rErr)
		}
		return nil, err
	}
	return tokenFromInternal(tk), nil
}

// tokenFromInternal maps an *internal.Token struct into
// a *Token struct.
func tokenFromInternal(t *internal.Token) *oauth2.Token {
	if t == nil {
		return nil
	}
	return (&oauth2.Token{
		AccessToken:  t.AccessToken,
		TokenType:    t.TokenType,
		RefreshToken: t.RefreshToken,
		Expiry:       t.Expiry,
		//raw:          t.Raw,
	}).WithExtra(t.Raw)
}
