package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

func (c *bcraClient) getDeudas(cuit string) (*deudaResponse, error) {
	digits := normalizeCUIT(cuit)
	url := fmt.Sprintf("%s/centraldedeudores/v1.0/Deudas/%s", bcraBaseURL, digits)
	resp, err := c.http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // CUIT sin antecedentes en el sistema
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bcra deudas: status %d", resp.StatusCode)
	}
	var out deudaResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
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

	deudas, deudasErr := c.getDeudas(cuit)
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

	logBCRAReport(cuit, result, deudas, deudasErr, usdErr)
	return result
}

func logBCRAReport(cuit string, r BCRACheck, deudas *deudaResponse, deudasErr, usdErr error) {
	sep := "────────────────────────────────────────────────"
	log.Printf("\n┌ BCRA CHECK · CUIT %s · %s\n%s", cuit, r.Timestamp.Format("2006-01-02 15:04:05"), sep)

	// ── Central de Deudores ──
	if deudasErr != nil {
		log.Printf("│ [Central de Deudores] ERROR: %v", deudasErr)
	} else if deudas == nil || len(deudas.Results.Periodos) == 0 {
		log.Printf("│ [Central de Deudores] Sin antecedentes en el sistema financiero")
	} else {
		desc := situacionDesc[r.PeorSituacion]
		if desc == "" {
			desc = "Desconocida"
		}
		log.Printf("│ [Central de Deudores]")
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
