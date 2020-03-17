package main

import (
	"gotools/log"
)

func main() {
	log.Init("./demo_error.log", "./demo_info.log")
	log.Info.Println("Info log...")
	log.Warning.Printf("Warning log...")
	log.Error.Println("Error log...")
}