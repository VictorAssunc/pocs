package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	ctx := context.Background()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/5", func(w http.ResponseWriter, r *http.Request) {
		// Simulating hard and slow work
		time.Sleep(5 * time.Second)
		w.Write([]byte("hello world"))
	})
	r.Get("/10", func(w http.ResponseWriter, r *http.Request) {
		// Simulating hard and slow work
		time.Sleep(10 * time.Second)
		w.Write([]byte("hello world"))
	})
	r.Get("/15", func(w http.ResponseWriter, r *http.Request) {
		// Simulating hard and slow work
		time.Sleep(15 * time.Second)
		w.Write([]byte("hello world"))
	})

	server := http.Server{
		Addr:        ":8888",
		Handler:     r,
		BaseContext: func(net.Listener) context.Context { return ctx },
	}

	go func() {
		fmt.Println("Server is running on port 8888")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c

	fmt.Println("Shutting down server...")
	ctx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}
