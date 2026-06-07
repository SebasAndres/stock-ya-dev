package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

type Store struct {
	db *sql.DB
}

func newStore(db *sql.DB) *Store {
	s := &Store{db: db}
	s.seedCatalog()
	return s
}

func newID() string {
	return fmt.Sprintf("%d%04d", time.Now().UnixMilli(), rand.Intn(10000))
}

func (s *Store) seedCatalog() {
	s.seedProviders()
	s.seedProducts()
	s.seedBranches()
}

func (s *Store) seedProviders() {
	rows := []Provider{
		{ID: "todos", Name: "Todos", Type: "icon"},
		{ID: "coto", Name: "Coto", Type: "img", Src: "/assets/coto.png"},
		{ID: "argenchino", Name: "Supermercado Chino", Type: "img", Src: "/assets/argenchino.jpeg", Badge: "新"},
		{ID: "diarco", Name: "Diarco", Type: "text", Text: "DIARCO"},
		{ID: "makro", Name: "Makro", Type: "text", Text: "MAKRO"},
	}
	for _, p := range rows {
		s.db.Exec(`INSERT IGNORE INTO providers (id,name,type,src,text_label,badge) VALUES (?,?,?,?,?,?)`,
			p.ID, p.Name, p.Type, p.Src, p.Text, p.Badge)
	}
}

func (s *Store) seedProducts() {
	catalog := []Product{
		{ID: "cc500x6", Name: "Coca-Cola 500ml (x6)", Price: 4200, Emoji: "🧃", Image: "/assets/coca.webp", Providers: []string{"coto", "makro"}},
		{ID: "fid", Name: "Fideos Marolio 500g", Price: 1100, Emoji: "🍝", Image: "/assets/fideos.webp", Providers: []string{"argenchino", "diarco", "makro"}},
		{ID: "arr", Name: "Arroz Gallo Oro 1kg", Price: 1400, Emoji: "🍚", Image: "/assets/arroz.webp", Providers: []string{"argenchino", "diarco"}},
		{ID: "cas", Name: "Casancream x36", Price: 3600, Emoji: "🍪", Image: "/assets/casancream.webp", Providers: []string{"argenchino", "coto"}},
		{ID: "ace", Name: "Aceite Natura 1.5L", Price: 2200, Emoji: "🫒", Image: "/assets/aceite.webp", Providers: []string{"coto", "makro", "diarco"}},
		{ID: "yer", Name: "Yerba Playadito 500g", Price: 1800, Emoji: "🧉", Providers: []string{"argenchino", "diarco"}},
		{ID: "azu", Name: "Azúcar Ledesma 1kg", Price: 950, Emoji: "🍬", Providers: []string{"coto", "argenchino", "makro"}},
	}
	for _, p := range catalog {
		s.db.Exec(`INSERT IGNORE INTO products (id,name,price,emoji,image) VALUES (?,?,?,?,?)`,
			p.ID, p.Name, p.Price, p.Emoji, p.Image)
		for _, pv := range p.Providers {
			s.db.Exec(`INSERT IGNORE INTO product_providers (product_id,provider_id) VALUES (?,?)`, p.ID, pv)
		}
	}
}

func (s *Store) seedBranches() {
	rows := []Branch{
		{ID: "centro", Name: "Sucursal Centro", Address: "Av. Corrientes 1234, CABA", Emoji: "🏙️"},
		{ID: "palermo", Name: "Sucursal Palermo", Address: "Thames 1800, Palermo", Emoji: "🌳"},
		{ID: "belgrano", Name: "Sucursal Belgrano", Address: "Cabildo 2500, Belgrano", Emoji: "🏘️"},
		{ID: "flores", Name: "Sucursal Flores", Address: "Av. Rivadavia 6000, Flores", Emoji: "🌸"},
	}
	for _, b := range rows {
		s.db.Exec(`INSERT IGNORE INTO branches (id,name,address,emoji) VALUES (?,?,?,?)`,
			b.ID, b.Name, b.Address, b.Emoji)
	}
}

// ── Businesses ──────────────────────────────────────────────────────────────

