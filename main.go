package main

import (
	"cmp"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var HTTP_RESP_ID = cmp.Or(os.Getenv("SOURCE_ID"), os.Getenv("RAILWAY_REPLICA_ID"), "unknown")

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/status-code/200", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, http.StatusText(http.StatusOK))
	})

	mux.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, HTTP_RESP_ID)
	})

	port := cmp.Or(os.Getenv("PORT"), "8080")

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Printf("starting server on port %s\n", port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("error starting server: %v\n", err)
			os.Exit(1)
		}
	}()

	<-shutdown
	fmt.Println("server shutdown initiated")

	ctx, cancel := context.WithTimeout(context.Background(), (10 * time.Second))
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("error during server shutdown: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("server stopped")
}
