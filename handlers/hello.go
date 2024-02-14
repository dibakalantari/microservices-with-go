package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Hello struct{
	l *log.Logger
}

// constructor function
func NewHello(l *log.Logger) *Hello{
	return &Hello{l}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Initial Request")
	d, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(rw, "Error happened", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(rw, "Input %s", d)
}