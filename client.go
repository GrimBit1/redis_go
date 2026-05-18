package main

import (
	"context"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	addr string
	conn net.Conn
}

var dialer = &net.Dialer{}

func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		addr: addr,
		conn: conn,
	}, nil
}

func (c *Client) Set(ctx context.Context, key, value string) error {
	conn, err := dialer.DialContext(ctx, "tcp", c.addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	for i := range 10 {

		data := fmt.Sprintf("*3\r\n$3\r\nSET\r\n$%d\r\n%s_%d\r\n$%d\r\n%s_%d\r\n", len(key)+1, key, i, len(value)+1, value, i)
		fmt.Println("data", data)

		_, err = conn.Write([]byte(data))
		if err != nil {
			return err
		}
		buf := make([]byte, 1024)
		conn.Read(buf)
		fmt.Println("buf", strings.TrimSpace(string(buf)))
	}
	return nil
}
func (c *Client) Get(ctx context.Context, key string) error {
	conn, err := dialer.DialContext(ctx, "tcp", c.addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	for i := range 10 {

		data := fmt.Sprintf("*2\r\n$3\r\nGET\r\n$%d\r\n%s_%d\r\n", len(key)+1, key, i)
		fmt.Println("data", data)

		_, err = conn.Write([]byte(data))
		if err != nil {
			return err
		}
		buf := make([]byte, 1024)
		conn.Read(buf)
		fmt.Println("buf", strings.TrimSpace(string(buf)))
	}
	return nil
}
