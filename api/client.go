package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gdey/forge-api-go-client/oauth"
	"github.com/gdey/forge-api-go-client/oauth/none"
	"github.com/gdey/forge-api-go-client/oauth/scopes"
)

const (
	ContentTypeJSON        = "application/json"
	ContentTypeFormEncoded = "application/x-www-form-urlencoded"
)

type Filterer interface {
	Add(url.Values) error
}

// ErrorResult reflects the body content when a request failed (g.e. Bad request or key conflict)
type ErrorResult struct {
	Reason     string `json:"reason"`
	StatusCode int
}

func (e ErrorResult) Error() string {
	return fmt.Sprintf("[%d]`%s`", e.StatusCode, e.Reason)
}

type Client struct {
	Client http.Client
	oauth.ForgeAuthenticator
}

func NewClient(auth oauth.ForgeAuthenticator) *Client {
	return &Client{
		ForgeAuthenticator: auth,
	}
}

func (c *Client) DoRawRequest(ctx context.Context, method string, scope scopes.Scope, paths []string, filters []Filterer, setHeaders func(http.Header) error, contentType string, body io.Reader) (*http.Response, error) {
	var (
		client http.Client
		auth   oauth.ForgeAuthenticator = none.None
	)
	if c != nil {
		client = c.Client
		if c.ForgeAuthenticator != nil {
			auth = c.ForgeAuthenticator
		}
	}
	path := auth.Path(paths...)

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		path,
		body,
	)
	if err != nil {
		return nil, err
	}

	if len(filters) > 0 {
		query := req.URL.Query()
		for _, filter := range filters {
			if filter == nil {
				continue
			}
			if err = filter.Add(query); err != nil {
				return nil, err
			}
		}
		req.URL.RawQuery = query.Encode()
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if setHeaders != nil {
		setHeaders(req.Header)
	}
	auth.SetAuthHeader(scope, req.Header)

	return client.Do(req)
}

func (c *Client) ProcessRawError(response *http.Response, result interface{}) (err error) {
	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		return ErrorResult{StatusCode: response.StatusCode, Reason: string(content)}
	}
	if result == nil {
		return nil
	}
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(result)
	if err != nil {
		return fmt.Errorf("JSON Decode: %w", err)
	}
	return nil
}

func (c *Client) ProcessJSONError(response *http.Response, result interface{}) (err error) {
	decoder := json.NewDecoder(response.Body)
	if response.StatusCode != http.StatusOK {
		err = ErrorResult{StatusCode: response.StatusCode}
		_ = decoder.Decode(&err)
		return err
	}
	if result == nil {
		return nil
	}
	err = decoder.Decode(result)
	if err != nil {
		return fmt.Errorf("JSON Decode: %w", err)
	}
	return nil
}

func (c *Client) DoRequest(ctx context.Context, method string, scope scopes.Scope, paths []string, result interface{}, filters []Filterer, contentType string, body io.Reader) error {

	res, err := c.DoRawRequest(ctx, method, scope, paths, filters, nil, contentType, body)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return c.ProcessJSONError(res, result)
}
func (c *Client) Post(ctx context.Context, scope scopes.Scope, paths []string, result interface{}, contentType string, body io.Reader) error {
	return c.DoRequest(ctx, http.MethodPost,
		scope,
		paths,
		result,
		nil,
		contentType,
		body,
	)
}
func (c *Client) Put(ctx context.Context, scope scopes.Scope, paths []string, result interface{}, contentType string, body io.Reader) error {
	return c.DoRequest(ctx, http.MethodPut,
		scope,
		paths,
		result,
		nil,
		contentType,
		body,
	)
}
func (c *Client) Get(ctx context.Context, scope scopes.Scope, paths []string, result interface{}, filters ...Filterer) error {
	return c.DoRequest(ctx, http.MethodGet,
		scope,
		paths,
		result,
		filters,
		ContentTypeJSON,
		nil,
	)
}
func (c *Client) Delete(ctx context.Context, scope scopes.Scope, paths []string) error {
	return c.DoRequest(ctx, http.MethodDelete,
		scope,
		paths,
		nil,
		nil,
		ContentTypeJSON,
		nil,
	)
}
