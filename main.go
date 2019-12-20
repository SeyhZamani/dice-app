package main

import (
	"context"
	"github.com/SeyhZamani/dice-game-app/app/handler"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	router := chi.NewRouter()

	router.Get("/match", handler.PostMatchHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Println("Server Starting...")
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	waitForShutdown(server)
}

func waitForShutdown(server *http.Server) {
	interruptChan := make(chan os.Signal, 1)

	// interrupt signal sent from terminal
	signal.Notify(interruptChan, os.Interrupt)
	// sigterm signal sent from kubernetes
	signal.Notify(interruptChan, syscall.SIGINT)
	signal.Notify(interruptChan, syscall.SIGTERM)

	<-interruptChan

	// crete deadline
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("HTTP Server Shutdown Error: %v", err)
	}

	log.Println("Shutting Down...")
	os.Exit(0)
}
