package data

import "testing"

func TestChecksValidation(t *testing.T) {
	product := &Product{
		Name:  "Vodka",
		Price: 1,
		SKU:   "abc-asd-xxx",
	}

	err := product.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
