package main

import (
	"bytes"
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

	data := fmt.Sprintf("*3\r\n$3\r\nSET\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n", len(key), key, len(value), value)
	fmt.Println("data", data)

	_, err := c.conn.Write([]byte(data))
	if err != nil {
		return err
	}
	buf := make([]byte, 1024)
	c.conn.Read(buf)
	fmt.Println("buf", strings.TrimSpace(string(buf)))
	return nil
}
func (c *Client) Get(ctx context.Context, key string) (string, error) {

	data := fmt.Sprintf("*2\r\n$3\r\nGET\r\n$%d\r\n%s\r\n", len(key), key)
	fmt.Println("data", data)

	_, err := c.conn.Write([]byte(data))
	if err != nil {
		return "", err
	}
	buf := make([]byte, 1024)
	val := ""
	for {
		n, err := c.conn.Read(buf)
		if err != nil {
			return "", err
		}
		currentData := buf[:n]
		fmt.Println("buf", currentData)
		i := bytes.Index(currentData, []byte(SEP))
		if i != -1 {
			val = string(currentData[:i])
			break
		} else {
			continue
		}
	}

	return val[1:], nil
}
func (c *Client) Close() error {
	return c.conn.Close()
}
