package main

import (
	"log"
	"net/http"
	"os"
	"testModule/handlers"
)

func main() {
	// -------------- First edition of the code before refactoring:
	// http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		// log.Println("Initial Commit")
		// d, err := io.ReadAll(r.Body)

		// if err != nil {
		// 	http.Error(rw, "Error happened", http.StatusBadRequest)
		// 	return
		// }
		// fmt.Fprintf(rw, "Input %s", d)
	// })

	// -------------- Refactored version using hello handler
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	helloHandler := handlers.NewHello(l)

	goodbyeHandler := handlers.NewGoodBye(l)

	serveMux := http.NewServeMux()
	serveMux.Handle("/", helloHandler)
	serveMux.Handle("/goodbye" , goodbyeHandler)

	http.ListenAndServe(":9090", serveMux)
}
