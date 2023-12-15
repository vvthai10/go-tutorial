package server

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"github.com/vvthai10/rabbitmq/rabbitmq"
)

const (
	_defaultWaitTime = 5 * time.Second
	_defaultAttempts = 10
	_defaultTimeout  = 2 * time.Second
)

// CallHandler
// Format of function to handle logic when have message to queues
type CallHandler func(*amqp.Delivery) (interface{}, error)

// Server
type Server struct {
	conn *rabbitmq.Connection
	error chan error
	stop chan struct{}
	router map[string]CallHandler

	timeout time.Duration
	// Add logger...
}

// New
func New(url, serverExchange string, router map[string]CallHandler, opts ...Option) (*Server, error) {
	cfg := rabbitmq.Config {
		URL: url,
		WaitTime: _defaultWaitTime,
		Attempts: _defaultAttempts,
	}

	s := &Server{
		conn: rabbitmq.New(serverExchange, cfg),
		error: make(chan error),
		stop: make(chan struct{}),
		router: router,
		timeout: _defaultTimeout,
		// Add logger...
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	err := s.conn.AttemptConnect()
	if err != nil {
		return nil, fmt.Errorf("New server - AttemptConnect: %w", err)
	}

	// Run to wait message
	go s.consumer()

	return s, nil
}

func (s *Server) consumer() {
	for {
		select {
		case <-s.stop:
			return
		case d, opened := <-s.conn.Delivery:
			if !opened {
				s.reconnect()
				return
			}
			_ = d.Ack(false) // confirm remove message in queue

			s.serverCall(&d)
		}
	}
}

// Handle logic when have message
func (s *Server) serverCall(d *amqp.Delivery) {
	callHandler, ok := s.router[d.Type]
	if !ok {
		s.publish(d, nil, rabbitmq.ErrBadHandler.Error())
		return
	}

	response, err := callHandler(d)
	if err != nil {
		s.publish(d, nil, rabbitmq.ErrInternalServer.Error())
		fmt.Println("Server - serverCall - callHandler: %w",err)
		return
	}

	body, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Server - serverCall - json.Marshal: %w",err)
	}
	s.publish(d, body, rabbitmq.Success)	
}

// Send message to publisher/client
func (s *Server) publish(d *amqp.Delivery, body []byte, status string) {
	err := s.conn.Channel.Publish(d.ReplyTo, "", false, false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: d.CorrelationId,
			Type:          status,
			Body:          body,
		})
	if err != nil {
		fmt.Println("Server - publish - s.conn.Channel.Publish: %w", err)
	}
}

func (s *Server) reconnect() {
	close(s.stop)

	err := s.conn.AttemptConnect()
	if err != nil {
		s.error <- err
		close(s.error)

		return
	}

	s.stop = make(chan struct{})

	go s.consumer()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.error
}

// Shutdown -.
func (s *Server) Shutdown() error {
	select {
	case <-s.error:
		return nil
	default:
	}

	close(s.stop)
	time.Sleep(s.timeout)

	err := s.conn.Connection.Close()
	if err != nil {
		return fmt.Errorf("Server - Shutdown - s.Connection.Close: %w", err)
	}

	return nil
}

