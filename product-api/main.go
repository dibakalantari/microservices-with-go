package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dibakalantari/microservices-with-go/handlers"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	gohandlers "github.com/gorilla/handlers"
)

func main() {
	l := log.New(os.Stdout, "products-api", log.LstdFlags)

	// Create the handlers
	productHandler := handlers.NewProducts(l)

	serveMux := mux.NewRouter()
	
	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", productHandler.ListAll)

	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products/{id:[0-9]+}", productHandler.UpdateProducts)
	putRouter.Use(productHandler.MiddlewareValidateProduct)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", productHandler.AddProduct)
	postRouter.Use(productHandler.MiddlewareValidateProduct)

	deleteRouter := serveMux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/products/{id:[0-9]+}", productHandler.DeleteProduct)
	
	// Redoc route
	ops := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(ops, nil)
	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// CORS Handlers
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	// filename regex: {filename:[a-zA-Z]+\\.[a-z]{3}}
	ph := serveMux.Get(http.MethodPost).Subrouter()
	ph.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", )
 
	// Create a new server
	server := &http.Server{
		Addr:         ":9090",
		Handler:      ch(serveMux),
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