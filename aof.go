package main

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"
)

type AOF struct {
	f *os.File
	w *bufio.Writer

	mu sync.Mutex
}

func NewAOF(path string) (*AOF, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &AOF{
		f: f,
		w: bufio.NewWriter(f),
	}, nil
}

func (a *AOF) WriteRaw(b []byte) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	_, err := a.w.Write(b)
	return err
}

func (a *AOF) FlushLoop() {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		slog.Info("[Writing to log]")
		a.mu.Lock()
		a.w.Flush() // flush bufio buffer to OS
		a.f.Sync()  // flush OS buffer to disk
		a.mu.Unlock()
	}
}

func (a *AOF) Replay(s *Store) error {
	_, err := a.f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(a.f)
	for {
		r := NewRequest()

		err := r.ParseReader(reader)
		if err != nil {
			break
		}
		cmd, err := r.ToCommand()
		if err != nil {
			// log it but keep going
			fmt.Println("warn: replay execute error", err)
			continue
		}

		_, err = cmd.Execute(s)
		if err != nil {
			// log it but keep going
			fmt.Println("warn: replay execute error", err)
			continue
		}
	}
	return nil
}
