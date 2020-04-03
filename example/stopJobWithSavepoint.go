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

	v, err := c.StopJobWithSavepoint("2bd452ba193d1575a4acc9ed09f896ea", "test", false)
	if err != nil {
		panic(err)
	}
	fmt.Println(v)
}
