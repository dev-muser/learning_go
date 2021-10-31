package handlers

import (
	"net/http"
	"strconv"

	"github.com/dev-muser/learning_go/data"
	"github.com/gorilla/mux"
)

func (p *Products) DeleteProducts(rw http.ResponseWriter, r *http.Request) {
	// this will always convert because of the router
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	p.l.Println("Handle DELETE Product", id)

	err := data.DeleteProduct(id)

	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}