# Feature: ML Scoring & Forecasting

**Status: specified, implementation deferred.** The three endpoints exist and return stub responses. This spec defines the intended behavior for when they are implemented.

**APIs:**
- `POST /api/ml/credit-score`
- `POST /api/ml/demand-forecast`
- `POST /api/ml/overdue-risk`

---

## Credit Score

**Purpose:** Recommend a credit limit for a business based on behavioral signals.

### Input signals (priority order)

| Signal | Source | Notes |
|---|---|---|
| Repayment punctuality | StockYa advance history | Days-to-pay distribution |
| Purchase frequency | StockYa advance history | Advances/month |
| Average ticket | StockYa advance history | Mean advance total |
| CUIT category | AFIP monotributo | Inferred annual revenue ceiling |
| QR payment history | Mercado Pago (future) | Estimated real revenue |

### Output

```json
{
  "score": 72,
  "riskLevel": "low",
  "creditLimitRecommendation": 150000,
  "message": "Based on 4 on-time repayments and consistent purchase frequency."
}
```

### Cold-start rules (no history)

| Available signal | Max recommended limit (ARS) |
|---|---|
| CUIT only | 10,000 |
| + AFIP category | 10% of declared monthly revenue |
| + 1 month of StockYa behavior | Model-driven |

---

## Demand Forecast

**Purpose:** Tell a store owner how much of each product to order in the next 30 days.

### Input

```json
{
  "businessId": "abc123",
  "productIds": ["cc225", "arr", "yer"]
}
```

### Output

```json
{
  "forecasts": [
    { "productId": "cc225", "predictedQty": 48, "confidence": 0.81 },
    { "productId": "arr",   "predictedQty": 20, "confidence": 0.74 }
  ],
  "windowDays": 30,
  "message": "Based on last 90 days of purchase history."
}
```

### Algorithm notes
- Use individual business history where available (n ≥ 10 advances).
- Fall back to network-level median for the same `BusinessType` in the same region when history is thin.
- `confidence` is lower when fewer historical data points are available.

---

## Overdue Risk

**Purpose:** Flag advances likely to miss their due date so operators can act proactively.

### Input

```json
{ "advanceId": "xyz789" }
```

### Output

```json
{
  "defaultProbability": 0.34,
  "daysAtRisk": 8,
  "riskLevel": "medium",
  "message": "This business had a 12-day late payment in its last advance."
}
```

### Risk levels

| defaultProbability | riskLevel |
|---|---|
| < 0.20 | low |
| 0.20 – 0.50 | medium |
| > 0.50 | high |

---

## Implementation notes (for when we build this)

- All three endpoints should be pure HTTP handlers — no side effects, no state mutations.
- Input validation: return 400 if `businessId`/`advanceId` is missing or unknown.
- Stub logic (no history): return `riskLevel: "pending"` and explain in `message`.
- The Go file `backend/ml.go` is intentionally isolated so ML logic can be swapped independently.
