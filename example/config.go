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

	// config test
	config, err := c.Config()
	if err != nil {
		panic(err)
	}
	fmt.Println(config)
}
