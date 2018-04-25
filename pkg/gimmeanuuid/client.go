package gimmeanuuid

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

// Client is the interface providing methods to interact with the gimme-an-uuid api.
type Client interface {
	TextPlainV1() (string, error)
	TextPlainV2(domain int) (string, error)
	TextPlainV3(namespace, name string) (string, error)
	TextPlainV4() (string, error)
	TextPlainV5(namespace, name string) (string, error)
}

type client struct {
	baseURL    url.URL
	httpClient *http.Client
}

const apiPath string = "/api/uuid"

// NewClient returns a new gimme an uuid client that uses the given http.Client.
func NewClient(httpClient *http.Client, baseURL url.URL) (Client, error) {
	if httpClient == nil {
		return nil, errors.New("httpClient has to be provided")
	}
	p := baseURL.EscapedPath()
	if p != apiPath {
		baseURL.Path = path.Join(p, apiPath)
	}
	return &client{httpClient: httpClient, baseURL: baseURL}, nil
}

func textPlainRequest(httpClient *http.Client, requestURL url.URL) (string, error) {
	req, err := http.NewRequest("GET", requestURL.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "text/plain")
	res, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	s := string(b)
	if s == "" {
		return "", errors.New("response body empty")
	}
	return s, nil
}

func (c *client) TextPlainV1() (string, error) {
	u := c.baseURL
	u.Path = path.Join(u.EscapedPath(), "v1")
	return textPlainRequest(c.httpClient, u)
}

func (c *client) TextPlainV2(domain int) (string, error) {
	u := c.baseURL
	u.Path = path.Join(u.EscapedPath(), "v2", string(domain))
	return textPlainRequest(c.httpClient, u)
}

func (c *client) TextPlainV3(namespace, name string) (string, error) {
	u := c.baseURL
	u.Path = path.Join(u.EscapedPath(), "v3", namespace, name)
	return textPlainRequest(c.httpClient, u)
}

func (c *client) TextPlainV4() (string, error) {
	u := c.baseURL
	u.Path = path.Join(u.EscapedPath(), "v4")
	return textPlainRequest(c.httpClient, u)
}

func (c *client) TextPlainV5(namespace, name string) (string, error) {
	u := c.baseURL
	u.Path = path.Join(u.EscapedPath(), "v5", namespace, name)
	return textPlainRequest(c.httpClient, u)
}
