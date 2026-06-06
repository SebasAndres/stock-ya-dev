package main

import (
	"encoding/json"
	"net/http"
)

type handlers struct {
	store *Store
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func decode(r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
}

// POST /api/businesses
func (h *handlers) createBusiness(w http.ResponseWriter, r *http.Request) {
	var req CreateBusinessReq
	if err := decode(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.CUIT == "" {
		writeError(w, http.StatusBadRequest, "name and cuit are required")
		return
	}
	biz := h.store.createBusiness(req)
	writeJSON(w, http.StatusCreated, biz)
}

// GET /api/businesses/{id}
func (h *handlers) getBusiness(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	biz, ok := h.store.getBusiness(id)
	if !ok {
		writeError(w, http.StatusNotFound, "business not found")
		return
	}
	writeJSON(w, http.StatusOK, biz)
}

// GET /api/businesses/{id}/dashboard
func (h *handlers) getDashboard(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	stats, ok := h.store.dashboardStats(id)
	if !ok {
		writeError(w, http.StatusNotFound, "business not found")
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

// GET /api/products?provider=id
func (h *handlers) listProducts(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("provider")
	products := h.store.listProductsByProvider(provider)
	writeJSON(w, http.StatusOK, products)
}

// GET /api/providers
func (h *handlers) listProviders(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.store.providers)
}

// GET /api/branches
func (h *handlers) listBranches(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.store.branches)
}

// POST /api/businesses/{id}/advances
func (h *handlers) createAdvance(w http.ResponseWriter, r *http.Request) {
	bizID := r.PathValue("id")
	if _, ok := h.store.getBusiness(bizID); !ok {
		writeError(w, http.StatusNotFound, "business not found")
		return
	}

	var req CreateAdvanceReq
	if err := decode(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	advance, err := h.store.createAdvance(bizID, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, advance)
}

// GET /api/businesses/{id}/advances
func (h *handlers) listAdvances(w http.ResponseWriter, r *http.Request) {
	bizID := r.PathValue("id")
	if _, ok := h.store.getBusiness(bizID); !ok {
		writeError(w, http.StatusNotFound, "business not found")
		return
	}
	advances := h.store.listAdvances(bizID)
	writeJSON(w, http.StatusOK, advances)
}

// POST /api/advances/{id}/pay
func (h *handlers) payAdvance(w http.ResponseWriter, r *http.Request) {
	advanceID := r.PathValue("id")
	advance, err := h.store.payAdvance(advanceID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, advance)
}

// POST /api/advances/{id}/delivery
func (h *handlers) requestDelivery(w http.ResponseWriter, r *http.Request) {
	advanceID := r.PathValue("id")

	var req DeliveryReq
	if err := decode(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.BranchID == "" {
		writeError(w, http.StatusBadRequest, "sucursalId is required")
		return
	}

	advance, err := h.store.requestDelivery(advanceID, req.BranchID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, advance)
}
