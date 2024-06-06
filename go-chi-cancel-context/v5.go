package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	ctx := context.Background()

	router := chi.NewRouter()
	// router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		go func() {
			select {
			case <-r.Context().Done():
				log.Println(r.Context().Err())
				return
			}
		}()

		time.Sleep(2 * time.Second)
		fmt.Println("It continues :(")
		w.Write([]byte("Hello, World!"))
	})

	server := http.Server{
		Addr:        ":4001",
		Handler:     router,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}
	log.Println("Server running on :4001")
	log.Fatal(server.ListenAndServe())
}
