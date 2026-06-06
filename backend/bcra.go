package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const bcraBaseURL = "https://api.bcra.gob.ar"

type bcraClient struct {
	http *http.Client
}

func newBCRAClient() *bcraClient {
	return &bcraClient{
		http: &http.Client{Timeout: 8 * time.Second},
	}
}

// deudaResponse mirrors the Central de Deudores v1.0 response.
type deudaResponse struct {
	Status  int `json:"status"`
	Results struct {
		Identificacion int64  `json:"identificacion"`
		Denominacion   string `json:"denominacion"`
		Periodos       []struct {
			Periodo   string `json:"periodo"`
			Entidades []struct {
				Situacion       int   `json:"situacion"`
				Monto           int64 `json:"monto"`
				DiasAtrasoPago  int   `json:"diasAtrasoPago"`
			} `json:"entidades"`
		} `json:"periodos"`
	} `json:"results"`
}

// cotizacionResponse mirrors /estadisticascambiarias/v1.0/Cotizaciones.
type cotizacionResponse struct {
	Status  int `json:"status"`
	Results struct {
		Detalle []struct {
			CodigoMoneda    string  `json:"codigoMoneda"`
			TipoCotizacion  float64 `json:"tipoCotizacion"`
		} `json:"detalle"`
	} `json:"results"`
}

// normalizeCUIT strips hyphens to get 11 raw digits.
func normalizeCUIT(cuit string) string {
	return strings.ReplaceAll(cuit, "-", "")
}

func (c *bcraClient) getDeudas(cuit string) (*deudaResponse, bool, error) {
	digits := normalizeCUIT(cuit)
	url := fmt.Sprintf("%s/centraldedeudores/v1.0/Deudas/%s", bcraBaseURL, digits)
	resp, err := c.http.Get(url)
	if err != nil {
		// Network failure — fall back to mock so the rest of the flow works.
		log.Printf("│ [Central de Deudores] API inaccesible (%v) — usando datos simulados", err)
		return mockDeudas(digits), true, nil
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, false, nil // CUIT sin antecedentes
	}
	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("bcra deudas: status %d", resp.StatusCode)
	}
	var out deudaResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, false, err
	}
	return &out, false, nil
}

// mockDeudas genera un perfil crediticio determinístico basado en el CUIT.
// El último dígito define el escenario: 0-5 → normal, 6-7 → sin antecedentes, 8 → seguimiento especial, 9 → con problemas.
func mockDeudas(digits string) *deudaResponse {
	last, _ := strconv.Atoi(digits[len(digits)-1:])
	period := time.Now().AddDate(0, -1, 0).Format("200601")

	type entidad struct {
		sit   int
		monto int64
		atraso int
	}

	scenarios := map[int][]entidad{
		0: {{1, 48_000, 0}},
		1: {{1, 120_000, 0}, {1, 35_000, 0}},
		2: {{1, 0, 0}},
		3: {{1, 75_000, 0}},
		4: {{1, 210_000, 0}, {1, 90_000, 0}},
		5: {{1, 15_000, 0}},
		6: nil, // sin antecedentes
		7: nil,
		8: {{2, 180_000, 32}, {1, 60_000, 0}},
		9: {{3, 420_000, 75}, {2, 95_000, 18}},
	}

	entidades := scenarios[last]
	if entidades == nil {
		return nil // sin antecedentes en el sistema
	}

	denominaciones := []string{
		"COMERCIOS REGIONALES SRL", "DISTRIBUIDORA NORTE SA", "ALMACENES UNIDOS SRL",
		"COMERCIO MINORISTA SRL", "GRUPO RETAIL AR SA",
	}
	denom := denominaciones[last%len(denominaciones)]

	type entRow struct {
		Situacion      int   `json:"situacion"`
		Monto          int64 `json:"monto"`
		DiasAtrasoPago int   `json:"diasAtrasoPago"`
	}
	rows := make([]entRow, len(entidades))
	for i, e := range entidades {
		rows[i] = entRow{e.sit, e.monto, e.atraso}
	}

	out := &deudaResponse{}
	out.Results.Denominacion = denom
	out.Results.Periodos = []struct {
		Periodo   string `json:"periodo"`
		Entidades []struct {
			Situacion      int   `json:"situacion"`
			Monto          int64 `json:"monto"`
			DiasAtrasoPago int   `json:"diasAtrasoPago"`
		} `json:"entidades"`
	}{
		{Periodo: period, Entidades: func() []struct {
			Situacion      int   `json:"situacion"`
			Monto          int64 `json:"monto"`
			DiasAtrasoPago int   `json:"diasAtrasoPago"`
		} {
			r := make([]struct {
				Situacion      int   `json:"situacion"`
				Monto          int64 `json:"monto"`
				DiasAtrasoPago int   `json:"diasAtrasoPago"`
			}, len(rows))
			for i, row := range rows {
				r[i].Situacion = row.Situacion
				r[i].Monto = row.Monto
				r[i].DiasAtrasoPago = row.DiasAtrasoPago
			}
			return r
		}()},
	}
	return out
}

