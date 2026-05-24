package main

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type RespType rune

const SEP = "\r\n"

const (
	Array        RespType = '*'
	BulkString   RespType = '$'
	SimpleString RespType = '+'
	Integer      RespType = ':'
	Error        RespType = '-'
)

type ReqState string

const (
	Init   ReqState = "init"
	Type   ReqState = "got type"
	Len    ReqState = "got len"
	Header ReqState = "reading header"
	Data   ReqState = "reading data"
	CMD    ReqState = "command complete"
	Done   ReqState = "done"
)

type Elem struct {
	Type  RespType
	Len   int
	Array []Elem
	Value []byte
}

type Request struct {
	state ReqState
	Elem
}

func NewRequest() *Request {
	return &Request{
		state: Init,
	}
}

func (r *Request) Parse(buf []byte) (int, error) {
	if len(buf) == 0 {
		return 0, nil
	}
	read := 0

	for r.state != Done {
		data := buf[read:]
		if len(data) < 2 {
			break
		}

		i := len(data)
		if r.state != Data {
			i = bytes.Index(data, []byte("\r\n"))

			if i == -1 {
				return read, nil
			}
			if i == 0 {
				read += len(SEP)
				if r.state == CMD {
					r.state = Done
				}
				return read, nil
			}
		}
		currentData := data[:i]
		fmt.Println(i, currentData)

		re, err := r.ParseMessage(currentData)
		if err != nil {
			return 0, err
		}
		read += re

	}

	return read, nil
}

func (r *Request) ParseMessage(data []byte) (int, error) {
	read := 0
	defer fmt.Println(string(data))

	if r.state == Init {
		r.Type = RespType(data[0])
		r.state = Type

		s := data[1:]
		lenInt, err := strconv.Atoi(string(s))
		if err != nil {
			return 0, err
		}
		r.Len = lenInt
		r.state = Header

		if r.Type == Array {
			r.Array = make([]Elem, 0, lenInt)
		}
		read += len(data) + len(SEP)
		return read, nil
	}
	if r.state == Header || r.state == Data {
		if r.Type == Array {
			switch data[0] {
			case byte(BulkString):
				fallthrough

			case byte(SimpleString):
				fallthrough

			case byte(Integer):
				fallthrough

			case byte(Array):
				elem := Elem{Type: RespType(data[0])}
				lenStr := data[1:]
				lenInt, err := strconv.Atoi(string(lenStr))
				if err != nil {
					return 0, err
				}
				elem.Len = lenInt
				if elem.Type == BulkString {
					elem.Value = make([]byte, 0, elem.Len)
				}
				r.Array = append(r.Array, elem)
				r.state = Data
				read += len(data) + len(SEP)
			default:
				if len(r.Array) == 0 {
					return 0, errors.New("not valid data")
				}
				read += min(len(data), r.Array[len(r.Array)-1].Len-len(r.Array[len(r.Array)-1].Value))
				r.Array[len(r.Array)-1].Value = append(r.Array[len(r.Array)-1].Value, data[:min(len(data), r.Array[len(r.Array)-1].Len-len(r.Array[len(r.Array)-1].Value))]...)

				if r.Array[len(r.Array)-1].Len == len(r.Array[len(r.Array)-1].Value) {
					if len(r.Array) == r.Len {
						r.state = CMD
					} else {
						r.state = Header
					}
					if len(data) != read && string(data[read:min(len(data), read+len(SEP))]) == SEP {
						read += len(SEP)
						if r.state == CMD {
							r.state = Done
						}
					}
					return read, nil
				}
				return read, nil
			}

		}
		return len(data) + len(SEP), nil
	}

	return read, nil
}

func (r *Request) ToCommand() (Command, error) {
	if r.Type != Array && len(r.Array) != 0 {
		return Command{}, errors.New("not valid request to convert to cmd")
	}
	cmd := Command{}
	switch strings.ToLower(string(r.Array[0].Value)) {
	case "set":
		cmd.Type = SET
	case "get":
		cmd.Type = GET
	case "delete":
		cmd.Type = DELETE
	case "hello":
		cmd.Type = HELLO
	default:
		cmd.Type = -1
	}

	cmd.Args = make([]string, 0, r.Len)

	for _, v := range r.Array[1:] {
		cmd.Args = append(cmd.Args, string(v.Value))
	}

	return cmd, nil

}
