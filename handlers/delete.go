package handlers

import (
	"net/http"
	"strconv"

	"github.com/dibakalantari/microservices-with-go/data"
	"github.com/gorilla/mux"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Returns a list of products
// responses:
// 201: noContent

// DeleteProduct deletes a product from the database
func(p *Products) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	p.l.Println("Handle DELETE product", id)

	err := data.DeleteProduct(id)

	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product Not Found", http.StatusInternalServerError)
		return
	}
}