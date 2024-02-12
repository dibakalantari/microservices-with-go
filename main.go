package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("Initial Commit")
		d, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(rw, "Error happened", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(rw, "Input %s", d)
	})

	http.ListenAndServe(":9090", nil)
}
