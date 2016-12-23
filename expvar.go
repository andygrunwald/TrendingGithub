package main

import (
	_ "expvar"
	"fmt"
	"log"
	"net"
	"net/http"
)

// StartExpvarServer will start a small tcp server for the expvar package.
// This server is only available via localhost on localhost:port/debug/vars
func StartExpvarServer(port int) error {
	sock, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return err
	}

	go func() {
		log.Printf("Expvar initialisation success: HTTP now available at http://localhost:%d/debug/vars", port)
		http.Serve(sock, nil)
	}()

	return nil
}
