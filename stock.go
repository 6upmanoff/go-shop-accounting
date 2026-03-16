package main

import (
	"net/http"
)

type StockRow struct {
	Name      string
	Purchased float64
	Sold      float64
	Stock     float64
}

func stockHandler(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query(`
	SELECT
    p.name,
    COALESCE(pu.total_purchased, 0) AS purchased,
    COALESCE(s.total_sold, 0) AS sold,
    COALESCE(pu.total_purchased, 0) - COALESCE(s.total_sold, 0) AS stock
FROM products p

LEFT JOIN (
    SELECT product_id, SUM(quantity_kg) AS total_purchased
    FROM purchases
    GROUP BY product_id
) pu ON pu.product_id = p.id

LEFT JOIN (
    SELECT product_id, SUM(quantity_kg) AS total_sold
    FROM sales
    GROUP BY product_id
) s ON s.product_id = p.id

ORDER BY p.name
	`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var items []StockRow

	for rows.Next() {

		var s StockRow

		err := rows.Scan(
			&s.Name,
			&s.Purchased,
			&s.Sold,
			&s.Stock,
		)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		items = append(items, s)
	}

	data := struct {
		Items []StockRow
	}{
		Items: items,
	}

	renderTemplate(w, "stock.html", data)
}
