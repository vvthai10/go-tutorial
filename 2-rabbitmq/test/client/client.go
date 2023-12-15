package main

import (
	"fmt"
	"log"
	"time"

	"github.com/vvthai10/rabbitmq/rabbitmq/client"
)

const (
	URL            = "amqp://guest:guest@localhost:5672/"
	serverExchange = "server_exchange"
	clientExchange = "client_exchange"
	requests = 10
)

func main() {
	client, err := client.New(URL, serverExchange, clientExchange)
	if err != nil {
		log.Fatal("Client - client.New: %w", err)
	}
	defer func() {
		err = client.Shutdown()
		if err != nil {
			log.Fatal("Client - client.Shutdown: %w", err)
		}
	}()

	for i := 0; i < requests; i++ {
		var res string
		err = client.RemoteCall("default", nil, &res)
		if err != nil {
			log.Fatal("Client -  client.RemoteCall: %w", err)
		}

		fmt.Println("Response from server: ", res)
		time.Sleep(6*time.Second)
	}
}