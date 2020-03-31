package main

import (
	"github.com/flink-go/api"
)

func main() {
	c, err := api.New("127.0.0.1:8081")
	if err != nil {
		panic(err)
	}

	// delete test
	err = c.DeleteJar("efb8367e-aa0d-4ceb-957f-0e8f46ed4b10_test.jar")
	if err != nil {
		panic(err)
	}
}
