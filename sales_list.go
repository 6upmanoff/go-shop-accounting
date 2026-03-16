package main

import (
	"net/http"
	"time"
)

type SaleRow struct {
	ID         int
	SaleDate   time.Time
	Product    string
	QtyKg      float64
	PricePerKg float64
	TotalPrice float64
	Payment    string
}

func listSalesHandler(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	if from == "" || to == "" {
		today := time.Now().Format("2006-01-02")
		from = today
		to = today
	}

	rows, err := db.Query(`
		SELECT
			s.id,
			s.sale_date,
			p.name,
			s.quantity_kg,
			s.price_per_kg,
			s.total_price,
			s.payment_type
		FROM sales s
		JOIN products p ON p.id = s.product_id
		WHERE s.sale_date BETWEEN $1 AND $2
		ORDER BY s.sale_date DESC, s.id DESC
	`, from, to)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var sales []SaleRow
	for rows.Next() {
		var s SaleRow
		if err := rows.Scan(
			&s.ID,
			&s.SaleDate,
			&s.Product,
			&s.QtyKg,
			&s.PricePerKg,
			&s.TotalPrice,
			&s.Payment,
		); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		sales = append(sales, s)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var totalAll, totalCash, totalDebt float64
	err = db.QueryRow(`
		SELECT
			COALESCE(SUM(total_price), 0),
			COALESCE(SUM(CASE WHEN payment_type='cash' THEN total_price ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN payment_type='debt' THEN total_price ELSE 0 END), 0)
		FROM sales
		WHERE sale_date BETWEEN $1 AND $2
	`, from, to).Scan(&totalAll, &totalCash, &totalDebt)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data := struct {
		From      string
		To        string
		TotalAll  float64
		TotalCash float64
		TotalDebt float64
		Sales     []SaleRow
	}{
		From:      from,
		To:        to,
		TotalAll:  totalAll,
		TotalCash: totalCash,
		TotalDebt: totalDebt,
		Sales:     sales,
	}

	renderTemplate(w, "sales.html", data)
}
