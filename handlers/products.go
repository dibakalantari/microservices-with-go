package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"testModule/data"
)

// Products is a http.Handler
type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if(r.Method == http.MethodGet) {
		p.getProducts(rw, r)
		return
	}

	if(r.Method == http.MethodPost) {
		p.addProduct(rw, r)
		return
	}

	if(r.Method == http.MethodPut) {
		// expect the id in the URI
		regex := regexp.MustCompile(`/([0-9]+)`) // we are looking for id, one or more numbers after slash 
		regexGroup := regex.FindAllStringSubmatch(r.URL.Path, -1)

		if len(regexGroup) != 1 {
			http.Error(rw, "Inavlid URL", http.StatusBadRequest)
			return
		}

		if len(regexGroup[0]) != 2 {
			http.Error(rw, "Inavlid URL", http.StatusBadRequest)
			return
		}

		idString := regexGroup[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(rw, "Invalid ID", http.StatusBadRequest)
			return
		}

		p.updateProducts(id, rw, r)
		return 
	}

	// catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	listOfProducts := data.GetProducts()
	// convert to json
	err := listOfProducts.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
} 

func (p* Products) addProduct(rw http.ResponseWriter, r *http.Request){
	p.l.Println("Handle POST Product")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)

	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	data.AddProduct(prod)
}

func (p* Products) updateProducts(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Product")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)

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
