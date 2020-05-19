package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

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
	if err := auth.SetAuthHeader(scope, req.Header); err != nil {
		return nil, fmt.Errorf("DoRawRequest:%w", err)
	}

	return client.Do(req)
}

func (c *Client) ProcessRawError(response *http.Response, result interface{}) (err error) {
	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		return ErrResult{StatusCode: response.StatusCode, Reason: string(content)}
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
		err = ErrResult{StatusCode: response.StatusCode}
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

START:
	res, err := c.DoRawRequest(ctx, method, scope, paths, filters, nil, contentType, body)
	if err != nil {
		return fmt.Errorf("error making request to %v %v : %w", method, strings.Join(paths, "/"), err)
	}
	defer res.Body.Close()
	err = c.ProcessJSONError(res, result)
	var errResult ErrResult
	if err != nil && errors.As(err, &errResult) {
		switch {
		case errResult.IsRateLimited():
			// we need to wait for a bit and then retry
			<-time.After(30 * time.Second)
			goto START
		case errResult.StatusCode == http.StatusUnsupportedMediaType:
			// This is lke a 500 error, however something is wrong with
			// the call. (We are using the wrong media type.) Did the
			// api change and we need to upgrade the api for this end-point
			return ErrAPIIncompatible{
				Err:          errResult,
				DeveloperMsg: fmt.Sprintf("incorrect MediaType: sent %v for %v(%v)", contentType, method, strings.Join(paths, "/")),
			}
		default:
			return errResult
		}
	}
	return err
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
