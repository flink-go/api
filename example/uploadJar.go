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

	// upload test
	resp, err := c.UploadJar("./testdata/test.jar")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
