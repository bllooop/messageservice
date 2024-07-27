package main

import (
	running "messageservice/pkg/server"

	_ "github.com/lib/pq"
)

func main() {
	running.Run()
}
