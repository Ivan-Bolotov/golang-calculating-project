package main

import (
	"github.com/Ivan-Bolotov/golang-calculating-project/internal/servers"
	"time"
)

func main() {
	go servers.StartNewHttpStorageServer()
	time.Sleep(time.Millisecond * 100)
	servers.StartNewHttpComputingServer(8081)
}
