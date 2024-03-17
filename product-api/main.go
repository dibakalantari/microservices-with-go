package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	protos "github.com/dibakalantari/microservices-with-go/currency/currency"
	gohandlers "github.com/gorilla/handlers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/dibakalantari/microservices-with-go/product-api/handlers"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "products-api", log.LstdFlags)

	conn, err := grpc.Dial("localhost:9092", grpc.WithTransportCredentials(insecure.NewCredentials())) // only for non-production environments, on production definitely use https
	if err != nil {
		panic(err)
	}	
	defer conn.Close()	
	// Create Client
	cc := protos.NewCurrencyClient(conn)

	// Create the handlers
	productHandler := handlers.NewProducts(l,cc)

	serveMux := mux.NewRouter()
	
	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", productHandler.ListAll)

	getRouter.HandleFunc("/products/{id:[0-9]+}", productHandler.ListSingle)

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
