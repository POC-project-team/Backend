package app

import (
	s "backend/internal/controller/rest/server"
	log "github.com/sirupsen/logrus"
)

func Run() {
	server := s.NewServer()
	log.Info("The server is up and running at ", server.Addr, "\n")

	// signal handler for correct shutdown
	done := make(chan bool)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Info(err.Error())
		}
		done <- true
	}()

	server.WaitShutdown()

	<-done
}
