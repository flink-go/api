package main

import (
	"github.com/flink-go/api"
)

func main() {
	c, err := api.New("127.0.0.1:8081")
	if err != nil {
		panic(err)
	}

	// shutdown test
	if err := c.Shutdown(); err != nil {
		panic(err)
	}
}
