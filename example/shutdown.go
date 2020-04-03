package main

import (
	"os"

	"github.com/flink-go/api"
)

func main() {
	c, err := api.New(os.Getenv("FILNK_API"))
	if err != nil {
		panic(err)
	}

	// shutdown test
	if err := c.Shutdown(); err != nil {
		panic(err)
	}
}
