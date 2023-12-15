package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/vvthai10/rabbitmq/rabbitmq"
)

// ErrConnectionClosed -.
var ErrConnectionClosed = errors.New("Client - RemoteCall - Connection closed")

const (
	_defaultWaitTime = 5 * time.Second
	_defaultAttempts = 10
	_defaultTimeout  = 2 * time.Second
)

// Message
type Message struct {
	Queue string
	Priority uint8
	ContentType string
	Body []byte
	ReplyTo string
	CorrelationID string
}

type pendingCall struct {
	done chan struct {}
	status string
	body []byte
}

type Client struct {
	conn *rabbitmq.Connection
	serverExchange string
	error chan error
	stop chan struct {}

	rw sync.RWMutex // make just working with one routine in time
	calls map[string]*pendingCall
	timeout time.Duration
}

func New(url, serverExchange, clientExchange string, opts ...Option) (*Client, error) {
	cfg := rabbitmq.Config {
		URL: url,
		WaitTime: _defaultWaitTime,
		Attempts: _defaultAttempts,
	}

	c := &Client{
		conn: rabbitmq.New(clientExchange, cfg),
		serverExchange: serverExchange,
		error: make(chan error),
		stop: make(chan struct{}),
		calls: make(map[string]*pendingCall),
		timeout: _defaultTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(c)
	}

	err := c.conn.AttemptConnect()
	if err != nil {
		return nil, fmt.Errorf("Client - c.conn.AttemptConnect: %w", err)
	}

	go c.consumer()

	return c, nil
}

func (c *Client) consumer() {
	for {
		select {
		case <-c.stop:
			return
		case d, opened := <-c.conn.Delivery:
			if !opened {
				c.reconnect()

				return
			}

			_ = d.Ack(false)

			// Write response to list call
			c.getCall(&d)
		}
	}
}

func (c *Client) getCall(d *amqp.Delivery) {
	c.rw.RLock()
	call, ok := c.calls[d.CorrelationId]
	c.rw.RUnlock()

	if !ok {
		return
	}

	call.status = d.Type
	call.body = d.Body
	close(call.done)
}

func (c *Client) publish(corrID, handler string, request interface{}) error {
	var (
		requestBody []byte
		err error
	)

	if request != nil {
		requestBody, err = json.Marshal(request)
		if err != nil {
			return err
		}
	}

	err = c.conn.Channel.Publish(c.serverExchange, "", false, false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrID,
			ReplyTo:       c.conn.ConsumerExchange,
			Type:          handler,
			Body:          requestBody,
		})
	if err != nil {
		return fmt.Errorf("c.Channel.Publish: %w", err)
	}

	return nil
}

func (c *Client) RemoteCall(handler string, request, response interface{}) error {
	select {
	case <-c.stop:
		time.Sleep(c.timeout)
		select {
		case <-c.stop:
			return ErrConnectionClosed
		default:
		}
	default:
	}

	corrID := uuid.New().String()
	
	err := c.publish(corrID, handler, request)
	if err != nil {
		return fmt.Errorf("Client - RemoteCall - c.publish: %w", err)
	}

	call := &pendingCall{done: make(chan struct{})}
	c.addCall(corrID, call)
	defer c.deleteCall(corrID)

	select {
	case <-time.After(c.timeout):
		return rabbitmq.ErrTimeout
	case <-call.done:
	}

	if call.status == rabbitmq.Success {
		err = json.Unmarshal(call.body, &response)
		if err != nil {
			return fmt.Errorf("Client - RemoteCall - json.Unmarshal: %w", err)
		}

		return nil
	}

	if call.status == rabbitmq.ErrBadHandler.Error() {
		return rabbitmq.ErrBadHandler
	}
	if call.status == rabbitmq.ErrInternalServer.Error() {
		return rabbitmq.ErrInternalServer
	}

	return nil
}

func (c *Client) addCall(corrID string, call *pendingCall) {
	c.rw.Lock()
	c.calls[corrID] = call
	c.rw.Unlock()
}

func (c *Client) deleteCall(corrID string) {
	c.rw.Lock()
	delete(c.calls, corrID)
	c.rw.Unlock()
}

func (c *Client) reconnect() {
	close(c.stop)

	err := c.conn.AttemptConnect()
	if err != nil {
		c.error <- err
		close(c.error)

		return
	}

	c.stop = make(chan struct{})

	go c.consumer()
}

// Notify -.
func (c *Client) Notify() <-chan error {
	return c.error
}

// Shutdown -.
func (c *Client) Shutdown() error {
	select {
	case <-c.error:
		return nil
	default:
	}

	close(c.stop)
	time.Sleep(c.timeout)

	err := c.conn.Connection.Close()
	if err != nil {
		return fmt.Errorf("rmq_rpc client - Client - Shutdown - c.Connection.Close: %w", err)
	}

	return nil
}