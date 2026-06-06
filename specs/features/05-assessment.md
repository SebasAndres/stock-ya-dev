# Feature: Assessment Pre-Crédito

Antes de acceder a la línea de crédito, cada tienda debe completar un assessment de verificación
de tres pasos: foto del local, preguntas de perfil y captura de ubicación GPS.

**API:** `POST /api/businesses/{id}/assessment`
**Domain spec:** see `specs/domain.md § Business`

---

## Contexto de negocio

El assessment sirve como KYC mínimo para StockYa. Tres señales combinadas reducen el riesgo:
- **Foto del local** — evidencia de que el negocio existe físicamente.
- **Antigüedad + volumen diario** — proxies rápidos de estabilidad y capacidad de pago.
- **Ubicación GPS** — geolocalización del negocio para verificación y routing logístico futuro.

---

## Flujo completo

```
Registro → view-assessment (paso 1: foto) → (paso 2: preguntas + ubicación) → (paso 3: confirmación) → dashboard
```

Si `assessmentStatus == "approved"`, el wizard se saltea automáticamente al entrar.

---

## Scenarios

### Scenario: Primera sesión post-registro

**Given** una tienda recién registrada con `assessmentStatus = "pending"`
**When** el usuario entra al dashboard
**Then** se redirige a `view-assessment` (paso 1)
**And** el botón "Adelantar stock" no es accesible hasta completar el assessment

---

### Scenario: Completar assessment exitosamente

**Given** una tienda en paso 2 del assessment
**When** el usuario envía:
```json
POST /api/businesses/{id}/assessment
{
  "photo": "<base64-jpeg>",
  "yearsOpen": 3,
  "dailyCustomers": 80,
  "latitude": -34.6037,
  "longitude": -58.3816
}
```
**Then** la respuesta es `200 OK` con el business actualizado
**And** `assessmentStatus = "approved"`
**And** el frontend navega al dashboard con crédito habilitado

---

### Scenario: Assessment ya completado

**Given** una tienda con `assessmentStatus = "approved"`
**When** el usuario inicia sesión (o refresca)
**Then** el wizard no se muestra; va directo al dashboard

---

### Scenario: Campos requeridos faltantes

**Given** un POST sin `yearsOpen` o sin `dailyCustomers`
**Then** la respuesta es `400 Bad Request`
```json
{ "error": "yearsOpen and dailyCustomers are required" }
```

---

### Scenario: Geolocalización denegada

**Given** el usuario deniega acceso a la ubicación del navegador
**When** está en el paso de preguntas
**Then** el campo de ubicación muestra "Ubicación no disponible"
**And** el assessment igualmente puede completarse (`latitude` y `longitude` quedan en 0)

---

## Acceptance criteria

- [x] `assessmentStatus` arranca en `"pending"` al crear el negocio
- [x] `POST /api/businesses/{id}/assessment` setea `assessmentStatus = "approved"` inmediatamente
- [x] El frontend bloquea el acceso a crédito si `assessmentStatus != "approved"`
- [x] La foto se envía como base64 en el body; el backend la almacena en memoria
- [x] La ubicación se captura vía `navigator.geolocation`; si se deniega, el assessment igual procede
- [ ] En producción: reemplazar aprobación automática por revisión manual o scoring ML (ver `specs/features/04-ml-scoring.md`)
