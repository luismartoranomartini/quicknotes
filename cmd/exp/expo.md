package main

import (
	"fmt"
	"os"

	// "github.com/joho/godotenv"
	"flag"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println()
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	fmt.Printf("%s:%s", host, port)
}
