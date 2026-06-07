package main

import (
	"encoding/json"
	"net/http"
)

type handlers struct {
	store *Store
	bcra  *bcraClient
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
	biz, err := h.store.createBusiness(req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create business")
		return
	}
	check := h.bcra.RunChecks(req.CUIT)
	h.store.setCreditLimit(biz.ID, CreditLimitFromCheck(check))
	biz.CreditLimit = CreditLimitFromCheck(check)
	writeJSON(w, http.StatusCreated, biz)
}

// GET /api/businesses/{id}
func (h *handlers) getBusiness(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	biz, err := h.store.getBusiness(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if biz == nil {
		writeError(w, http.StatusNotFound, "business not found")
		return
	}
	writeJSON(w, http.StatusOK, biz)
}

// GET /api/businesses/{id}/dashboard
func (h *handlers) getDashboard(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	stats, err := h.store.dashboardStats(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if stats == nil {
		writeError(w, http.StatusNotFound, "business not found")
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

// GET /api/products?provider=id
func (h *handlers) listProducts(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("provider")
	products, err := h.store.listProductsByProvider(provider)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, products)
}

// GET /api/providers
func (h *handlers) listProviders(w http.ResponseWriter, r *http.Request) {
	providers, err := h.store.listProviders()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, providers)
}

// GET /api/branches
func (h *handlers) listBranches(w http.ResponseWriter, r *http.Request) {
	branches, err := h.store.listBranches()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, branches)
}

// POST /api/businesses/{id}/advances
func (h *handlers) createAdvance(w http.ResponseWriter, r *http.Request) {
	bizID := r.PathValue("id")
	biz, err := h.store.getBusiness(bizID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if biz == nil {
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
	biz, err := h.store.getBusiness(bizID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if biz == nil {
		writeError(w, http.StatusNotFound, "business not found")
		return
	}
	advances, err := h.store.listAdvances(bizID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
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

// POST /api/businesses/{id}/assessment
func (h *handlers) submitAssessment(w http.ResponseWriter, r *http.Request) {
	bizID := r.PathValue("id")
	var req SubmitAssessmentReq
	if err := decode(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	biz, err := h.store.submitAssessment(bizID, req)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, biz)
}

// POST /api/auth/login
func (h *handlers) login(w http.ResponseWriter, r *http.Request) {
	var req LoginReq
	if err := decode(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.CUIT == "" {
		writeError(w, http.StatusBadRequest, "cuit is required")
		return
	}
	biz, err := h.store.getByCUIT(req.CUIT)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if biz == nil {
		writeError(w, http.StatusNotFound, "negocio no encontrado — registrate primero")
		return
	}
	check := h.bcra.RunChecks(req.CUIT)
	writeJSON(w, http.StatusOK, LoginResponse{Business: biz, BCRACheck: check})
}

// POST /api/advances/{id}/delivery
func (h *handlers) requestDelivery(w http.ResponseWriter, r *http.Request) {
	advanceID := r.PathValue("id")
	advance, err := h.store.requestDelivery(advanceID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, advance)
}
