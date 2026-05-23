package main

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"time"
)

const DefaultListenAddr = ":6379"

type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	ln    net.Listener
	store Store
}

func NewServer(cfg Config) *Server {
	if cfg.ListenAddr == "" {
		cfg.ListenAddr = DefaultListenAddr
	}
	return &Server{
		Config: cfg,
		store:  NewStore(),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	s.ln = ln
	return s.acceptLoop()
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		fmt.Println("conn", conn, time.Now())
		if err != nil {
			slog.Error("accept err", "error", err)
			return err
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) error {
	buf := make([]byte, 2048)
	read := 0
	r := NewRequest()
	defer conn.Close()
	for {
		n, err := conn.Read(buf[read:])
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				fmt.Println("conn closed by client")
			} else {
				return err
			}
		}

		read += n

		re, err := r.Parse(buf[:read])
		if err != nil {
			return err
		}

		if re != 0 {
			copy(buf, buf[re:read])
			read -= re
		}

		if r.state == Done {
			cmd, err := r.ToCommand()
			if err != nil {
				return err
			}
			fmt.Println("cmd", cmd)
			switch cmd.Type {
			case SET:
				err := s.store.Set(cmd.Args[0], cmd.Args[1])
				if err != nil {
					conn.Write([]byte("-" + err.Error()))
					continue
				}
				fmt.Println("done", time.Now())
			case GET:
				val, ok := s.store.Get(cmd.Args[0])
				if !ok {
					conn.Write([]byte("-not found"))
				} else {
					conn.Write([]byte("+" + val + SEP))
				}
				fmt.Println("done", time.Now())

				continue
			}
			conn.Write([]byte("+OK\r\n"))
			r = NewRequest()
			read = 0
		}
	}
}

func main() {
	s := NewServer(Config{})

	log.Fatal(s.Start())

}
