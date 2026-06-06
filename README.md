# StockYa

<img src="demo.png">

Plataforma B2B de crédito embebido para almacenes y kioscos argentinos. Los negocios piden adelantos de stock contra una línea de crédito y pagan a 30 días — sin ir al banco, sin garantías externas.

Inspirado en el modelo de *warehouse receipt financing* (los certificados de depósito de soja en silos argentinos), pero para mercadería de consumo masivo.

---

## El Problema

Los almacenes y kioscos argentinos enfrentan dos fricciones opuestas:

- **Sin liquidez:** necesitan reponer stock antes de haber cobrado las ventas anteriores.
- **Con liquidez pero sin espacio:** tienen pesos que se licúan con la inflación pero no pueden acumular mercadería.

Sin historial crediticio formal, quedan excluidos del sistema financiero tradicional.

---

## El Modelo

StockYa actúa como **buffer de inventario y emisor de crédito embebido**:

```
[Almacén] ── pide adelanto de stock ──▶ [StockYa]
                                              │
                              compra al mayorista (Coto, Diarco, Makro…)
                                              │
                   ┌──────────────────────────┴──────────────────────────┐
                   ▼                                                      ▼
           envio: entrega directa                          deposito: retiro desde
           al local (sin fee)                             sucursal StockYa (+2%)
```

| Actor | Situación | Qué hace StockYa |
|---|---|---|
| Almacén sin liquidez | Necesita stock ya | Otorga adelanto; cobra tasa a 30 días |
| Almacén con liquidez | Quiere precio de volumen sin ir al mayorista | Agrega demanda; entrega en 48 hs |

### Garantía anti-default

Se entrega el stock completo **menos la última unidad**, que se libera recién cuando el pago está completo. Esto alinea incentivos sin garantías externas.

### Red de referidos y scoring

Cada cliente que refiere a un nuevo almacén alimenta un **grafo de confianza** que el motor de scoring usa para asignar tasas. Un cliente referido por alguien con buen historial hereda un prior positivo — resolviendo el problema de cold start.

---

## Stack

| Capa | Tecnología |
|---|---|
| Backend | Go — `net/http` estándar, sin frameworks |
| Frontend | HTML + JS vanilla, una sola página |
| Persistencia | In-memory (MVP) |
| API | REST, contrato en `specs/openapi.yaml` |

---

## Correr localmente

```bash
cd backend
go run .
# Servidor en http://localhost:8080
```

El frontend se sirve estático desde el mismo proceso en `/`.

---

## API

El contrato completo está en [`specs/openapi.yaml`](specs/openapi.yaml). Endpoints principales:

| Método | Path | Descripción |
|---|---|---|
| `POST` | `/api/businesses` | Registrar un negocio |
| `GET` | `/api/businesses/{id}/dashboard` | Stats y crédito disponible |
| `POST` | `/api/businesses/{id}/advances` | Crear un adelanto de stock |
| `GET` | `/api/businesses/{id}/advances` | Listar adelantos |
| `POST` | `/api/advances/{id}/pay` | Registrar repago |
| `POST` | `/api/advances/{id}/delivery` | Pedir entrega desde depósito |
| `GET` | `/api/products` | Catálogo de productos |
| `GET` | `/api/providers` | Mayoristas disponibles |
| `GET` | `/api/branches` | Sucursales de retiro |
| `POST` | `/api/ml/credit-score` | Score crediticio *(stub)* |
| `POST` | `/api/ml/demand-forecast` | Predicción de demanda *(stub)* |
| `POST` | `/api/ml/overdue-risk` | Riesgo de default *(stub)* |

---

## Estructura del proyecto

```
backend/        Go — modelos, handlers, store in-memory, stubs ML
frontend/       HTML + assets (imágenes de productos)
specs/
  openapi.yaml          Contrato de la API (fuente de verdad)
  domain.md             Entidades, invariantes, reglas de negocio
  features/
    01-onboarding.md
    02-catalog.md
    03-credit-flow.md
    04-ml-scoring.md    Spec completo; implementación pendiente
docs/
  sample_case.md        Caso de uso numérico: 4 almacenes, Mes 1
  ideas.md              Análisis de mercado y decisiones de diseño
```

---

## Estado de implementación

| Feature | Estado |
|---|---|
| Registro de negocios | Implementado |
| Catálogo de productos y mayoristas | Implementado |
| Ciclo completo de adelantos (crear / pagar / pedir entrega) | Implementado |
| Logística dual (`envio` / `deposito`) | Implementado |
| ML scoring y forecasting | Stub — spec en `specs/features/04-ml-scoring.md` |
| Persistencia en base de datos | Pendiente (hoy: in-memory) |
| Enforcement del límite de crédito | Pendiente (known gap en `specs/domain.md`) |
| Aging de status (`ok → warn → danger`) | Pendiente |
