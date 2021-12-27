package main

import (
	"fmt"

	post05 "github.com/virsel/go/postgres"
)

func main() {
	post05.Hostname = "localhost"
	fmt.Println(post05.Port)
	fmt.Println(post05.Hostname)
}
