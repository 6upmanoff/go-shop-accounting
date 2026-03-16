package main

import (
	"net/http"
)

type DebtRow struct {
	SaleID      int
	ClientName  string
	ProductName string
	QtyKg       float64
	TotalPrice  float64
	PaidAmount  float64
	Remaining   float64
	SaleDate    string
}

func debtsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
	SELECT
	s.id,
	c.name,
	p.name,
	s.quantity_kg,
	s.total_price,
	COALESCE(SUM(dp.amount), 0) AS paid_amount,
	s.total_price - COALESCE(SUM(dp.amount), 0) AS remaining,
	s.sale_date
FROM sales s
JOIN products p ON p.id = s.product_id
JOIN clients c ON c.id = s.client_id
LEFT JOIN debt_payments dp ON dp.sale_id = s.id
WHERE s.payment_type = 'debt'
GROUP BY s.id, c.name, p.name, s.quantity_kg, s.total_price, s.sale_date
HAVING s.total_price - COALESCE(SUM(dp.amount), 0) > 0
ORDER BY s.sale_date DESC, s.id DESC
	`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var debts []DebtRow
	var totalRemaining float64

	for rows.Next() {
		var d DebtRow

		err := rows.Scan(
			&d.SaleID,
			&d.ClientName,
			&d.ProductName,
			&d.QtyKg,
			&d.TotalPrice,
			&d.PaidAmount,
			&d.Remaining,
			&d.SaleDate,
		)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		totalRemaining += d.Remaining
		debts = append(debts, d)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data := struct {
		Debts          []DebtRow
		TotalRemaining float64
	}{
		Debts:          debts,
		TotalRemaining: totalRemaining,
	}

	renderTemplate(w, "debts.html", data)
}
