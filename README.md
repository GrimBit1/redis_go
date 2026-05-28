# redis_go

A minimal Redis-like server implemented in Go for learning systems programming and the Redis protocol.

This project implements a small in-memory key-value store with a RESP2 protocol parser and a TCP server using goroutines.

Features

- RESP2 protocol parser
- TCP server handling connections with goroutines
- Pipelining support for multiple commands in one request
- TTL/expiry with background sweeper
- Append-only log (AOF) for persistence
- Thread-safe key-value store (sync.RWMutex)
- Compatible with `redis-cli` and `go-redis`

Quick start

Prerequisites

- Go 1.25+ installed

Run the server

```bash
# from repository root
go run ./main.go
```

The server listens on `:6379` by default.

Use `redis-cli` or `go-redis` to interact:

```bash
# using redis-cli
redis-cli -p 6379 SET foo bar
redis-cli -p 6379 GET foo

# using go-redis (example client available in ./client)
go run ./client
```

Development

Run tests:

```bash
go test ./...
```

Build:

```bash
go build .
```

Implementation notes

- The RESP parser started as a custom stateful parser handling partial TCP frames from `net.Conn`; switching to `bufio.Reader` simplified parsing by allowing blocking reads and a linear parsing flow.
- `store.go` keeps entries in memory with optional TTL and uses `sync.RWMutex` for concurrency safety. Expired entries are treated as missing (the background sweeper can remove them periodically).
- `cmd.go` implements core Redis commands (GET, SET, DEL, PING, ECHO, HELLO), including TTL parsing for `EX` and `PX` variants.

Still Contributing
