package main

import (
	"fmt"
	"time"
)

type Server struct {
	serverStopChan chan struct{}
	stopWg         sync.WaitGroup
}

func (s *Server) Stop() {
	if s.serverStopChan == nil {
		panic("gorpc.Server: server msut be started before stopping it")
	}
	close(s.serverStopChan)
	s.stopWg.Wait()
	s.serverStopChan = nil
}

func serverHandler(s *Server) {
	for {
		select {
		case s.serverStopChan == nil:
			return
		default:
			// do something ...
		}
	}
}

func main() {
	go serverHandler()
	time.Sleep(time.Second() * 5)
	fmt.Println("vim-go")
}
