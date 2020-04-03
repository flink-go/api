package main

import (
	"os"

	"github.com/flink-go/api"
)

func main() {
	c, err := api.New(os.Getenv("FLINK_API"))
	if err != nil {
		panic(err)
	}

	// stop job test
	err = c.StopJob("8ea123d2bdc3064f36b92889e43803ee")
	if err != nil {
		panic(err)
	}
}
