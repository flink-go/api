# Flink Monitoring API Golang library

Detail doc: https://ci.apache.org/projects/flink/flink-docs-stable/monitoring/rest_api.html

```
package main

import (
	"fmt"

	"github.com/flink-go/api"
)

func main() {
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
