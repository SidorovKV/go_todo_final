package main

import (
	"context"
	"errors"
	"go_todo_final/config"
	"go_todo_final/server"
	"go_todo_final/services/todolist"
	"go_todo_final/sqlitestorage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.New()

	storage, err := sqlitestorage.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer storage.Close()

	todoService := todolist.New(storage)

	ourSserver := server.New(cfg, todoService)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := ourSserver.Start(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Println("Server stopped.")
			} else {
				log.Fatal(err)
			}
		}
	}()

	log.Printf("Service started on port %s.\n", cfg.Port)

	<-quit

	if err := ourSserver.Stop(context.Background()); err != nil {
		log.Fatal(err)
	}

	log.Println("Service stopped.")
}
