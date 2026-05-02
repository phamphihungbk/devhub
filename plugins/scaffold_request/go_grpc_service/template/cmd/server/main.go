package main

import (
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":[[PORT]]")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[[SERVICE_NAME]] gRPC service listening on %s", listener.Addr().String())
	select {}
}
