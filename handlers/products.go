package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dev-muser/learning_go/data"
	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Product")
	listOfProducts := data.GetProducts()
	err := listOfProducts.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")

	prod := r.Context().Value(KeyProduct{}).(data.Product) //safe to cast it
	data.AddProduct(&prod)
}

func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}
	p.l.Println("Handle PUT Products", id)
	prod := r.Context().Value(KeyProduct{}).(data.Product) //safe to cast it

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

// the way to use a context is to define a key (type - preffered or string)
type KeyProduct struct{}

// to avoid duplicate code, like prod.FromJSON, we can put it on a function or use a Middleware
// USE function in Gorilla allows you apply Middleware.
// MiddleWare is a http handler and can chain multiplers handlers together (like Cors - validate requests)

//MiddlewareProductValidation validates the product in the request and calls next if ok
func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		//validate/sanitize the product
		err = prod.Validate()
		if err != nil {
			p.l.Println("[ERROR] validating product", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating prodcut: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		// add the product to the context to use it later on.
		// Request has Context.
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		// call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, req)
	})
}
