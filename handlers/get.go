package handlers

import (
	"net/http"
	"github.com/dibakalantari/microservices-with-go/data"
)

// swagger:route method url tag name 
// What does it do?
// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
// 200: productsResponse
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	listOfProducts := data.GetProducts()
	// convert to json
	err := listOfProducts.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
} 