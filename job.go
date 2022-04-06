package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type kv struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// JobManagerConfig returns the cluster configuration of
// job manager server.
func (c *Client) JobManagerConfig() ([]kv, error) {
	var r []kv
	req, err := http.NewRequest(
		"GET",
		c.url("/jobmanager/config"),
		nil,
	)
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

type metricValue struct {
	ID    string `json:"id"`
	Value string `json:"value,omitempty"`
}

// JobManagerMetrics provides access to job manager
// metrics.
func (c *Client) JobManagerMetrics(ids []string) ([]metricValue, error) {
	var r []metricValue
	req, err := http.NewRequest(
		"GET",
		c.url("/jobmanager/metrics"),
		nil,
	)
	if err != nil {
		return r, err
	}
	q := req.URL.Query()
	if len(ids) > 0 {
		q.Add("get", strings.Join(ids, ","))
	}

	b, err := c.client.Do(req)
	if err != nil {
		return r, err
	}
	err = json.Unmarshal(b, &r)
	return r, err
}

type jobsResp struct {
	Jobs []job `json:"jobs"`
}

type job struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// Jobs returns an overview over all jobs and their
// current state.
func (c *Client) Jobs() (jobsResp, error) {
	var r jobsResp
	req, err := http.NewRequest(
		"GET",
		c.url("/jobs"),
		nil,
	)
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

// SubmitJob submits a job.
func (c *Client) SubmitJob() error {
	return fmt.Errorf("not implement")
}

type JobMetricsOpts struct {
	// Metrics (optional): string values to select
	// specific metrics.
	Metrics []string

	// Agg (optional): list of aggregation modes which
	// should be calculated. Available aggregations are:
	// "min, max, sum, avg".
	Agg []string

	// Jobs (optional): job list of 32-character
	// hexadecimal strings to select specific jobs.
	Jobs []string
}

// JobMetrics provides access to aggregated job metrics.
func (c *Client) JobMetrics(opts JobMetricsOpts) (map[string]interface{}, error) {
	var r map[string]interface{}
	req, err := http.NewRequest(
		"GET",
		c.url("/jobs/metrics"),
		nil,
	)
	if err != nil {
		return r, err
	}
	q := req.URL.Query()
	if len(opts.Metrics) > 0 {
		q.Add("get", strings.Join(opts.Metrics, ","))
	}

	if len(opts.Agg) > 0 {
		q.Add("agg", strings.Join(opts.Agg, ","))
	}

	if len(opts.Jobs) > 0 {
		q.Add("jobs", strings.Join(opts.Jobs, ","))
	}

	req.URL.RawQuery = q.Encode()

	b, err := c.client.Do(req)
	if err != nil {
		return r, err
	}
	err = json.Unmarshal(b, &r)
	return r, err
}

type overviewResp struct {
	Jobs []jobOverview `json:"jobs"`
}

type jobOverview struct {
	ID               string `json:"jid"`
	Name             string `json:"name"`
	State            string `json:"state"`
	Start            int64  `json:"start-time"`
	End              int64  `json:"end-time"`
	Duration         int64  `json:"duration"`
	LastModification int64  `json:"last-modification"`
	Tasks            status `json:"tasks"`
}

type status struct {
	Total       int `json:"total,omitempty"`
	Created     int `json:"created"`
	Scheduled   int `json:"scheduled"`
	Deploying   int `json:"deploying"`
	Running     int `json:"running"`
	Finished    int `json:"finished"`
	Canceling   int `json:"canceling"`
	Canceled    int `json:"canceled"`
	Failed      int `json:"failed"`
	Reconciling int `json:"reconciling"`
}

// JobsOverview returns an overview over all jobs.
func (c *Client) JobsOverview() (overviewResp, error) {
	var r overviewResp
	req, err := http.NewRequest(
		"GET",
		c.url("/jobs/overview"),
		nil,
	)
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

type jobResp struct {
	ID          string `json:"jid"`
	Name        string `json:"name"`
	IsStoppable bool   `json:"isStoppable"`
	State       string `json:"state"`

	Start    int64 `json:"start-time"`
	End      int64 `json:"end-time"`
	Duration int64 `json:"duration"`
	Now      int64 `json:"now"`

	Timestamps   timestamps `json:"timestamps"`
	Vertices     []vertice  `json:"vertices"`
	StatusCounts status     `json:"status-counts"`
	Plan         plan       `json:"plan"`
}

type timestamps struct {
	Canceled    int64 `json:"CANCELED"`
	Suspended   int64 `json:"SUSPENDED"`
	Finished    int64 `json:"FINISHED"`
	Canceling   int64 `json:"CANCELLING"`
	Running     int64 `json:"RUNNING"`
	Restaring   int64 `json:"RESTARTING"`
	Reconciling int64 `json:"RECONCILING"`
	Created     int64 `json:"CREATED"`
	Failed      int64 `json:"FAILED"`
	Failing     int64 `json:"FAILING"`
}

type vertice struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Status      string                 `json:"status"`
	Parallelism int                    `json:"parallelism"`
	Start       int64                  `json:"start-time"`
	End         int64                  `json:"end-time"`
	Duration    int64                  `json:"duration"`
	Tasks       status                 `json:"tasks"`
	Metrics     map[string]interface{} `json:"metrics"`
}

// Job returns details of a job.
func (c *Client) Job(jobID string) (jobResp, error) {
	var r jobResp
	uri := fmt.Sprintf("/jobs/%s", jobID)
	req, err := http.NewRequest(
		"GET",
		c.url(uri),
		nil,
	)
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

// StopJob terminates a job.
func (c *Client) StopJob(jobID string) error {
	uri := fmt.Sprintf("/jobs/%s", jobID)
	req, err := http.NewRequest(
		"PATCH",
		c.url(uri),
		nil,
	)
	if err != nil {
		return err
	}
	_, err = c.client.Do(req)
	return err
}

type checkpointsResp struct {
	Counts  counts                     `json:"counts"`
	Summary summary                    `json:"summary"`
	Latest  latest                     `json:"latest"`
	History []failedCheckpointsStatics `json:"history"`
}

type counts struct {
	Restored   int `json:"restored"`
	Total      int `json:"total"`
	InProgress int `json:"in_progress"`
	Completed  int `json:"completed"`
	Failed     int `json:"failed"`
}

type summary struct {
	StateSize         statics `json:"state_size`
	End2EndDuration   statics `json:"end_to_end_duration"`
	AlignmentBuffered statics `json:"alignment_buffered"`
}

type statics struct {
	Min int `json:"min"`
	Max int `json:"max"`
	Avg int `json:"avg"`
}

type latest struct {
	Completed completedCheckpointsStatics `json:"completed"`
	Savepoint savepointsStatics           `json:"savepoint"`
	Failed    failedCheckpointsStatics    `json:"failed"`
	Restored  restoredCheckpointsStatics  `json:"restored"`
}

type completedCheckpointsStatics struct {
	ID                      string                 `json:"id"`
	Status                  string                 `json:"status"`
	IsSavepoint             bool                   `json:"is_savepoint"`
	TriggerTimestamp        int64                  `json:"trigger_timestamp"`
	LatestAckTimestamp      int64                  `json:"latest_ack_timestamp"`
	StateSize               int64                  `json:"state_size"`
	End2EndDuration         int64                  `json:"end_to_end_duration"`
	AlignmentBuffered       int64                  `json:"alignment_buffered"`
	NumSubtasks             int64                  `json:"num_subtasks"`
	NumAcknowledgedSubtasks int64                  `json:"num_acknowledged_subtasks"`
	tasks                   taskCheckpointsStatics `json:"tasks"`
	ExternalPath            string                 `json:"external_path"`
	Discarded               bool                   `json:"discarded"`
}

type savepointsStatics struct {
	ID                      int                    `json:"id"`
	Status                  string                 `json:"status"`
	IsSavepoint             bool                   `json:"is_savepoint"`
	TriggerTimestamp        int64                  `json:"trigger_timestamp"`
	LatestAckTimestamp      int64                  `json:"latest_ack_timestamp"`
	StateSize               int64                  `json:"state_size"`
	End2EndDuration         int64                  `json:"end_to_end_duration"`
	AlignmentBuffered       int64                  `json:"alignment_buffered"`
	NumSubtasks             int64                  `json:"num_subtasks"`
	NumAcknowledgedSubtasks int64                  `json:"num_acknowledged_subtasks"`
	tasks                   taskCheckpointsStatics `json:"tasks"`
	ExternalPath            string                 `json:"external_path"`
	Discarded               bool                   `json:"discarded"`
}
type taskCheckpointsStatics struct {
	ID     string `json:"id"`
	Status string `json:"status"`

	LatestAckTimestamp int64 `json:"latest_ack_timestamp"`

	FailureTimestamp int64  `json:"failure_timestamp"`
	FailureMessage   string `json:"failure_message"`

	StateSize               int64 `json:"state_size"`
	End2EndDuration         int64 `json:"end_to_end_duration"`
	AlignmentBuffered       int64 `json:"alignment_buffered"`
	NumSubtasks             int64 `json:"num_subtasks"`
	NumAcknowledgedSubtasks int64 `json:"num_acknowledged_subtasks"`
}

type failedCheckpointsStatics struct {
	ID                      int64                  `json:"id"`
	Status                  string                 `json:"status"`
	IsSavepoint             bool                   `json:"is_savepoint"`
	TriggerTimestamp        int64                  `json:"trigger_timestamp"`
	LatestAckTimestamp      int64                  `json:"latest_ack_timestamp"`
	StateSize               int64                  `json:"state_size"`
	End2EndDuration         int64                  `json:"end_to_end_duration"`
	AlignmentBuffered       int64                  `json:"alignment_buffered"`
	NumSubtasks             int64                  `json:"num_subtasks"`
	NumAcknowledgedSubtasks int64                  `json:"num_acknowledged_subtasks"`
	tasks                   taskCheckpointsStatics `json:"tasks"`
}

type restoredCheckpointsStatics struct {
	ID               int64  `json:"id"`
	RestoreTimestamp int64  `json:"restore_timestamp"`
	IsSavepoint      bool   `json:"is_savepoint"`
	ExternalPath     string `json:"external_path"`
}

// Checkpoints returns checkpointing statistics for a job.
func (c *Client) Checkpoints(jobID string) (checkpointsResp, error) {
	var r checkpointsResp
	uri := fmt.Sprintf("/jobs/%s/checkpoints", jobID)
	req, err := http.NewRequest(
		"GET",
		c.url(uri),
		nil,
	)
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

type savePointsResp struct {
	RequestID string `json:"request-id"`
}

// SavePoints triggers a savepoint, and optionally cancels the
// job afterwards. This async operation would return a
// 'triggerid' for further query identifier.
func (c *Client) SavePoints(jobID string, saveDir string, cancleJob bool) (savePointsResp, error) {
	var r savePointsResp

	type savePointsReq struct {
		SaveDir   string `json:"target-directory"`
		CancleJob bool   `json:"cancel-job"`
	}

	d := savePointsReq{
		SaveDir:   saveDir,
		CancleJob: cancleJob,
	}
	data := new(bytes.Buffer)
	json.NewEncoder(data).Encode(d)
	uri := fmt.Sprintf("/jobs/%s/savepoints", jobID)
	req, err := http.NewRequest(
		"POST",
		c.url(uri),
		data,
	)
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

type stopJobResp struct {
	RequestID string `json:"request-id"`
}

// StopJob stops a job with a savepoint. Optionally, it can also
// emit a MAX_WATERMARK before taking the savepoint to flush out
// any state waiting for timers to fire. This async operation
// would return a 'triggerid' for further query identifier.
func (c *Client) StopJobWithSavepoint(jobID string, saveDir string, drain bool) (stopJobResp, error) {
	var r stopJobResp
	type stopJobReq struct {
		SaveDir string `json:"targetDirectory"`
		Drain   bool   `json:"drain"`
	}

	d := stopJobReq{
		SaveDir: saveDir,
		Drain:   drain,
	}
	data := new(bytes.Buffer)
	json.NewEncoder(data).Encode(d)
	uri := fmt.Sprintf("/jobs/%s/stop", jobID)
	req, err := http.NewRequest(
		"POST",
		c.url(uri),
		data,
	)
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
