package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Store struct {
	mu          sync.RWMutex
	businesses  map[string]*Business
	advances    map[string]*Advance
	bizAdvances map[string][]string // businessID → []advanceID

	products  []Product
	providers []Provider
	branches  []Branch
}

func newStore() *Store {
	s := &Store{
		businesses:  make(map[string]*Business),
		advances:    make(map[string]*Advance),
		bizAdvances: make(map[string][]string),
	}
	s.seedCatalog()
	return s
}

func (s *Store) seedCatalog() {
	s.providers = []Provider{
		{ID: "todos", Name: "Todos", Type: "icon"},
		{ID: "coto", Name: "Coto", Type: "img", Src: "/assets/coto.png"},
		{ID: "argenchino", Name: "Supermercado Chino", Type: "img", Src: "/assets/argenchino.jpeg", Badge: "新"},
		{ID: "diarco", Name: "Diarco", Type: "text", Text: "DIARCO"},
		{ID: "makro", Name: "Makro", Type: "text", Text: "MAKRO"},
	}
	s.products = []Product{
		{ID: "cc225", Name: "Coca-Cola 2.25L", Price: 2800, Emoji: "🥤", Image: "/assets/coca.webp", Providers: []string{"coto", "diarco"}},
		{ID: "cc500x6", Name: "Coca-Cola 500ml (x6)", Price: 4200, Emoji: "🧃", Image: "/assets/coca.webp", Providers: []string{"coto", "makro"}},
		{ID: "fid", Name: "Fideos Marolio 500g", Price: 1100, Emoji: "🍝", Image: "/assets/fideos.webp", Providers: []string{"argenchino", "diarco", "makro"}},
		{ID: "arr", Name: "Arroz Gallo Oro 1kg", Price: 1400, Emoji: "🍚", Image: "/assets/arroz.webp", Providers: []string{"argenchino", "diarco"}},
		{ID: "cas", Name: "Casancream x36", Price: 3600, Emoji: "🍪", Image: "/assets/casancream.webp", Providers: []string{"argenchino", "coto"}},
		{ID: "ace", Name: "Aceite Natura 1.5L", Price: 2200, Emoji: "🫒", Image: "/assets/aceite.webp", Providers: []string{"coto", "makro", "diarco"}},
		{ID: "yer", Name: "Yerba Playadito 500g", Price: 1800, Emoji: "🧉", Providers: []string{"argenchino", "diarco"}},
		{ID: "azu", Name: "Azúcar Ledesma 1kg", Price: 950, Emoji: "🍬", Providers: []string{"coto", "argenchino", "makro"}},
	}
	s.branches = []Branch{
		{ID: "centro", Name: "Sucursal Centro", Address: "Av. Corrientes 1234, CABA", Emoji: "🏙️"},
		{ID: "palermo", Name: "Sucursal Palermo", Address: "Thames 1800, Palermo", Emoji: "🌳"},
		{ID: "belgrano", Name: "Sucursal Belgrano", Address: "Cabildo 2500, Belgrano", Emoji: "🏘️"},
		{ID: "flores", Name: "Sucursal Flores", Address: "Av. Rivadavia 6000, Flores", Emoji: "🌸"},
	}
}

func (s *Store) newID() string {
	return fmt.Sprintf("%d%04d", time.Now().UnixMilli(), rand.Intn(10000))
}

func (s *Store) createBusiness(req CreateBusinessReq) *Business {
	s.mu.Lock()
	b := &Business{
		ID:          s.newID(),
		Name:        req.Name,
		CUIT:        req.CUIT,
		Phone:       req.Phone,
		Address:     req.Address,
		Type:        req.Type,
		CreditLimit: 50000,
		CreditUsed:  0,
		CreatedAt:   time.Now(),
	}
	s.businesses[b.ID] = b
	s.mu.Unlock()
	return b
}

func (s *Store) getBusiness(id string) (*Business, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.businesses[id]
	return b, ok
}

func (s *Store) listProductsByProvider(providerID string) []Product {
	if providerID == "" || providerID == "todos" {
		return s.products
	}
	out := make([]Product, 0)
	for _, p := range s.products {
		for _, pv := range p.Providers {
			if pv == providerID {
				out = append(out, p)
				break
			}
		}
	}
	return out
}

