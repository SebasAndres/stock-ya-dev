# Contexto del proyecto

**StockYa** — plataforma B2B de crédito embebido para almacenes y kioscos argentinos.
Los negocios piden adelantos de stock contra una línea de crédito; pagan a 30 días.

## Dominio

- **Actores:** StockYa (operador del buffer de inventario), kioscos/almacenes (usuarios de crédito o ahorro en mercadería)
- **Mecanismo central:** adelantos de stock respaldados por compras reales a mayoristas (Coto, Diarco, Makro, etc.)
- **Analogía de referencia:** warehouse receipt financing / BNPL B2B

## Decisiones de diseño tomadas

- La garantía anti-default retiene la última unidad hasta completar el pago — no se usa garantía externa
- Los tokens representan unidades físicas concretas (e.g. 1 token = 1 Coca-Cola 2.25L), no valor monetario
- Logística doble: envío directo (`envio`) o depósito en sucursal + retiro posterior (`deposito`, con fee 2%)

---

## Workflow: Spec-Driven Development

Este proyecto usa spec-driven development. Antes de tocar código, se escribe (o actualiza) el spec.

### Reglas

1. **Spec primero.** Antes de implementar cualquier feature o endpoint nuevo, actualizar o crear el spec correspondiente en `specs/`.
2. **API contract.** Todo endpoint nuevo va primero en `specs/openapi.yaml`. El YAML es fuente de verdad; el código implementa lo que dice.
3. **Feature spec.** Cualquier comportamiento no trivial tiene un archivo en `specs/features/`. Usar escenarios (Given/When/Then) para los casos felices y edge cases.
4. **Domain invariants.** Cualquier regla de negocio nueva va en `specs/domain.md` antes de implementarse.
5. **Si el código difiere del spec**, el spec se actualiza primero con la justificación, luego el código.

### Estructura de `specs/`

```
specs/
  openapi.yaml              # Contrato completo de la API (source of truth)
  domain.md                 # Entidades, invariantes, reglas de negocio
  features/
    01-onboarding.md        # Registro de negocios
    02-catalog.md           # Catálogo de productos y proveedores
    03-credit-flow.md       # Ciclo de vida del adelanto de crédito
    04-ml-scoring.md        # ML scoring y forecasting (spec completo, impl. pendiente)
```

### Cómo agregar una feature nueva

```
1. Escribir o actualizar el spec en specs/features/XX-nombre.md
2. Si agrega endpoints: actualizar specs/openapi.yaml
3. Si agrega entidades o reglas: actualizar specs/domain.md
4. Revisar el spec con el usuario
5. Implementar en backend/ y/o frontend/
6. Marcar acceptance criteria como completos en el feature spec
```

### Notas para Claude

- Cuando el usuario pide una feature nueva, **preguntar si quiere el spec primero** antes de escribir código.
- Al escribir specs: ser concreto con ejemplos de requests/responses reales (números de ARS reales, IDs de productos reales del catálogo).
- Los gaps conocidos van en `specs/domain.md § Known Gaps`, no como TODOs en el código.
