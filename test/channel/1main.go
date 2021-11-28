package main

import (
	"fmt"
	"sync"
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
		case <-s.serverStopChan:
			return
		default:
			// do something ...
		}
	}
}

func main() {
	var s Server
	go serverHandler(&s)
	time.Sleep(time.Second * 5)
	s.Stop()
	fmt.Println("vim-go")
	select {}
}
