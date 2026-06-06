package main

import "net/http"

// POST /api/ml/credit-score
// Input: { "businessId": "..." }
// Returns a creditworthiness score and limit recommendation.
func (h *handlers) creditScore(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"score":                    0,
		"riskLevel":                "pending",
		"creditLimitRecommendation": 0,
		"message":                  "not implemented",
	})
}

// POST /api/ml/demand-forecast
// Input: { "businessId": "...", "productIds": ["cc225", ...] }
// Returns predicted reorder quantities per product for the next 30 days.
func (h *handlers) demandForecast(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"forecasts": []any{},
		"windowDays": 30,
		"message":   "not implemented",
	})
}

// POST /api/ml/overdue-risk
// Input: { "advanceId": "..." }
// Returns probability of default and days-at-risk estimate.
func (h *handlers) overdueRisk(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"defaultProbability": 0.0,
		"daysAtRisk":         0,
		"riskLevel":          "pending",
		"message":            "not implemented",
	})
}
