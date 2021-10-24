// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oauth2

import (
	"context"
	"errors"
	"net/url"

	"golang.org/x/oauth2"
)

// tokenRefresher is a TokenSource that makes "grant_type"=="refresh_token"
// HTTP requests to renew a token using a RefreshToken.
type tokenRefresher struct {
	ctx          context.Context // used to get HTTP requests
	conf         *Config
	refreshToken string
}

// WARNING: Token is not safe for concurrent access, as it
// updates the tokenRefresher's refreshToken field.
// Within this package, it is used by reuseTokenSource which
// synchronizes calls to this method with its own mutex.
//
// https://developer.withings.com/api-reference#operation/oauth2-refreshaccesstoken
func (tf *tokenRefresher) Token() (*oauth2.Token, error) {
	if tf.refreshToken == "" {
		return nil, errors.New("oauth2: token expired and refresh token is not set")
	}

	tk, err := retrieveToken(tf.ctx, tf.conf, url.Values{
		"action":        {"requesttoken"},
		"grant_type":    {"refresh_token"},
		"refresh_token": {tf.refreshToken},
	})

	if err != nil {
		return nil, err
	}
	if tf.refreshToken != tk.RefreshToken {
		tf.refreshToken = tk.RefreshToken
	}
	return tk, err
}
