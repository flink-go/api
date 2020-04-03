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

	// job manager config test
	config, err := c.JobManagerConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(config)
}
