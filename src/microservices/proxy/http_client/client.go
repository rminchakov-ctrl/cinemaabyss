package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) DoRequest(request *http.Request) (body []byte, err error) {
	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("can't do http request: %w", err)
	}

	if resp.Body == nil {
		return nil, fmt.Errorf("empty reader")
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)

	var dst json.RawMessage
	if err = json.Unmarshal(body, &dst); err != nil {
		return nil, fmt.Errorf("can't unmarshal response body: %w", err)
	}

	return dst, nil
}

func (c *Client) NewGetRequest(
	ctx context.Context, path string,
) (*http.Request, error) {
	u := fmt.Sprintf("%s/%s", c.baseURL, path)
	return http.NewRequestWithContext(ctx, http.MethodGet, u, nil)

}

func (c *Client) NewDeleteRequest(
	ctx context.Context, path string,
) (*http.Request, error) {
	u := fmt.Sprintf("%s/%s", c.baseURL, path)
	return http.NewRequestWithContext(ctx, http.MethodDelete, u, nil)

}

func (c *Client) NewPostRequest(
	ctx context.Context, path string, body any,
) (*http.Request, error) {
	return newRequestWithBody(ctx, path, http.MethodPost, c.baseURL, body)
}

func (c *Client) NewPutRequest(
	ctx context.Context, path string, body any,
) (*http.Request, error) {
	return newRequestWithBody(ctx, path, http.MethodPut, c.baseURL, body)
}

func (c *Client) NewPatchRequest(
	ctx context.Context, path string, body any,
) (*http.Request, error) {
	return newRequestWithBody(ctx, path, http.MethodPatch, c.baseURL, body)
}

func newRequestWithBody(ctx context.Context, path, method, baseURL string, body any) (*http.Request, error) {
	u := fmt.Sprintf("%s/%s", baseURL, path)

	var bodyReader io.Reader
	if body != nil {
		bodyRaw, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling error: %w", err)
		}

		bodyReader = bytes.NewReader(bodyRaw)
	}

	request, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	return request, nil
}