func (s *Store) createBusiness(req CreateBusinessReq) (*Business, error) {
	b := &Business{
		ID: newID(), Name: req.Name,
		CUIT: normalizeCUIT(req.CUIT), Phone: req.Phone,
		Address: req.Address, Type: req.Type,
		CreditLimit: 50000, AssessmentStatus: AssessmentPending,
		CreatedAt: time.Now(),
	}
	_, err := s.db.Exec(
		`INSERT INTO businesses (id,name,cuit,phone,address,type,credit_limit,credit_used,assessment_status,created_at) VALUES (?,?,?,?,?,?,?,?,?,?)`,
		b.ID, b.Name, b.CUIT, b.Phone, b.Address, string(b.Type), b.CreditLimit, b.CreditUsed, string(b.AssessmentStatus), b.CreatedAt)
	return b, err
}

func (s *Store) getBusiness(id string) (*Business, error) {
	b := &Business{}
	err := s.db.QueryRow(
		`SELECT id,name,cuit,COALESCE(phone,''),COALESCE(address,''),type,credit_limit,credit_used,assessment_status,created_at FROM businesses WHERE id=?`, id).
		Scan(&b.ID, &b.Name, &b.CUIT, &b.Phone, &b.Address, &b.Type, &b.CreditLimit, &b.CreditUsed, &b.AssessmentStatus, &b.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return b, err
}

func (s *Store) getByCUIT(cuit string) (*Business, error) {
	b := &Business{}
	err := s.db.QueryRow(
		`SELECT id,name,cuit,COALESCE(phone,''),COALESCE(address,''),type,credit_limit,credit_used,assessment_status,created_at FROM businesses WHERE cuit=?`,
		normalizeCUIT(cuit)).
		Scan(&b.ID, &b.Name, &b.CUIT, &b.Phone, &b.Address, &b.Type, &b.CreditLimit, &b.CreditUsed, &b.AssessmentStatus, &b.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return b, err
}

func (s *Store) setCreditLimit(bizID string, limit int64) error {
	_, err := s.db.Exec(`UPDATE businesses SET credit_limit=? WHERE id=?`, limit, bizID)
	return err
}

func (s *Store) submitAssessment(bizID string, _ SubmitAssessmentReq) (*Business, error) {
	res, err := s.db.Exec(`UPDATE businesses SET assessment_status=? WHERE id=?`, string(AssessmentApproved), bizID)
	if err != nil {
		return nil, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return nil, fmt.Errorf("business not found")
	}
	return s.getBusiness(bizID)
}

// ── Catalog ─────────────────────────────────────────────────────────────────

func (s *Store) listProviders() ([]Provider, error) {
	rows, err := s.db.Query(`SELECT id,name,type,COALESCE(src,''),COALESCE(text_label,''),COALESCE(badge,'') FROM providers`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Provider
	for rows.Next() {
		var p Provider
		if err := rows.Scan(&p.ID, &p.Name, &p.Type, &p.Src, &p.Text, &p.Badge); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *Store) listBranches() ([]Branch, error) {
	rows, err := s.db.Query(`SELECT id,name,COALESCE(address,''),COALESCE(emoji,'') FROM branches`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Branch
	for rows.Next() {
		var b Branch
		if err := rows.Scan(&b.ID, &b.Name, &b.Address, &b.Emoji); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func (s *Store) listProductsByProvider(providerID string) ([]Product, error) {
	var rows *sql.Rows
	var err error
	if providerID == "" || providerID == "todos" {
		rows, err = s.db.Query(`SELECT id,name,price,COALESCE(emoji,''),COALESCE(image,'') FROM products`)
	} else {
		rows, err = s.db.Query(
			`SELECT DISTINCT p.id,p.name,p.price,COALESCE(p.emoji,''),COALESCE(p.image,'') FROM products p JOIN product_providers pp ON p.id=pp.product_id WHERE pp.provider_id=?`,
			providerID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Emoji, &p.Image); err != nil {
			return nil, err
		}
		if p.Providers, err = s.loadProductProviders(p.ID); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func (s *Store) loadProductProviders(productID string) ([]string, error) {
	rows, err := s.db.Query(`SELECT provider_id FROM product_providers WHERE product_id=?`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, rows.Err()
}

func (s *Store) productByID(id string) (*Product, error) {
	p := &Product{}
	err := s.db.QueryRow(`SELECT id,name,price,COALESCE(emoji,''),COALESCE(image,'') FROM products WHERE id=?`, id).
		Scan(&p.ID, &p.Name, &p.Price, &p.Emoji, &p.Image)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return p, err
}

func (s *Store) branchByID(id string) (*Branch, error) {
	b := &Branch{}
	err := s.db.QueryRow(`SELECT id,name,COALESCE(address,''),COALESCE(emoji,'') FROM branches WHERE id=?`, id).
		Scan(&b.ID, &b.Name, &b.Address, &b.Emoji)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return b, err
}

// ── Advances ─────────────────────────────────────────────────────────────────

func newAdvance(bizID string, items []AdvanceItem, base int64, logistics LogisticsType) *Advance {
	var storageCost int64
	if logistics == LogisticsWarehouse {
		storageCost = base * 2 / 100
	}
	now := time.Now()
	return &Advance{
		ID: newID(), BusinessID: bizID, Items: items,
		Total: base + storageCost, Logistics: logistics, StorageCost: storageCost,
		AdvanceDate: now.Format("02/01/2006"), DueDate: now.AddDate(0, 0, 30).Format("02/01/2006"),
		Status: StatusCurrent, CreatedAt: now,
	}
}

func (s *Store) createAdvance(bizID string, req CreateAdvanceReq) (*Advance, error) {
	items, base, err := s.buildItems(req.Cart)
	if err != nil {
		return nil, err
	}
	if base == 0 {
		return nil, fmt.Errorf("cart is empty")
	}
	a := newAdvance(bizID, items, base, req.Logistics)
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var creditLimit, creditUsed int64
	err = tx.QueryRow(`SELECT credit_limit, credit_used FROM businesses WHERE id=? FOR UPDATE`, bizID).
		Scan(&creditLimit, &creditUsed)
	if err != nil {
		return nil, err
	}
	if creditUsed+a.Total > creditLimit {
		return nil, fmt.Errorf("crédito insuficiente: disponible $%d, solicitado $%d", creditLimit-creditUsed, a.Total)
	}
	if err := s.insertAdvanceTx(tx, a); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(`UPDATE businesses SET credit_used=credit_used+? WHERE id=?`, a.Total, bizID); err != nil {
		return nil, err
	}
	return a, tx.Commit()
}

func (s *Store) insertAdvanceTx(tx *sql.Tx, a *Advance) error {
	_, err := tx.Exec(
		`INSERT INTO advances (id,business_id,total,logistics,storage_cost,advance_date,due_date,status,delivery_status,target_branch,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
		a.ID, a.BusinessID, a.Total, string(a.Logistics), a.StorageCost, a.AdvanceDate, a.DueDate, string(a.Status), string(a.DeliveryStatus), a.TargetBranch, a.CreatedAt)
	if err != nil {
		return err
	}
	for _, item := range a.Items {
		if _, err = tx.Exec(
			`INSERT INTO advance_items (advance_id,product_id,name,emoji,image,qty,price) VALUES (?,?,?,?,?,?,?)`,
			a.ID, item.ProductID, item.Name, item.Emoji, item.Image, item.Qty, item.Price); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) buildItems(cart map[string]int) ([]AdvanceItem, int64, error) {
	var items []AdvanceItem
	var base int64
	for productID, qty := range cart {
		if qty <= 0 {
			continue
		}
		p, err := s.productByID(productID)
		if err != nil {
			return nil, 0, err
		}
		if p == nil {
			return nil, 0, fmt.Errorf("product %s not found", productID)
		}
		items = append(items, AdvanceItem{ProductID: p.ID, Name: p.Name, Emoji: p.Emoji, Image: p.Image, Qty: qty, Price: p.Price})
		base += p.Price * int64(qty)
	}
	return items, base, nil
}

func (s *Store) loadAdvanceItems(advanceID string) ([]AdvanceItem, error) {
	rows, err := s.db.Query(
		`SELECT product_id,name,COALESCE(emoji,''),COALESCE(image,''),qty,price FROM advance_items WHERE advance_id=?`, advanceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AdvanceItem
	for rows.Next() {
		var item AdvanceItem
		if err := rows.Scan(&item.ProductID, &item.Name, &item.Emoji, &item.Image, &item.Qty, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) loadAdvance(id string) (*Advance, error) {
	a := &Advance{}
	err := s.db.QueryRow(
		`SELECT id,business_id,total,logistics,storage_cost,advance_date,due_date,status,delivery_status,target_branch,created_at FROM advances WHERE id=?`, id).
		Scan(&a.ID, &a.BusinessID, &a.Total, &a.Logistics, &a.StorageCost, &a.AdvanceDate, &a.DueDate, &a.Status, &a.DeliveryStatus, &a.TargetBranch, &a.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	a.Items, err = s.loadAdvanceItems(a.ID)
	return a, err
}

func (s *Store) listAdvances(bizID string) ([]*Advance, error) {
	rows, err := s.db.Query(
		`SELECT id,business_id,total,logistics,storage_cost,advance_date,due_date,status,delivery_status,target_branch,created_at FROM advances WHERE business_id=? ORDER BY created_at DESC`, bizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var advances []*Advance
	for rows.Next() {
		a := &Advance{}
		err = rows.Scan(&a.ID, &a.BusinessID, &a.Total, &a.Logistics, &a.StorageCost, &a.AdvanceDate, &a.DueDate, &a.Status, &a.DeliveryStatus, &a.TargetBranch, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		if a.Items, err = s.loadAdvanceItems(a.ID); err != nil {
			return nil, err
		}
		advances = append(advances, a)
	}
	return advances, rows.Err()
}

func (s *Store) checkDeliveryEligible(advanceID string) error {
	var logistics LogisticsType
	var deliveryStatus DeliveryStatus
	err := s.db.QueryRow(`SELECT logistics,delivery_status FROM advances WHERE id=?`, advanceID).
		Scan(&logistics, &deliveryStatus)
	if err == sql.ErrNoRows {
		return fmt.Errorf("advance not found")
	}
	if err != nil {
		return err
	}
	if logistics != LogisticsWarehouse {
		return fmt.Errorf("advance is not stored in warehouse")
	}
	if deliveryStatus == DeliveryInTransit {
		return fmt.Errorf("delivery already in transit")
	}
	return nil
}

func (s *Store) requestDelivery(advanceID, branchID string) (*Advance, error) {
	branch, err := s.branchByID(branchID)
	if err != nil {
		return nil, err
	}
	if branch == nil {
		return nil, fmt.Errorf("branch not found")
	}
	if err := s.checkDeliveryEligible(advanceID); err != nil {
		return nil, err
	}
	if _, err := s.db.Exec(`UPDATE advances SET delivery_status=?,target_branch=? WHERE id=?`,
		string(DeliveryInTransit), branch.Name, advanceID); err != nil {
		return nil, err
	}
	return s.loadAdvance(advanceID)
}

func (s *Store) payAdvance(advanceID string) (*Advance, error) {
	a, err := s.loadAdvance(advanceID)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, fmt.Errorf("advance not found")
	}
	if a.Status == StatusPaid {
		return nil, fmt.Errorf("advance already paid")
	}
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`UPDATE advances SET status=? WHERE id=?`, string(StatusPaid), advanceID); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(`UPDATE businesses SET credit_used=GREATEST(0,credit_used-?) WHERE id=?`, a.Total, a.BusinessID); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	a.Status = StatusPaid
	return a, nil
}

func (s *Store) dashboardStats(bizID string) (*DashboardStats, error) {
	biz, err := s.getBusiness(bizID)
	if err != nil {
		return nil, err
	}
	if biz == nil {
		return nil, nil
	}
	var totalAdvances, totalPaid, activeProducts int
	s.db.QueryRow(`SELECT COUNT(*), COALESCE(SUM(status='paid'),0) FROM advances WHERE business_id=?`, bizID).
		Scan(&totalAdvances, &totalPaid)
	s.db.QueryRow(
		`SELECT COUNT(*) FROM advance_items ai JOIN advances a ON ai.advance_id=a.id WHERE a.business_id=? AND a.status!='paid'`, bizID).
		Scan(&activeProducts)
	return &DashboardStats{Business: biz, ActiveProducts: activeProducts, TotalAdvances: totalAdvances, TotalPaid: totalPaid}, nil
}
