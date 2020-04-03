package main

import (
	"fmt"
	"os"

	"github.com/flink-go/api"
)

func main() {
	c, err := api.New(os.Getenv("FLINK_API"))
	if err != nil {
		panic(err)
	}

	// jobs overview test
	jobs, err := c.JobsOverview()
	if err != nil {
		panic(err)
	}
	fmt.Println(jobs)
}
