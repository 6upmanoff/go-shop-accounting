package main

import (
	"net/http"
)

func addProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		category := r.FormValue("category")

		_, err := db.Exec(`
			INSERT INTO products (name, category)
			VALUES ($1, $2)
		`, name, category)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	renderTemplate(w, "product_add.html", nil)
}
