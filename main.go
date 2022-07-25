package main

import (
	"github.com/bootjp/s3-redis/app"
	"log"
)

func main() {
	err := app.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
