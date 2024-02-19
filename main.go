package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"testModule/handlers"
	"time"
)

func main() {
	// -------------- First edition of the code before refactoring:
	// http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
	// log.Println("Initial Commit")
	// d, err := io.ReadAll(r.Body)

	// if err != nil {
	// 	http.Error(rw, "Error happened", http.StatusBadRequest)
	// 	return
	// }
	// fmt.Fprintf(rw, "Input %s", d)
	// })

	// http.ListenAndServe(":9090", serveMux)

	// -------------- Refactored version using hello handler
	l := log.New(os.Stdout, "products-api", log.LstdFlags)

	// Create the handlers
	productHandler := handlers.NewProducts(l)

	serveMux := http.NewServeMux()
	serveMux.Handle("/", productHandler)

	// Create a new server
	server := &http.Server{
		Addr:         ":9090",
		Handler:      serveMux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// Start the server
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	sig := <-signalChannel
	l.Println("Received Terminate, graceful shutdown", sig)

	timeout := time.Now().Add(30 * time.Second)
	timeoutContext, _ := context.WithDeadline(context.Background(), timeout)
	server.Shutdown(timeoutContext) // when we call Shutdown server will not accept any more requests, and let the in progress requests finish then it will shut down the server
}
