package main

import (
	"flag"
	"fmt"
	"log"

	clnttcp "swe-challenge-tcp/client"
	srvtcp "swe-challenge-tcp/server"

	"github.com/pkg/errors"
)

func main() {
	fmt.Println("Starting...")
	connect := flag.String("connect", "", "IP address of process to join. If empty, go into the listen mode.")
	flag.Parse()
	// Si la bandera de conexión está configurada, ingrese al modo cliente

	// Connect the client to the server
	if *connect != "" {
		err := clnttcp.Client(*connect)
		if err != nil {
			log.Println("Error:", errors.WithStack(err))
		}
		log.Println("Client done.")
		return
	}
	// Connect to the server
	err := srvtcp.Server()
	if err != nil {
		log.Println("Error:", errors.WithStack(err))
	}
	log.Println("Server done.")
}
