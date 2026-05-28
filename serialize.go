package main

import "fmt"

func ToErrorString(s string) []byte {
	return []byte(string(Error) + s + SEP)
}

func ToSimpleString(s string) []byte {
	return []byte(string(SimpleString) + s + SEP)
}
func ToBulkString(s string) []byte {
	return []byte(fmt.Sprintf("%c%d\r\n%s\r\n", BulkString, len(s), s))
}

func NullBulkString() []byte {
	return []byte("$-1\r\n")
}

func ToInteger(i int) []byte {
	return []byte(fmt.Sprintf("%c%d%s", Integer, i, SEP))
}