func (c *bcraClient) getCotizacionUSD() (float64, error) {
	url := fmt.Sprintf("%s/estadisticascambiarias/v1.0/Cotizaciones", bcraBaseURL)
	resp, err := c.http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	var out cotizacionResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return 0, err
	}
	for _, d := range out.Results.Detalle {
		if d.CodigoMoneda == "USD" {
			return d.TipoCotizacion, nil
		}
	}
	return 0, fmt.Errorf("USD not found in cotizaciones")
}

// CreditLimitFromCheck maps BCRA situation to an initial credit limit in ARS.
func CreditLimitFromCheck(check BCRACheck) int64 {
	switch {
	case !check.DeudorDisponible || check.PeorSituacion <= 1:
		return 50_000
	case check.PeorSituacion == 2:
		return 35_000
	case check.PeorSituacion == 3:
		return 15_000
	default: // 4, 5, 6
		return 5_000
	}
}

var situacionDesc = map[int]string{
	1: "Normal",
	2: "Con seguimiento especial / Riesgo bajo",
	3: "Con problemas / Riesgo medio",
	4: "Con alto riesgo de insolvencia",
	5: "Irrecuperable",
	6: "Irrecuperable por disposición técnica",
}

// RunChecks executes all BCRA checks for a CUIT, logs a full report, and returns a BCRACheck.
func (c *bcraClient) RunChecks(cuit string) BCRACheck {
	result := BCRACheck{Timestamp: time.Now()}

	deudas, isMock, deudasErr := c.getDeudas(cuit)
	usd, usdErr := c.getCotizacionUSD()

	if deudasErr == nil {
		result.DeudorDisponible = true
		if deudas != nil && len(deudas.Results.Periodos) > 0 {
			latest := deudas.Results.Periodos[0]
			result.Denominacion = deudas.Results.Denominacion
			result.UltimoPeriodo = latest.Periodo
			for _, ent := range latest.Entidades {
				if ent.Situacion > result.PeorSituacion {
					result.PeorSituacion = ent.Situacion
				}
				result.DeudaTotal += ent.Monto
			}
		}
	}

	if usdErr == nil {
		result.TipoCambioUSD = usd
	}

	logBCRAReport(cuit, result, deudas, isMock, deudasErr, usdErr)
	return result
}

func logBCRAReport(cuit string, r BCRACheck, deudas *deudaResponse, isMock bool, deudasErr, usdErr error) {
	sep := "────────────────────────────────────────────────"
	log.Printf("\n┌ BCRA CHECK · CUIT %s · %s\n%s", cuit, r.Timestamp.Format("2006-01-02 15:04:05"), sep)

	// ── Central de Deudores ──
	if deudasErr != nil {
		log.Printf("│ [Central de Deudores] ERROR: %v", deudasErr)
	} else if deudas == nil || len(deudas.Results.Periodos) == 0 {
		source := "Central de Deudores"
		if isMock {
			source += " [MOCK]"
		}
		log.Printf("│ [%s] Sin antecedentes en el sistema financiero", source)
	} else {
		source := "Central de Deudores"
		if isMock {
			source += " [MOCK]"
		}
		desc := situacionDesc[r.PeorSituacion]
		if desc == "" {
			desc = "Desconocida"
		}
		log.Printf("│ [%s]", source)
		log.Printf("│   Denominación   : %s", r.Denominacion)
		log.Printf("│   Último período : %s", r.UltimoPeriodo)
		log.Printf("│   Peor situación : %d — %s", r.PeorSituacion, desc)
		log.Printf("│   Deuda total    : $%d ARS", r.DeudaTotal)
		log.Printf("│   Detalle por entidad (período %s):", r.UltimoPeriodo)
		for i, ent := range deudas.Results.Periodos[0].Entidades {
			entDesc := situacionDesc[ent.Situacion]
			if entDesc == "" {
				entDesc = "?"
			}
			log.Printf("│     [%d] sit=%d (%s)  monto=$%d ARS  atraso=%d días",
				i+1, ent.Situacion, entDesc, ent.Monto, ent.DiasAtrasoPago)
		}
	}

	// ── Cotizaciones ──
	if usdErr != nil {
		log.Printf("│ [Cotizaciones BCRA] ERROR: %v", usdErr)
	} else {
		log.Printf("│ [Cotizaciones BCRA] USD referencia: $%.2f ARS", r.TipoCambioUSD)
	}

	log.Printf("└%s", sep)
}
