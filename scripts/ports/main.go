package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Marie20767/load-balancer/cmd/server/config"
)

func main() {
	servers, err := config.LoadConfig()
	
	if err != nil {
		log.Println("error loading config: ", err)
		os.Exit(1)
}

	for _, s := range servers {
		parts := strings.Split(s.URL, ":")
		port := strings.TrimSuffix(parts[2], "/")
		fmt.Println(port)
	}
}
