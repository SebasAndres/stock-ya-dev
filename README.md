# StockYa

Sistema de tokenización de stock físico para financiamiento de kioscos y almacenes, inspirado en el modelo de *warehouse receipt financing* (similar a los certificados de depósito de soja en silos argentinos).

---

## El Problema

Los kioscos y almacenes enfrentan dos fricciones opuestas:

- **Sin capital:** necesitan stock pero no tienen liquidez para comprarlo.
- **Con capital pero sin espacio:** tienen pesos que se licúan con la inflación pero no pueden acumular mercadería.

---

## El Modelo

**XXXXX** actúa como depósito físico centralizado y emisor de tokens 1:1 respaldados por stock real.

```
Stock físico (Coca-Cola, etc.)  →  Token ERC-20 / digital  →  Circula entre kioscos
```

### Casos de uso

| Actor | Acción | Beneficio |
|---|---|---|
| Kiosco sin capital | Toma tokens como crédito → retira físico → paga en pesos + interés | Acceso a stock sin liquidez inicial |
| Kiosco con capital, sin espacio | Compra tokens → los deja en depósito → cobra interés | Plazo fijo en mercadería, cobertura inflacionaria |

### Garantía anti-default

Se entrega el stock completo **menos la última unidad**, que se libera recién cuando el pago está completo. Esto alinea incentivos sin necesidad de garantías externas.

---

## Referencia del mundo real

Modelo análogo al **warehouse receipt financing** — la misma mecánica que usan los productores de soja en Argentina para financiarse con sus granos en silo sin venderlos.

---
