package main

import (
	"github.com/Dashboard/urlShortener/application"
)

func main() {
	a := application.Application{}
	if err := a.Init("config/config.yaml"); err != nil {
		panic(err)
	}
	a.Run()
}
