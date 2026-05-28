package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"
)

type CommandType string

const (
	GET    CommandType = "get"
	SET    CommandType = "set"
	DELETE CommandType = "delete"
	HELLO  CommandType = "hello"
	PING   CommandType = "ping"
	ECHO   CommandType = "echo"
)

type Command struct {
	Type    CommandType
	Writing bool
	Args    []string
}

func (c *Command) ToResp() ([]byte, error) {
	resp := bytes.Buffer{}

	_, err := fmt.Fprintf(&resp, "*%d\r\n", len(c.Args)+1)
	if err != nil {
		return nil, err
	}
	_, err = fmt.Fprintf(&resp, "$%d\r\n%s\r\n", len(c.Type), c.Type)
	if err != nil {
		return nil, err
	}

	for _, v := range c.Args {
		_, err = fmt.Fprintf(&resp, "$%d\r\n%s\r\n", len(v), v)
		if err != nil {
			return nil, err
		}
	}

	return resp.Bytes(), nil
}

func (cmd *Command) Execute(s *Store) ([]byte, error) {
	switch cmd.Type {
	case SET:
		if len(cmd.Args) < 2 {
			return nil, errors.New("wrong number of arguments for 'set' command")
		}
		ttl := time.Duration(0)
		var err error
		// Timestamp command
		if len(cmd.Args) == 4 {
			expiryType := cmd.Args[2]
			switch strings.ToLower(expiryType) {
			case "ex":
				ttl, err = time.ParseDuration(cmd.Args[3] + "s")
				if err != nil {
					return nil, err
				}

			case "px":
				ttl, err = time.ParseDuration(cmd.Args[3] + "ms")
				if err != nil {
					return nil, err
				}

			default:
				return nil, (errors.New("syntax error"))
			}

		}
		err = s.Set(cmd.Args[0], cmd.Args[1], time.Duration(ttl))
		if err != nil {
			return nil, err
		}
		fmt.Println("done", time.Now())
		return ToSimpleString("OK"), nil

	case GET:
		val, ok := s.Get(cmd.Args[0])
		if !ok {
			return NullBulkString(), nil
		} else {
			return ToSimpleString(val), nil
		}

	case HELLO:
		return []byte("*14\r\n" +
			"$6\r\nserver\r\n$5\r\nredis\r\n" +
			"$7\r\nversion\r\n$5\r\n6.0.0\r\n" +
			"$5\r\nproto\r\n$1\r\n2\r\n" +
			"$2\r\nid\r\n$1\r\n1\r\n" +
			"$4\r\nmode\r\n$10\r\nstandalone\r\n" +
			"$4\r\nrole\r\n$6\r\nmaster\r\n" +
			"$7\r\nmodules\r\n*0\r\n"), nil
	case PING:
		return ToSimpleString("PONG"), nil
	case ECHO:
		if len(cmd.Args) < 1 {
			return nil, errors.New("wrong number of arguments")
		}
		val := cmd.Args[0]
		return ToBulkString(val), nil
	case DELETE:
		if len(cmd.Args) < 1 {
			return nil, errors.New("wrong number of arguments")
		}
		err := s.Delete(cmd.Args[0])
		if err != nil {
			return nil, err
		}
		return ToSimpleString("OK"), nil
	default:
		return ToSimpleString("OK"), nil
	}
}
