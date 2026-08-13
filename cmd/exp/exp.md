package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Server struct {
		Port      int
		Host      string
		StaticDir string
	}
}

func main() {
	file, err := os.Open("./config.json")
	if err != nil {
		panic(err)
	}
	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Dir Static: %s\n %s%d\n", config.Server.StaticDir, config.Server.Host, config.Server.Port)
}
