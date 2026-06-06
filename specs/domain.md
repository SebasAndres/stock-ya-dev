# Domain Model — StockYa

## Core Entities

### Business
A registered kiosco, almacén, verdulería, or similar store. Owns a credit line.

| Field | Type | Notes |
|---|---|---|
| id | string | Auto-generated |
| name | string | Required |
| cuit | string | Required; Argentine tax ID |
| phone | string | Optional |
| address | string | Optional |
| type | BusinessType | almacen \| kiosco \| verduleria \| otro |
| creditLimit | int64 | ARS cents; starts at 50,000 on registration |
| creditUsed | int64 | ARS cents; sum of active (unpaid) advance totals |
| createdAt | time | Auto-set |

### Advance
A stock draw against a business's credit line. Represents physical products to be sourced from wholesale and either delivered or held in a StockYa warehouse.

| Field | Type | Notes |
|---|---|---|
| id | string | Auto-generated |
| businessId | string | FK → Business |
| items | AdvanceItem[] | Products, quantities, unit prices at time of creation |
| total | int64 | ARS cents: Σ(item.price × item.qty) + storageCost |
| logistics | LogisticsType | envio (direct ship) or deposito (warehouse hold) |
| depositoCost | int64 | 2% of product total when logistics=deposito; 0 otherwise |
| adelantoDate | string | DD/MM/YYYY creation date |
| dueDate | string | DD/MM/YYYY; adelantoDate + 30 days |
| status | AdvanceStatus | ok → warn → danger → paid |
| deliveryStatus | DeliveryStatus | "" or en_transito |
| targetSucursal | string | Branch name set when delivery is requested |
| createdAt | time | Auto-set |

### Product
An SKU available in the wholesale catalog.

| Field | Type | Notes |
|---|---|---|
| id | string | Stable slug (e.g. "cc225") |
| name | string | Human-readable |
| price | int64 | ARS cents; represents current wholesale price |
| emoji | string | UI display |
| providers | []string | Provider IDs that carry this product |

### Provider
A wholesale source (Coto, Diarco, Makro, etc.). Products belong to one or more providers.

### Branch
A physical StockYa location that can receive warehouse-held stock via delivery.

---

## Invariants

These must always hold after any mutation:

1. **Credit limit**: `business.creditUsed ≤ business.creditLimit`
   - Enforced at advance creation (currently unchecked — see known gaps below).

2. **Paid is terminal**: once `advance.status = paid`, no further mutations are allowed.
   - Attempting to pay an already-paid advance returns a 400 error.

3. **Delivery requires warehouse**: `requestDelivery` only applies when `advance.logistics = deposito`.
   - Direct-ship advances (`envio`) cannot request additional delivery.

4. **Single delivery request**: once `deliveryStatus = en_transito`, the same advance cannot request delivery again.

5. **Non-empty cart**: an advance must have at least one item with `qty > 0`; an empty cart returns a 400 error.

6. **Price lock**: `AdvanceItem.price` is set at creation time from the current catalog price.
   Subsequent product price changes do not retroactively affect existing advances.

---

## Business Rules

### Storage cost
When `logistics = deposito`, a warehousing fee of **2% of the product subtotal** is added to `advance.total`.

```
storageCost = Σ(item.price × item.qty) × 0.02
total       = Σ(item.price × item.qty) + storageCost
```

### Credit accounting
- On advance creation: `business.creditUsed += advance.total`
- On advance payment: `business.creditUsed -= advance.total` (floor at 0)

### Due date
Always `adelantoDate + 30 calendar days`. Status transitions (`ok → warn → danger`) are computed from this date but are currently managed client-side.

---

## Known Gaps (spec vs. current implementation)

| Gap | Description |
|---|---|
| Credit limit enforcement | `createAdvance` does not yet reject when `creditUsed + total > creditLimit`. Tracked as a TODO. |
| Status aging | `advance.status` is set to `ok` at creation and never auto-updated to `warn`/`danger` as the due date passes. A background job or on-read computation is needed. |
| Persistence | All state is in-memory; a server restart loses all data. |
| ML endpoints | `/ml/*` return stub responses; behavioral scoring is not yet implemented. |
