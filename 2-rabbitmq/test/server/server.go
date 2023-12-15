package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
	"github.com/vvthai10/rabbitmq/rabbitmq/server"
)

// List func handle logic of rabbitmq
func NewRouter() map[string]server.CallHandler {
	routes := map[string]server.CallHandler{
		"default": defaultHandle() ,
		"secret": secretHandle(),
	}
	return routes
}

func defaultHandle() server.CallHandler{
	return func(d *amqp.Delivery) (interface{}, error) {
		// Doing...
		fmt.Println("Server - defaultHandle: ", d.Body)
		return "Message client is handle by default", nil
	}
}

func secretHandle() server.CallHandler{
	return func(d *amqp.Delivery) (interface{}, error) {
		// Doing...
		fmt.Println("Server - defaultHandle: ", d.Body)
		return "Message client is handle by secret", nil
	}
}

const (
	URL = "amqp://guest:guest@localhost:5672/"
	serverExchange = "server_exchange"
	clientExchange = "clientExchange"
)

func main() {
	router := NewRouter()
	server, err := server.New(URL, serverExchange, router)
	if err != nil {
		fmt.Println("Run - server.New: ", err)
	}

	err = <-server.Notify()
	if err != nil {
		fmt.Println("Run - server.Notify: ", err)
	}
	

	err = server.Shutdown()
	if err != nil {
		log.Fatal("Run - server.Shutdown: %w", err)
	}
}