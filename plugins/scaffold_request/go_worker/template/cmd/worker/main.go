package main

import (
	"log"
	"time"
)

func main() {
	log.Println("[[SERVICE_NAME]] worker started")
	for {
		log.Println("processing background work")
		time.Sleep(30 * time.Second)
	}
}
