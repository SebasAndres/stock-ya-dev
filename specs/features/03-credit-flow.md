# Feature: Credit Advance Lifecycle

The core financial flow: a store draws stock on credit, chooses how it's held, optionally requests delivery, then pays.

**APIs:**
- `POST /api/businesses/{id}/advances` — create advance
- `GET /api/businesses/{id}/advances` — list advances
- `POST /api/advances/{id}/pay` — settle payment
- `POST /api/advances/{id}/delivery` — request delivery for warehouse-held stock
- `GET /api/businesses/{id}/dashboard` — summary stats

---

## Scenario: Direct-ship advance (envio)

**Given** a registered business with available credit  
**When** they POST `{ cart: {cc225: 10}, logistics: "envio" }`  
**Then** an advance is created with:
  - `status = ok`
  - `deliveryStatus = ""`
  - `depositoCost = 0`
  - `total = 10 × 2800 = 28000`
  - `dueDate = adelantoDate + 30 days`
**And** `business.creditUsed` increases by `28000`

---

## Scenario: Warehouse advance with storage fee (deposito)

**Given** the same cart but `logistics: "deposito"`  
**Then** the advance is created with:
  - `depositoCost = 28000 × 0.02 = 560`
  - `total = 28000 + 560 = 28560`
  - `deliveryStatus = ""` (delivery not yet requested)

---

## Scenario: Request delivery for a warehouse advance

**Given** a warehouse advance exists with `deliveryStatus = ""`  
**When** `POST /api/advances/{id}/delivery` with `{ sucursalId: "palermo" }` is called  
**Then** the advance is updated with:
  - `deliveryStatus = "en_transito"`
  - `targetSucursal = "Sucursal Palermo"`

**When** delivery is requested again on the same advance  
**Then** the response is `400 Bad Request` ("delivery already in transit")

---

## Scenario: Delivery not applicable for direct-ship advances

**Given** an advance with `logistics = "envio"`  
**When** `POST /api/advances/{id}/delivery` is called  
**Then** the response is `400 Bad Request` ("advance is not stored in warehouse")

---

## Scenario: Pay an advance

**Given** an unpaid advance with `total = 28000`  
**When** `POST /api/advances/{id}/pay` is called  
**Then**:
  - `advance.status = "paid"`
  - `business.creditUsed` decreases by `28000` (floor at 0)

**When** the same advance is paid again  
**Then** the response is `400 Bad Request` ("advance already paid")

---

## Scenario: Empty cart is rejected

**When** `POST /api/businesses/{id}/advances` with `cart: {}` or all quantities 0  
**Then** the response is `400 Bad Request` ("cart is empty")

---

## Scenario: Unknown product in cart

**When** the cart contains a product ID that does not exist in the catalog  
**Then** the response is `400 Bad Request` ("product {id} not found")

---

## Scenario: Dashboard reflects current state

**Given** a business with 2 advances (1 paid, 1 active with 3 items)  
**When** `GET /api/businesses/{id}/dashboard` is called  
**Then**:
  - `totalAdelantos = 2`
  - `totalPagados = 1`
  - `activeProducts = 3`

---

## Acceptance criteria

- [ ] `creditUsed` is never negative after payment
- [ ] Paid advances are immutable (no re-payment, no delivery update)
- [ ] All monetary values in ARS cents (int64)
- [ ] Due date is always `adelantoDate + 30 calendar days`
- [ ] Status aging (ok → warn → danger) is a tracked gap — see `domain.md`
