package handlers

import (
	"net/http"

	"github.com/dev-muser/learning_go/data"
)

// swagger:route DELETE /products/{id} products deleteProduct
// responses:
// 		201: noContent
// Delete deletes a product from datastore
func (p *Products) Delete(rw http.ResponseWriter, r *http.Request) {
	id := getProductID(r)
	p.l.Println("Handle DELETE Product", id)

	err := data.DeleteProduct(id)

	if err == data.ErrProductNotFound {
		p.l.Println("[ERROR] deleting record id does not exist")

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	if err != nil {
		p.l.Println("[ERROR] deleting record", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
