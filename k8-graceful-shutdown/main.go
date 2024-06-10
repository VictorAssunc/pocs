package main

import (
	"context"
	"errors"
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
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(31 * time.Second)
		w.Write([]byte("hello world"))
	})
	router.Get("/medium", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.Write([]byte("hello world"))
	})
	router.Get("/fast", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	ctx := context.Background()
	server := http.Server{
		Addr:        ":8888",
		Handler:     router,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	go func() {
		log.Println("Server is running on port 8888")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c

	start := time.Now()
	log.Println("Shutting down server...")
	ctx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Printf("Server was shutdown - took %s\n", time.Since(start).String())
}
