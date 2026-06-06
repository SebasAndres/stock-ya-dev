# Feature: Catalog Browsing

Store owners browse available wholesale products before placing an advance.

**APIs:** `GET /api/products`, `GET /api/providers`, `GET /api/branches`

---

## Scenario: Browse all products

**When** `GET /api/products` is called without a query parameter  
**Then** all products in the catalog are returned  
**And** each product includes `id`, `name`, `price`, `emoji`, and `providers[]`

---

## Scenario: Filter by provider

**When** `GET /api/products?provider=coto` is called  
**Then** only products whose `providers` array includes `"coto"` are returned

**When** `GET /api/products?provider=todos` is called  
**Then** all products are returned (same as no filter)

---

## Scenario: List providers

**When** `GET /api/providers` is called  
**Then** all providers are returned with their display metadata (`type`, `src`, `text`, `badge`)

---

## Scenario: List delivery branches

**When** `GET /api/branches` is called  
**Then** all StockYa branches are returned with `id`, `name`, `address`, `emoji`

---

## Notes

- Prices are in ARS cents (integer). The frontend formats them as `$X.XXX`.
- The catalog is static seed data in MVP. Prices do not change at runtime.
- Provider `"todos"` is a UI sentinel — it is not a real supplier, just an all-filter option.
