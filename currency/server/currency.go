package server

import (
	"context"
	"io"
	"time"

	protos "github.com/dibakalantari/microservices-with-go/currency/currency"
	"github.com/hashicorp/go-hclog"
)

// Currency is a gRPC server it implements the methods defined by the CurrencyServer interface
type Currency struct {
	protos.UnimplementedCurrencyServer
	log           hclog.Logger
}

func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	return &protos.RateResponse{Rate: 0.5}, nil
}

func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {
	// handle client messages
	go func ()  {
		for {
			rr, err := src.Recv()
			if err == io.EOF {
				c.log.Info("CLient has closed connection")
				break;
			}
	
			if err != nil {
				c.log.Error("Unable to read from client", "error", err)
				break;
			}

			c.log.Info("Handle Client Request", "request", rr)
		}
	}()

	// handle server responses
	for {
		err := src.Send(&protos.RateResponse{Rate: 12.1})
		if err != nil {
			return err
		}

		time.Sleep(5 * time.Second)
	}
}