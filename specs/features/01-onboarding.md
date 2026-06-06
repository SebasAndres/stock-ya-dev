# Feature: Business Onboarding

A new store registers with StockYa to gain access to the catalog and a credit line.

**API:** `POST /api/businesses`
**Domain spec:** see `specs/domain.md § Business`

---

## Scenario: Successful registration

**Given** a new almacén owner with a valid CUIT  
**When** they POST to `/api/businesses` with `name`, `cuit`, `type`, `phone`, `address`  
**Then** the response is `201 Created` with the full `Business` object  
**And** `creditLimit = 50000` (default starting line)  
**And** `creditUsed = 0`

---

## Scenario: Missing required fields

**Given** a POST body without `name` or without `cuit`  
**When** the request is received  
**Then** the response is `400 Bad Request` with `{ "error": "name and cuit are required" }`

---

## Scenario: Retrieve registered business

**Given** a business was registered and its `id` returned  
**When** `GET /api/businesses/{id}` is called  
**Then** the response is `200 OK` with the same business record  

**When** an unknown `id` is used  
**Then** the response is `404 Not Found`

---

## Acceptance criteria

- [ ] `id` is unique and stable across calls
- [ ] `createdAt` is set to the server time at registration
- [ ] `type` defaults gracefully when omitted (field is optional)
- [ ] Credit limit enforcement gap is documented in `domain.md` until resolved
