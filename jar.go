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
	"strconv"
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
	if err != nil {
		return r, err
	}
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
	if err != nil {
		return r, err
	}
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
	if err != nil {
		return r, err
	}
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

type planResp struct {
	Plan plan `json:"plan"`
}

type plan struct {
	JID   string `json:"jid"`
	Name  string `json:"name"`
	Nodes []node `json:"nodes"`
}

type node struct {
	ID               string  `json:"id"`
	Parallelism      int     `json:"parallelism"`
	Operator         string  `json:"operator"`
	OperatorStrategy string  `json:"operator_strategy"`
	Description      string  `json:"description"`
	Inputs           []input `json:"inputs"`
}

type input struct {
	Num          int    `json:"num"`
	ID           string `json:"id"`
	ShipStrategy string `json:"ship_strategy"`
	Exchange     string `json:"exchange"`
}

// PlanJar returns the dataflow plan of a job contained
// in a jar previously uploaded via '/jars/upload'.
// Todo: support more args.
func (c *Client) PlanJar(jarid string) (planResp, error) {
	var r planResp
	uri := fmt.Sprintf("/jars/%s/plan", jarid)
	req, err := http.NewRequest("GET", c.url(uri), nil)
	if err != nil {
		return r, err
	}
	b, err := c.client.Do(req)
	if err != nil {
		return r, err
	}
	err = json.Unmarshal(b, &r)
	return r, err
}

type runResp struct {
	ID string
}

type RunOpts struct {
	// JarID: String value that identifies a jar. When
	// uploading the jar a path is returned, where the
	// filename is the ID.
	JarID string

	// AllowNonRestoredState(optional): Boolean value that
	// specifies whether the job submission should be
	// rejected if the savepoint contains state that
	// cannot be mapped back to the job.
	AllowNonRestoredState bool

	// SavepointPath (optional): String value that
	// specifies the path of the savepoint to restore the
	// job from.
	SavepointPath string

	// programArg (optional): list of program arguments.
	ProgramArg []string

	// EntryClass (optional): String value that specifies
	// the fully qualified name of the entry point class.
	// Overrides the class defined in the jar file
	// manifest.
	EntryClass string

	// Parallelism (optional): Positive integer value that
	// specifies the desired parallelism for the job.
	Parallelism int
}

// RunJar submits a job by running a jar previously
// uploaded via '/jars/upload'.
func (c *Client) RunJar(opts RunOpts) (runResp, error) {
	var r runResp
	uri := fmt.Sprintf("/jars/%s/run", opts.JarID)
	req, err := http.NewRequest("POST", c.url(uri), nil)
	q := req.URL.Query()
	if opts.SavepointPath != "" {
		q.Add("savepointPath", opts.SavepointPath)
		q.Add("allowNonRestoredState", strconv.FormatBool(opts.AllowNonRestoredState))
	}
	if len(opts.ProgramArg) > 0 {
		q.Add("programArg", strings.Join(opts.ProgramArg, ","))
	}
	if opts.EntryClass != "" {
		q.Add("entry-class", opts.EntryClass)
	}
	if opts.Parallelism > 0 {
		q.Add("parallelism", strconv.Itoa(opts.Parallelism))
	}
	req.URL.RawQuery = q.Encode()
	if err != nil {
		return r, err
	}
	b, err := c.client.Do(req)
	if err != nil {
		return r, err
	}
	err = json.Unmarshal(b, &r)
	return r, err
}
