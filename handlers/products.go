package handlers

import (
	"log"
	"net/http"
	"strconv"
	"testModule/data"

	"github.com/gorilla/mux"
)

// Products is a http.Handler
type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	listOfProducts := data.GetProducts()
	// convert to json
	err := listOfProducts.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
} 

func (p* Products) AddProduct(rw http.ResponseWriter, r *http.Request){
	p.l.Println("Handle POST Product")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)

	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	data.AddProduct(prod)
}

func (p* Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert", http.StatusBadRequest)
	}

	p.l.Println("Handle PUT Product")

	prod := &data.Product{}
	err = prod.FromJSON(r.Body)

	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

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
