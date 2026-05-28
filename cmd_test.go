package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestParseCmd(t *testing.T) {
	input := "*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$7\r\nmyvalue\r\n"

	// 1. Wrap the string in strings.NewReader to get an io.Reader
	stringReader := strings.NewReader(input)

	// 2. Pass that reader to bufio.NewReader
	reader := bufio.NewReader(stringReader)
	request := NewRequest()
	err := request.ParseReader(reader)
	if err != nil {
		t.Error(err)
	}

	cmd, err := request.ToCommand()
	if err != nil {
		t.Error(err)
	}

	if cmd.Type != SET {
		t.Error("wrong type")
	}
}
