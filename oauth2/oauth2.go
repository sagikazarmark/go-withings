// Package oauth2 extends "golang.org/x/oauth2" to support using OAuth2 to access Withings.
//
// Notable deviations from standard OAuth2:
//
// The "scope" parameter is comma separated, not space separated.
//
// The token endpoint requires an additional "action" parameter to
// differentiate between authorization code exchange and token refresh.
//
// The token response payload is wrapped in a custom response envelope
// applied to every API response.
// More details about this change: https://support.withings.com/hc/en-us/articles/360016745358--BREAKING-Deprecating-access-and-refresh-tokens-endpoints
//
// Read more about the Withings OAuth2 API here: https://developer.withings.com/api-reference#tag/oauth2
package oauth2

import (
	"context"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
)

// Endpoint is the Public Cloud endpoint for Withings.
//
// https://developer.withings.com/developer-guide/getting-started/register-to-withings-api#public-endpoint
var Endpoint = oauth2.Endpoint{
	AuthURL:   "https://account.withings.com/oauth2_user/authorize2",
	TokenURL:  "https://wbsapi.withings.net/v2/oauth2",
	AuthStyle: oauth2.AuthStyleInParams,
}

// Endpoint is the HIPAA endpoint for Withings.
//
// https://developer.withings.com/developer-guide/getting-started/register-to-withings-hipaa-api#hipaa-endpoint
var EndpointHIPAA = oauth2.Endpoint{
	AuthURL:   "https://account.us.withingsmed.com/oauth2_user/authorize2",
	TokenURL:  "https://wbsapi.us.withingsmed.net/v2/oauth2",
	AuthStyle: oauth2.AuthStyleInParams,
}

// ModeDemo automatically logs the user in as a demo user.
//
// https://developer.withings.com/developer-guide/data-api/demo-user#demo-user
var ModeDemo = oauth2.SetAuthURLParam("mode", "demo")

// Config wraps a golang.org/x/oauth2.Config struct to extend its
// functionality with Withings specific behavior.
type Config struct {
	*oauth2.Config
}

// AuthCodeURL returns a URL to OAuth 2.0 provider's consent page
// that asks for permissions for the required scopes explicitly.
//
// It wraps the same function from golang.org/x/oauth2.Config and
// makes the scope parameter a comma separated string.
//
// https://developer.withings.com/api-reference#operation/oauth2-authorize
func (c *Config) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	if len(c.Scopes) > 0 {
		opts = append([]oauth2.AuthCodeOption{oauth2.SetAuthURLParam("scope", strings.Join(c.Scopes, ","))}, opts...)
	}

	return c.Config.AuthCodeURL(state, opts...)
}

// Exchange converts an authorization code into a token.
//
// It reimplements parts of the same function from golang.org/x/oauth2.Config
// to make it compatible with the Withings API.
//
// https://developer.withings.com/api-reference#operation/oauth2-getaccesstoken
func (c *Config) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	v := url.Values{
		"action":     {"requesttoken"},
		"grant_type": {"authorization_code"},
		"code":       {code},
	}
	if c.RedirectURL != "" {
		v.Set("redirect_uri", c.RedirectURL)
	}
	return retrieveToken(ctx, c, v)
}
