package handlers

import (
	"log"
	"net/http"

	"github.com/dev-muser/learning_go/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	listOfProducts := data.GetProducts()
	// converting to json
	err := listOfProducts.ToJSON(rw)
	// d, err := json.Marshal(listOfProducts)
	// reason to change to encoder is that i dont have to buffer anything
	//into memory. Dont have to allocate memory for that data.
	// encoder is also faster than using marshal
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}
