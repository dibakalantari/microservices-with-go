package main

import (
	protos "currency/currency"
	"currency/server"
	"net"
	"os"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()

	gs := grpc.NewServer()
	c := &server.Currency{}

	protos.RegisterCurrencyServer(gs, c) 

	reflection.Register(gs) 

	l, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Error("Unable to listen","error" , err)
		os.Exit(1)
	}

	gs.Serve(l)
}
