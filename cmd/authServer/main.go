package main

import (
	"fmt"
	"github.com/sha1sof/authRESTAPI/internal/config"
)

func main() {
	cfg := config.MustLoadConfig()
	fmt.Println(cfg)

	//TODO: init logger

	//TODO: init storage

	//TODO: init server
}
