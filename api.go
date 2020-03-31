package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Client reprents flink REST API client
type Client struct {
	// Addr reprents flink job manager server address
	Addr string

	client *httpClient
}

// New returns a flink client
func New(addr string) (*Client, error) {
	return &Client{
		Addr:   addr,
		client: newHttpClient(),
	}, nil
}

func (c *Client) url(path string) string {
	if strings.HasPrefix(c.Addr, "http") {
		return fmt.Sprintf("%s%s", c.Addr, path)
	}
	return fmt.Sprintf("http://%s%s", c.Addr, path)
}

// Shutdown shutdown the flink cluster
func (c *Client) Shutdown() error {
	req, err := http.NewRequest("DELETE", c.url("/cluster"), nil)
	if err != nil {
		return err
	}
	_, err = c.client.Do(req)
	return err
}

type configResp struct {
	RefreshInterval int64    `json:"refresh-interval"`
	TimezoneName    string   `json:"timezone-name"`
	TimezoneOffset  int64    `json:"timezone-offset"`
	FlinkVersion    string   `json:"flink-version"`
	FlinkRevision   string   `json:"flink-revision"`
	Features        features `json:"features"`
}
type features struct {
	WebSubmit bool `json:"web-submit"`
}

// Config returns the configuration of the WebUI
func (c *Client) Config() (configResp, error) {
	var r configResp
	req, err := http.NewRequest("GET", c.url("/config"), nil)
	if err != nil {
		return r, err
	}
	b, err := c.client.Do(req)
	err = json.Unmarshal(b, &r)
	return r, err
}

type uploadResp struct {
	FileName string `json:"filename"`
	Status   string `json:"status"`
}

// Upload uploads jar file
func (c *Client) UploadJar(fpath string) (uploadResp, error) {
	var r uploadResp
	file, err := os.Open(fpath)
	if err != nil {
		return r, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("jarfile", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()

	req, err := http.NewRequest("POST", c.url("/jars/upload"), body)
	if err != nil {
		return r, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	b, err := c.client.Do(req)
	err = json.Unmarshal(b, &r)
	return r, err
}

type jarsResp struct {
	Address string    `json:"address"`
	Files   []jarFile `json:"files"`
}

type jarFile struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Uploaded int64   `json:"uploaded"`
	Entries  []entry `json:"entry"`
}

type entry struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Jars eturns a list of all jars previously uploaded
// via '/jars/upload'
func (c *Client) Jars() (jarsResp, error) {
	var r jarsResp
	req, err := http.NewRequest("GET", c.url("/jars"), nil)
	if err != nil {
		return r, err
	}
	b, err := c.client.Do(req)
	err = json.Unmarshal(b, &r)
	return r, err
}

// DeleteJar deletes a jar file
func (c *Client) DeleteJar(jarid string) error {
	uri := fmt.Sprintf("/jars/%s", jarid)
	req, err := http.NewRequest("DELETE", c.url(uri), nil)
	if err != nil {
		return err
	}
	_, err = c.client.Do(req)
	return err
}
