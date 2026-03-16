package main

import (
	"net/http"
)

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	var profit float64

	query := `
	SELECT
	  COALESCE((SELECT SUM(total_price) FROM sales), 0)
	- COALESCE((SELECT SUM(total_cost) FROM purchases), 0)
	- COALESCE((SELECT SUM(amount) FROM expenses), 0)
	`

	if err := db.QueryRow(query).Scan(&profit); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data := struct {
		Profit float64
	}{
		Profit: profit,
	}

	renderTemplate(w, "dashboard.html", data)
}
