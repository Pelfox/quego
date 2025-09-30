package main

// Example application demonstrating the usage of the quego server.

import (
	"fmt"
	"time"

	"github.com/Pelfox/quego"
	"github.com/Pelfox/quego/models"
)

func main() {
	server, err := quego.NewServer()
	if err != nil {
		panic(err)
	}

	server.RegisterFunction("hello-world", func(trigger *models.Trigger) error {
		fmt.Println("Function triggered!")
		time.Sleep(10 * time.Second)
		fmt.Println("Function completed!")
		return nil
	})

	if err := server.Start(":8080"); err != nil {
		panic(err)
	}
}
