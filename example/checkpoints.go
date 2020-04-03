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

	// Checkpoints test
	v, err := c.Checkpoints("8355c4efb63ddd4ea26a7829156a4c58")
	if err != nil {
		panic(err)
	}
	fmt.Println(v)
}
