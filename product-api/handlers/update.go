package handlers

import (
	"net/http"
	"strconv"
	"github.com/dibakalantari/microservices-with-go/product-api/data"

	"github.com/gorilla/mux"
)

func (p* Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert", http.StatusBadRequest)
	}

	p.l.Println("Handle PUT Product")
	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}