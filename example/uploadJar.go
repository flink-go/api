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

	// upload test
	resp, err := c.UploadJar("./testdata/test.jar")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
