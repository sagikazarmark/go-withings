package withings

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/mitchellh/mapstructure"
)

const (
	endpoint      = "https://wbsapi.withings.net/"
	endpointHIPAA = "https://wbsapi.us.withingsmed.net/"

	userAgent = "go-withings"
)

// A Client manages communication with the Withings API.
type Client struct {
	client *http.Client // HTTP client used to communicate with the API.

	// Base URL for API requests. Defaults to the public Withings API, but can be
	// set to a different URL to use the Withings HIPAA endpoint.
	// BaseURL should always be specified with a trailing slash.
	BaseURL *url.URL

	// User agent used when communicating with the Withings API.
	UserAgent string

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// Services used for talking to different parts of the Withings API.
	Measure *MeasureService
	Heart   *HeartService
	Sleep   *SleepService
	Notify  *NotifyService
}

type service struct {
	client *Client
}

// NewClient returns a new Withings API client for the Public endpoint.
// Provide an http.Client that will perform the authentication
// (such as that provided by the golang.org/x/oauth2 library).
func NewClient(httpClient *http.Client) *Client {
	return newClient(httpClient, endpoint)
}

// NewHIPAAClient returns a new Withings API client for the HIPAA endpoint.
// Provide an http.Client that will perform the authentication
// (such as that provided by the golang.org/x/oauth2 library).
func NewHIPAAClient(httpClient *http.Client) *Client {
	return newClient(httpClient, endpointHIPAA)
}

func newClient(httpClient *http.Client, endpoint string) *Client {
	baseURL, _ := url.Parse(endpoint)

	c := &Client{
		client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: userAgent,
	}

	c.common.client = c

	c.Measure = (*MeasureService)(&c.common)
	c.Heart = (*HeartService)(&c.common)
	c.Sleep = (*SleepService)(&c.common)
	c.Notify = (*NotifyService)(&c.common)

	return c
}

// Response is a Withings API response. This wraps the standard http.Response
// returned from Withings and provides convenient access to things like
// pagination offset.
type Response struct {
	HttpResponse *http.Response

	// Status code returned from the Withings API.
	//
	// Withings API docs: https://developer.withings.com/api-reference#section/Response-status
	Status int

	// These fields provide information whether there is more data to fetch.
	// If more is true, sending a new request with the offset will return
	// the next set of results.
	More   bool
	Offset int
}

// newResponse creates a new Response for the provided http.Response.
// r must not be nil.
func newResponse(r *http.Response) *Response {
	response := &Response{HttpResponse: r}

	return response
}

// PostForm sends a POST request to the API
// with data's keys and values URL-encoded as the request body.
//
// The Content-Type header is set to application/x-www-form-urlencoded.
// To set other headers, use NewRequest and Do.
//
// A relative URL can be provided in url,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.
func (c *Client) PostForm(ctx context.Context, url string, data url.Values, v interface{}) (resp *Response, err error) {
	req, err := c.NewRequest(ctx, http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.Do(req, v)
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.
func (c *Client) NewRequest(ctx context.Context, method string, urlStr string, body io.Reader) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}

	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}

var errNonNilContext = errors.New("context must be non-nil")

// BareDo sends an API request and lets you handle the api response. If an error
// or API Error occurs, the error will contain more information. Otherwise you
// are supposed to read and close the response's Body.
func (c *Client) BareDo(req *http.Request) (*Response, error) {
	if req.Context() == nil {
		return nil, errNonNilContext
	}

	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		default:
		}

		return nil, err
	}

	response := newResponse(resp)

	return response, err
}

type apiResponse struct {
	Status int `json:"status"`

	Body struct {
		More   bool `json:"more"`
		Offset int  `json:"offset"`
	} `json:"body"`
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.
// If v is nil, and no error hapens, the response is returned as is.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.BareDo(req)
	if err != nil {
		return resp, err
	}
	defer resp.HttpResponse.Body.Close()

	body := map[string]interface{}{}

	err = json.NewDecoder(resp.HttpResponse.Body).Decode(&body)
	if err != nil && err != io.EOF { // nolint: errorlint // ignore EOF errors caused by empty response body
		return resp, err
	}

	var apiResp apiResponse

	err = decode(body, &apiResp)
	if err != nil {
		return resp, err
	}

	err = decode(body, v)
	if err != nil {
		return resp, err
	}

	resp.Status = apiResp.Status
	resp.More = apiResp.Body.More
	resp.Offset = apiResp.Body.Offset

	return resp, err
}

func decode(input interface{}, output interface{}) error {
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   output,
		TagName:  "json",
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}
