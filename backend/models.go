package main

import "time"

type BusinessType string

const (
	Almacen    BusinessType = "almacen"
	Kiosco     BusinessType = "kiosco"
	Verduleria BusinessType = "verduleria"
	Other      BusinessType = "otro"
)

type Business struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	CUIT        string       `json:"cuit"`
	Phone       string       `json:"phone"`
	Address     string       `json:"address"`
	Type        BusinessType `json:"type"`
	CreditLimit int64        `json:"creditLimit"`
	CreditUsed  int64        `json:"creditUsed"`
	CreatedAt   time.Time    `json:"createdAt"`
}

type Provider struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Src   string `json:"src,omitempty"`
	Text  string `json:"text,omitempty"`
	Badge string `json:"badge,omitempty"`
}

type Product struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Price     int64    `json:"price"`
	Emoji     string   `json:"emoji"`
	Image     string   `json:"image,omitempty"`
	Providers []string `json:"providers"`
}

type LogisticsType string

const (
	LogisticsDelivery  LogisticsType = "envio"
	LogisticsWarehouse LogisticsType = "deposito"
)

type AdvanceStatus string

const (
	StatusCurrent AdvanceStatus = "ok"
	StatusWarning AdvanceStatus = "warn"
	StatusOverdue AdvanceStatus = "danger"
	StatusPaid    AdvanceStatus = "paid"
)

type DeliveryStatus string

const (
	DeliveryNone      DeliveryStatus = ""
	DeliveryInTransit DeliveryStatus = "en_transito"
)

type AdvanceItem struct {
	ProductID string `json:"id"`
	Name      string `json:"name"`
	Emoji     string `json:"emoji"`
	Image     string `json:"image,omitempty"`
	Qty       int    `json:"qty"`
	Price     int64  `json:"price"`
}

type Advance struct {
	ID             string         `json:"id"`
	BusinessID     string         `json:"businessId"`
	Items          []AdvanceItem  `json:"items"`
	Total          int64          `json:"total"`
	Logistics      LogisticsType  `json:"logistics"`
	StorageCost    int64          `json:"depositoCost"`
	AdvanceDate    string         `json:"adelantoDate"`
	DueDate        string         `json:"dueDate"`
	Status         AdvanceStatus  `json:"status"`
	DeliveryStatus DeliveryStatus `json:"deliveryStatus,omitempty"`
	TargetBranch   string         `json:"targetSucursal,omitempty"`
	CreatedAt      time.Time      `json:"createdAt"`
}

type Branch struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Emoji   string `json:"emoji"`
}

type DashboardStats struct {
	Business       *Business `json:"business"`
	ActiveProducts int       `json:"activeProducts"`
	TotalAdvances  int       `json:"totalAdelantos"`
	TotalPaid      int       `json:"totalPagados"`
}

// Request bodies

type CreateBusinessReq struct {
	Name    string       `json:"name"`
	CUIT    string       `json:"cuit"`
	Phone   string       `json:"phone"`
	Address string       `json:"address"`
	Type    BusinessType `json:"type"`
}

type CreateAdvanceReq struct {
	Cart      map[string]int `json:"cart"`
	Logistics LogisticsType  `json:"logistics"`
}

type DeliveryReq struct {
	BranchID string `json:"sucursalId"`
}
