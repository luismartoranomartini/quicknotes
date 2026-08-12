package main

import (
	"flag"
	"fmt"
)

func main() {
	var port string
	var verbose bool
	var valor int

	// port := flag.StringVar(&port, "port", "7000", "Server port")
	flag.StringVar(&port, "port", "7000", "Server port")
	flag.StringVar(&port, "p", "7000", "Server port")
	flag.IntVar(&valor, "valor", 0, "some value")

	flag.BoolVar(&verbose, "v", false, "Verbose mode")

	flag.Parse()

	if verbose {

		fmt.Println("Server is running on port", port)
		fmt.Println("Valor", valor)
	} else {
		fmt.Println(port)
	}

	// port := fmt.Println("Server is running on port", *port)
}
