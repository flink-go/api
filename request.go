package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type httpClient struct {
	client http.Client
}

func newHttpClient() *httpClient {
	return &httpClient{
		client: http.Client{},
	}
}

func (c *httpClient) Do(req *http.Request) ([]byte, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if int(resp.StatusCode/100) != 2 {
		return nil, fmt.Errorf("http status not 2xx: %d %s", resp.StatusCode, string(body))
	}
	return body, nil
}