func (s *Store) productByID(id string) (Product, bool) {
	for _, p := range s.products {
		if p.ID == id {
			return p, true
		}
	}
	return Product{}, false
}

func (s *Store) branchByID(id string) (Branch, bool) {
	for _, b := range s.branches {
		if b.ID == id {
			return b, true
		}
	}
	return Branch{}, false
}

func (s *Store) createAdvance(bizID string, req CreateAdvanceReq) (*Advance, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.businesses[bizID]; !ok {
		return nil, fmt.Errorf("business not found")
	}

	items, base, err := s.buildItems(req.Cart)
	if err != nil {
		return nil, err
	}
	if base == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	var storageCost int64
	if req.Logistics == LogisticsWarehouse {
		storageCost = base * 2 / 100
	}

	now := time.Now()
	a := &Advance{
		ID:          s.newID(),
		BusinessID:  bizID,
		Items:       items,
		Total:       base + storageCost,
		Logistics:   req.Logistics,
		StorageCost: storageCost,
		AdvanceDate: now.Format("02/01/2006"),
		DueDate:     now.AddDate(0, 0, 30).Format("02/01/2006"),
		Status:      StatusCurrent,
		CreatedAt:   now,
	}

	s.advances[a.ID] = a
	s.bizAdvances[bizID] = append(s.bizAdvances[bizID], a.ID)
	s.businesses[bizID].CreditUsed += a.Total
	return a, nil
}

func (s *Store) buildItems(cart map[string]int) ([]AdvanceItem, int64, error) {
	var items []AdvanceItem
	var base int64
	for productID, qty := range cart {
		if qty <= 0 {
			continue
		}
		p, ok := s.productByID(productID)
		if !ok {
			return nil, 0, fmt.Errorf("product %s not found", productID)
		}
		items = append(items, AdvanceItem{
			ProductID: p.ID,
			Name:      p.Name,
			Emoji:     p.Emoji,
			Image:     p.Image,
			Qty:       qty,
			Price:     p.Price,
		})
		base += p.Price * int64(qty)
	}
	return items, base, nil
}

func (s *Store) listAdvances(bizID string) []*Advance {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ids := s.bizAdvances[bizID]
	out := make([]*Advance, 0, len(ids))
	for _, id := range ids {
		if a, ok := s.advances[id]; ok {
			out = append(out, a)
		}
	}
	return out
}

func (s *Store) requestDelivery(advanceID, branchID string) (*Advance, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	a, ok := s.advances[advanceID]
	if !ok {
		return nil, fmt.Errorf("advance not found")
	}
	if a.Logistics != LogisticsWarehouse {
		return nil, fmt.Errorf("advance is not stored in warehouse")
	}
	if a.DeliveryStatus == DeliveryInTransit {
		return nil, fmt.Errorf("delivery already in transit")
	}

	branch, ok := s.branchByID(branchID)
	if !ok {
		return nil, fmt.Errorf("branch not found")
	}

	a.DeliveryStatus = DeliveryInTransit
	a.TargetBranch = branch.Name
	return a, nil
}

func (s *Store) dashboardStats(bizID string) (*DashboardStats, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	biz, ok := s.businesses[bizID]
	if !ok {
		return nil, false
	}

	ids := s.bizAdvances[bizID]
	var activeProducts, paid int
	for _, id := range ids {
		a, exists := s.advances[id]
		if !exists {
			continue
		}
		if a.Status == StatusPaid {
			paid++
		} else {
			activeProducts += len(a.Items)
		}
	}

	return &DashboardStats{
		Business:       biz,
		ActiveProducts: activeProducts,
		TotalAdvances:  len(ids),
		TotalPaid:      paid,
	}, true
}

func (s *Store) payAdvance(advanceID string) (*Advance, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	a, ok := s.advances[advanceID]
	if !ok {
		return nil, fmt.Errorf("advance not found")
	}
	if a.Status == StatusPaid {
		return nil, fmt.Errorf("advance already paid")
	}

	a.Status = StatusPaid
	biz := s.businesses[a.BusinessID]
	biz.CreditUsed -= a.Total
	if biz.CreditUsed < 0 {
		biz.CreditUsed = 0
	}
	return a, nil
}
