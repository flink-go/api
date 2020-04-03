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

	opts := api.RunOpts{
		JarID: "8c0c2226-b532-4d9b-b698-8aa649694bb9_test.jar",
	}
	// run test
	resp, err := c.RunJar(opts)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
