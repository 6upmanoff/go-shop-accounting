package main

import (
	"net/http"
	"strconv"
)

type DebtPaymentPage struct {
	SaleID      int
	ClientName  string
	ProductName string
	TotalPrice  float64
	PaidAmount  float64
	Remaining   float64
}

func payDebtHandler(w http.ResponseWriter, r *http.Request) {
	saleIDStr := r.URL.Query().Get("sale_id")
	if saleIDStr == "" {
		http.Error(w, "sale_id не указан", 400)
		return
	}

	saleID, err := strconv.Atoi(saleIDStr)
	if err != nil {
		http.Error(w, "некорректный sale_id", 400)
		return
	}

	if r.Method == http.MethodPost {
		amount := r.FormValue("amount")

		_, err := db.Exec(`
			INSERT INTO debt_payments (sale_id, amount, payment_date)
			VALUES ($1, $2, CURRENT_DATE)
		`, saleID, amount)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		http.Redirect(w, r, "/debts", http.StatusSeeOther)
		return
	}

	var debt DebtPaymentPage

	err = db.QueryRow(`
		SELECT
			c.name,
			p.name,
			s.total_price,
			COALESCE(SUM(dp.amount), 0) AS paid_amount,
			s.total_price - COALESCE(SUM(dp.amount), 0) AS remaining
		FROM sales s
		JOIN products p ON p.id = s.product_id
		JOIN clients c ON c.id = s.client_id
		LEFT JOIN debt_payments dp ON dp.sale_id = s.id
		WHERE s.id = $1
		GROUP BY c.name, p.name, s.total_price
	`, saleID).Scan(
		&debt.ClientName,
		&debt.ProductName,
		&debt.TotalPrice,
		&debt.PaidAmount,
		&debt.Remaining,
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	debt.SaleID = saleID

	renderTemplate(w, "debt_payment.html", debt)
}
