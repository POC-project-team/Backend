package app

import (
	s "backend/internal/controller/rest/server"
	db "backend/internal/repository/sqlite"
	log "github.com/sirupsen/logrus"
)

func Run() {
	database := db.NewSQLDataBase()
	server := s.NewServer(database)
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
