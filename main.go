package main

import (
	"log"
	"net/http"
	"text/template"
)

func main() {
	if err := initDB(); err != nil {
		log.Fatal(err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/sales/add", addSaleHandler)
	http.HandleFunc("/purchase", addPurchaseHandler)
	http.HandleFunc("/", dashboardHandler)
	http.HandleFunc("/sales", listSalesHandler)
	http.HandleFunc("/stock", stockHandler)
	http.HandleFunc("/product/add", addProductHandler)
	http.HandleFunc("/debts", debtsHandler)
	http.HandleFunc("/debts/pay", payDebtHandler)

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	t, err := template.ParseFiles(
		"templates/base.html",
		"templates/"+name,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
