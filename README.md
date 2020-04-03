# Flink Monitoring API Golang library

Detail doc: https://ci.apache.org/projects/flink/flink-docs-stable/monitoring/rest_api.html

Status: Beta


```
package main

import (
	"fmt"

	"github.com/flink-go/api"
)

func main() {
	// Your flink server HTTP API
	c, err := api.New("127.0.0.1:8081")
	if err != nil {
		panic(err)
	}

	// get cluster config
	config, err := c.Config()
	if err != nil {
		panic(err)
	}
	fmt.Println(config)
}
```

More examples in [example](/example) dir.
### Cluster API

* shutdown cluster
* list config


### Jar File API

* upload jar file
* list jar files
* delete jar file
* plan jar file
* run jar file

### Job API

* job manager config
* job manager metrics
* list all jobs
* stop a job
* job overview
* job detail

### checkpoints

* get all checkpoints of a job
* stop a job with a savepoint

### TODO:

* vertices
* checkpoints/config
* /jobs/:jobid/checkpoints/details/:checkpointid
* /jobs/:jobid/config
* /jobs/:jobid/exceptions
* /jobs/:jobid/execution-result
* /jobs/:jobid/metrics
* /jobs/:jobid/plan
* /jobs/:jobid/rescaling
* /jobs/:jobid/rescaling/:triggerid
* overview
* /savepoint-disposal
* /taskmanagers

