package handlers

import (
	"net/http"
	"github.com/dibakalantari/microservices-with-go/data"
)

func (p* Products) AddProduct(rw http.ResponseWriter, r *http.Request){
	p.l.Println("Handle POST Product")

	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	data.AddProduct(prod)
}