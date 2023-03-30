package main

// Build a REST API using only the standard library and no 3rd party routing libraries/dependencies.

// The REST API will regulate access to the Product resource. We'll create a Product Handler to handle the routing for the product REST API endpoints.

// We won't use and populate a DB to persist the products for now. We'll use a slice to store products in memory.

// when you allocate an instance of product handler to a variable, that variable can change the original values of the fields in the producthandler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// Model representing the Product resource
type Product struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Slice that holds Products in memory a la slice of Products
type Products []Product

// Implements the handler interface and handles requests to the products API endpoint and all the routing for products.
// struct that has a slice, products of type Products which holds Product structs
type productHandler struct {
	// A lock that allows one to lock access to the productHandler's critical section i.e. product slice when a request is interacting with the handler
	// to prevent a race condition where one request modifies the product slice before another can read from it causing inconsistency. A scenario
	// which could occur when requests access the products slice concurrently or in parallel as each request to the http server spins up a new goroutine.
	sync.Mutex
	products Products
}

// ServeHTTP is defined on a pointer to the productHandler and as such productHandler now implements the handler interface.
// Handles the request differently depending on the HTTP Request method.
func (ph *productHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		ph.post(w, r)
	case "GET":
		ph.get(w, r)
	case "PUT":
		ph.put(w, r)
	case "DELETE":
		ph.delete(w, r)
	default:
		// if the request method does not match any of the above, respond with an error.
		returnErrorResponse(w, http.StatusMethodNotAllowed, "Invalid HTTP method")
	}
}

func main() {
	port := ":8081"

	// allocate an instance of a pointer to productHandler to pHanlder
	pHandler := &productHandler{
		// Product slice literal
		products: Products{
			Product{"food", 10.00},
			Product{"car", 250.00},
			Product{"gadgets", 50.00},
		},
	}

	// registers a variable, pHandler whose type, productHandler implements the handler interface as the handler for the /products route
	http.Handle("/products", pHandler)
	http.Handle("/products/", pHandler)

	// registered an inlined anonymous function that as the handler for the root path "/"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hey!")
	})

	log.Fatal(http.ListenAndServe(port, nil))
}

// define http methods on a pointer to the productHandler and as such an instance of a pointer to the productHandler or a pointer receiver can call post
func (ph *productHandler) post(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Post!")
}

// method defined on productHandler that handles GET requests for the related Url pattern
// handles GET on /products for all products and /products/ for a specific product
func (ph *productHandler) get(w http.ResponseWriter, r *http.Request) {
	// concurrency handling
	// Unlock access to the Product slice when get is done using the mutex
	defer ph.Unlock()
	// Lock access to the Product slice such that only this GET request can interact with the Product slice at this time.
	ph.Lock()

	// return the products list
	returnJSONResponse(w, http.StatusOK, ph.products)
}

func (ph *productHandler) put(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Put!")
}

func (ph *productHandler) delete(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Delete!")
}

// The type of the data argument can be any type and so to represent that, we use the interface{}. In kotlin this would be <Any>
func returnJSONResponse(w http.ResponseWriter, code int, data interface{}) {
	response, err := json.Marshal(data)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// response is of type application/json
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func returnErrorResponse(w http.ResponseWriter, code int, msg string) {
	returnJSONResponse(w, code, map[string]string{"error": msg})
}
