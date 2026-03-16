package main

import (
	"net/http"
)

type Product struct {
	ID   int
	Name string
}

type Client struct {
	ID   int
	Name string
}

func getClients() ([]Client, error) {
	rows, err := db.Query(`SELECT id, name FROM clients ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []Client
	for rows.Next() {
		var c Client
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}
	return clients, rows.Err()
}

func getProducts() ([]Product, error) {
	rows, err := db.Query(`SELECT id, name FROM products ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, rows.Err()
}

func addSaleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		productID := r.FormValue("product_id")
		price := r.FormValue("price")
		qty := r.FormValue("qty")
		payment := r.FormValue("payment")
		clientID := r.FormValue("client_id")

		var clientValue interface{}

		if payment == "debt" {
			if clientID == "" {
				http.Error(w, "Для продажи в долг нужно выбрать клиента", http.StatusBadRequest)
				return
			}
			clientValue = clientID
		} else {
			clientValue = nil
		}

		_, err := db.Exec(`
			INSERT INTO sales (product_id, price_per_kg, quantity_kg, payment_type, sale_date, client_id)
			VALUES ($1, $2, $3, $4, CURRENT_DATE, $5)
		`, productID, price, qty, payment, clientValue)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		http.Redirect(w, r, "/sales", http.StatusSeeOther)
		return
	}

	products, err := getProducts()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	clients, err := getClients()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data := struct {
		Products []Product
		Clients  []Client
	}{
		Products: products,
		Clients:  clients,
	}

	renderTemplate(w, "sales_add.html", data)
}
