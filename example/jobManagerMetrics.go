package main

import (
	"fmt"
	"os"

	"github.com/flink-go/api"
)

func main() {
	c, err := api.New(os.Getenv("FILNK_API"))
	if err != nil {
		panic(err)
	}

	// job manager metrics test
	config, err := c.JobManagerMetrics()
	if err != nil {
		panic(err)
	}
	fmt.Println(config)
}
