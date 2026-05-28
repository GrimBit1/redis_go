package main

import (
	"bufio"
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
	aof   *AOF
}

func NewServer(cfg Config) (*Server, error) {
	if cfg.ListenAddr == "" {
		cfg.ListenAddr = DefaultListenAddr
	}
	aof, err := NewAOF("./tmp.log")
	if err != nil {
		return nil, err
	}
	return &Server{
		Config: cfg,
		store:  NewStore(),
		aof:    aof,
	}, nil
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	s.ln = ln
	go s.aof.FlushLoop()
	s.store.startExpirySweep()
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

	reader := bufio.NewReader(conn)

	defer conn.Close()
	for {
		r := NewRequest()

		err := r.ParseReader(reader)
		if err != nil {
			return err
		}
		cmd, err := r.ToCommand()
		if err != nil {
			return err
		}
		fmt.Println("cmd", cmd)

		res, err := cmd.Execute(&s.store)
		if err != nil {
			conn.Write(ToErrorString("ERR " + err.Error()))
		} else {
			conn.Write(res)
			if cmd.Writing {
				resp, err := cmd.ToResp()
				if err != nil {
					return err
				}
				err = s.aof.WriteRaw(resp)
				if err != nil {
					return err
				}
			}
		}
	}
}

func main() {
	s, err := NewServer(Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := s.aof.Replay(&s.store); err != nil {
		log.Fatal(err)
	}

	log.Fatal(s.Start())

}
