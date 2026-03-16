package main

import (
	"net/http"
)

func addPurchaseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		productID := r.FormValue("product_id")
		price := r.FormValue("price")
		qty := r.FormValue("qty")

		_, err := db.Exec(`
			INSERT INTO purchases (product_id, price_per_kg, quantity_kg, purchase_date)
			VALUES ($1, $2, $3, CURRENT_DATE)
		`, productID, price, qty)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	products, err := getProducts()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data := struct {
		Products []Product
	}{
		Products: products,
	}

	renderTemplate(w, "purchase.html", data)
}
