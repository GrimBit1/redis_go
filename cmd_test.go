package main

import "testing"

func TestParseCmd(t *testing.T) {
	cmd := `*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$7\r\nmyvalue\r\n`

	_, err := parseCommand((cmd))
	if err != nil {
		t.Error(err)
	}
}
