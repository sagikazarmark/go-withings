package internal

// withingsTokenResponse wraps tokenJSON.
//
// https://support.withings.com/hc/en-us/articles/360016745358--BREAKING-Deprecating-access-and-refresh-tokens-endpoints
type withingsTokenResponse struct {
	Status int       `json:"status"`
	Body   tokenJSON `json:"body"`
}
