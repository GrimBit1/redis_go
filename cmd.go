package main

import "fmt"

type CommandType int

const (
	GET    CommandType = 0
	SET    CommandType = 1
	DELETE CommandType = 2
	HELLO  CommandType = 3
)

type Command struct {
	Type CommandType
	Args []string
}

func parseCommand(msg string) (Command, error) {
	t := msg[0]
	switch t {
	// Bulk string
	case '*':
		fmt.Println("bulk string")
	}
	return Command{}, nil
}
