package main

import (
	"fmt"
	"testing"

	"github.com/dibakalantari/microservices-with-go/sdk/client"
	"github.com/dibakalantari/microservices-with-go/sdk/client/products"
)

func TestOurClient(t *testing.T) {
	cfg := client.DefaultTransportConfig().WithHost("localhost:9090")
	c := client.NewHTTPClientWithConfig(nil,cfg)

	params := products.NewListProductsParams()
	prod, err := c.Products.ListProducts(params)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%#v", prod.GetPayload()[0])
}