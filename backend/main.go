package main

import (
	"log"
	"net/http"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	store := newStore(openDB())
	h := &handlers{store: store, bcra: newBCRAClient()}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/auth/login", h.login)
	mux.HandleFunc("POST /api/businesses", h.createBusiness)
	mux.HandleFunc("GET /api/businesses/{id}", h.getBusiness)
	mux.HandleFunc("GET /api/businesses/{id}/dashboard", h.getDashboard)
	mux.HandleFunc("GET /api/businesses/{id}/advances", h.listAdvances)
	mux.HandleFunc("POST /api/businesses/{id}/advances", h.createAdvance)
	mux.HandleFunc("POST /api/businesses/{id}/assessment", h.submitAssessment)

	mux.HandleFunc("GET /api/products", h.listProducts)
	mux.HandleFunc("GET /api/providers", h.listProviders)
	mux.HandleFunc("GET /api/branches", h.listBranches)

	mux.HandleFunc("POST /api/advances/{id}/pay", h.payAdvance)
	mux.HandleFunc("POST /api/advances/{id}/delivery", h.requestDelivery)

	mux.HandleFunc("POST /api/ml/credit-score", h.creditScore)
	mux.HandleFunc("POST /api/ml/demand-forecast", h.demandForecast)
	mux.HandleFunc("POST /api/ml/overdue-risk", h.overdueRisk)

	mux.Handle("/", http.FileServer(http.Dir("../frontend")))

	log.Println("StockYa backend running on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(mux)))
}
