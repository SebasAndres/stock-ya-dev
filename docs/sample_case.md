# Caso de Uso: StockYa — Lanzamiento y Primeros Ciclos de Crédito

Este documento describe el flujo operativo y financiero del primer mes de StockYa a través de cuatro almacenes iniciales. Su objetivo es validar la viabilidad del modelo antes de implementar los endpoints de scoring y logística.

---

## 1. Parámetros del modelo (Mes 0)

| Símbolo | Valor | Descripción |
|---------|-------|-------------|
| **Z**   | $2.000.000 ARS/mes | Capital operacional mensual total |
| **A**   | $1.400.000 ARS (70 % de Z) | Pool destinado a microcréditos |
| **B**   | $300.000 ARS (15 % de Z) | Infraestructura de datos y desarrollo del modelo de scoring |
| **C**   | $300.000 ARS (15 % de Z) | Logística inicial: depósito transitorio + fletes |
| **m₀**  | 10 clientes | Almacenes/kioscos semilla, seleccionados por red de confianza |
| **d**   | 0,85 | Fracción esperada de clientes que repagan en término (mes 0) |
| **t**   | 4,5 % TEM | Tasa mensual que paga StockYa para financiarse (crédito pyme) |
| **t'**  | ≥ t + 0,1 % = 4,6 % TEM (mínimo contractual) | Tasa mensual cobrada al cliente; crece según scoring de riesgo |
| **g**   | 1,3 | Factor de crecimiento mensual: m_{i+1} = m_i × g |
| **fee_dep** | 2 % del subtotal de productos | Fee de logística y depósito en sucursal StockYa |
| **ref_bonus** | $10.000 ARS en límite de crédito | Beneficio al referenciador por cada cliente nuevo que complete su primer adelanto |

### Restricción de asignación

$$
Z = A + B + C \quad \Rightarrow \quad \$2.000.000 = \$1.400.000 + \$300.000 + \$300.000 \quad \checkmark
$$

### Línea de crédito inicial por cliente

Cada uno de los m₀ clientes accede a un **primer adelanto máximo de A/m₀**:

$$
\text{límite\_inicial} = \frac{A}{m_0} = \frac{\$1.400.000}{10} = \$140.000 \text{ ARS}
$$

El primer adelanto está limitado a **$50.000 ARS** por política de cold-start (fricción intencional). Este límite sube automáticamente tras el primer repago en término.

---

## 2. Restricciones formales

### 2.1 Viabilidad del libro de crédito

Para que el negocio de crédito sea autosustentable (sin considerar B y C, que son costos de infraestructura), la expectativa de cobranza debe superar el costo de fondeo:

$$
\begin{align*}
d \times A \times t' &> t \times A \\
\implies d \times t' &> t \\
\implies 0{,}85 \times t' &> 0{,}045 \\
\implies t' &> 5{,}29\% \text{ TEM}
\end{align*}
$$

**Conclusión:** La tasa mínima rentable para el portafolio agragado es **t' > 5,3 % TEM**. Todo cliente con score que implique una tasa < 5,3 % TEM debe rechazarse o reclasificarse.

La tasa de los clientes iniciales sin historial se fija en **t'₀ = 6 % TEM**, generando un margen bruto de:

$$
\begin{align*}
\text{margen\_bruto} &= d \times A \times (t' - t) \\
&= 0{,}85 \times \$1.400.000 \times (0{,}06 - 0{,}045) \\
&= 0{,}85 \times \$1.400.000 \times 0{,}015 \\
&= \$17.850 \text{ ARS / mes}
\end{align*}
$$

### 2.2 Restricción de cobertura con fondeo externo

Si los adelantos superan el pool A (por ejemplo en un mes de fuerte demanda), StockYa puede financiarse externamente a tasa t. La condición de solvencia es:

$$
\begin{align*}
\text{ingreso\_esperado} &\geq \text{costo\_externo} \\
d \times \text{cartera\_total} \times t' &\geq t \times \text{exceso\_financiado}
\end{align*}
$$

donde $\text{exceso\_financiado} = \max(0,\ \text{cartera\_total} - A)$.

### 2.3 Crecimiento de la red

El número de clientes crece exponencialmente con tasa mensual $g = 1{,}3$:

$$
m_i = m_0 \cdot g^i = 10 \cdot 1{,}3^i
$$

$$
m_0 = 10, \quad m_1 = 13, \quad m_2 = 17, \quad m_3 = 22, \quad \ldots
$$

La tasa de crecimiento continua equivalente es $r = \ln(g) = \ln(1{,}3) \approx 0{,}262$ por mes.

La cartera de crédito crece proporcionalmente si el crédito promedio por cliente se mantiene constante.

---

## 3. Flujo narrativo

### Semana 1 — Almacén #1: Primera operación

**Contexto.** Almacén #1 es un almacén familiar de barrio en Flores. Quiere comprar $40.000 ARS en Coca-Cola 2.25L porque anticipa demanda alta y el precio mayorista está conveniente. No tiene liquidez inmediata y tampoco tiene espacio de almacenamiento propio.

**Registro y onboarding.**

1. Almacén #1 crea cuenta en StockYa con CUIT activo (monotributista categoría D).
2. El sistema asigna una línea de crédito inicial de $50.000 ARS (primer adelanto, cold-start cap).
3. Scoring inicial: sin historial propio → tasa t'₁ = **6,0 % TEM** (piso para clientes nuevos).

**Solicitud del adelanto.**

| Ítem | Detalle |
|------|---------|
| Producto | Coca-Cola 2.25L — 16 unidades |
| Precio mayorista unitario | $2.500 ARS |
| Subtotal productos | $40.000 ARS |
| Logística | `deposito` — retiro posterior desde sucursal |
| Fee depósito | $40.000 × 2 % = **$800 ARS** |
| **Total del adelanto** | **$40.800 ARS** |
| Tasa aplicada (t'₁) | 6,0 % TEM |
| Vencimiento | 30 días desde la fecha del adelanto |
| **Monto a repagar** | $40.800 × 1,06 = **$43.248 ARS** |

El fee de depósito ($800) forma parte del `advance.total` y queda bloqueado en la línea de crédito.

**Tracking de stock.**

Almacén #1 puede ver en la plataforma cuántas unidades tiene en depósito. Cuando las necesita, solicita la entrega a la sucursal más cercana (logística tercerizada). La entrega actualiza `deliveryStatus` a `en_transito`.

**Repago.**

Al cabo de 30 días, con las ventas de las Coca-Colas, Almacén #1 repaga **$43.248 ARS**. Su `advance.status` pasa a `paid` y su `creditUsed` baja a $0.

**Efecto en el portafolio StockYa (Mes 1, cliente único):**

$$
\begin{align*}
\text{Ingreso bruto} &= \$43.248 - \$40.800 = \$2.448 \text{ ARS} \\
\text{Costo de fondeo} &= \$40.800 \times 4{,}5\% = \$1.836 \text{ ARS} \\
\text{Margen neto} &= \$2.448 - \$1.836 = \$612 \text{ ARS}
\end{align*}
$$

---

### Semana 2 — Almacén #2: El comprador adelantado

**Contexto.** Almacén #2, también del barrio, tiene mejor caja que #1. Vio cómo le fue a #1 y entiende que puede usar StockYa para comprar a precio de volumen sin desembolsar todo junto ni ir al mayorista.

**Perfil: comprador adelantado (no financiado).** Almacén #2 paga por adelantado y no necesita crédito: usa la plataforma como agregador de demanda y buffer de inventario.

**Operación.**

| Ítem | Detalle |
|------|---------|
| Producto | Fideos 500g — 100 paquetes |
| Precio mayorista unitario | $1.000 ARS |
| Subtotal | $100.000 ARS |
| Logística | `deposito` — sin entrega inmediata |
| Pago | Inmediato (no es crédito) |
| Fee depósito | $100.000 × 2 % = $2.000 ARS |
| **Total pagado** | **$102.000 ARS** |

StockYa recibe $102.000 ARS de contado, compra los fideos al mayorista, y los almacena. Almacén #2 retira cuando quiere.

**Referidos.** Almacén #2, satisfecho con el servicio, recomienda la plataforma a Almacén #3 y Almacén #4. Al completar el primer adelanto de cada referido, Almacén #2 recibe **+$10.000 ARS en línea de crédito por cada uno** (total potencial: +$20.000 ARS).

**Efecto en el grafo de confianza:**

```
[#1] ──referido por red inicial──▶ [StockYa seed]
[#2] ──referido por #1──────────▶ [StockYa seed]
[#3] ──referido por #2──────────▶ nodo nuevo
[#4] ──referido por #2──────────▶ nodo nuevo
```

StockYa registra internamente esta estructura como un grafo dirigido temporal. Los atributos del nodo incluyen: tipo de cliente (adelantado / financiado), historial de pagos, ticket promedio, categorías de stock.

---

### Semana 3 — Almacén #3: Entrega urgente desde el buffer

**Contexto.** Almacén #3 se quedó sin stock de fideos. Necesita reposición urgente y, a diferencia de #2, no tiene caja. Viene referido por #2.

**Ventaja del scoring por referido.** Al llegar referido por un cliente con buen comportamiento (Almacén #2 pagó en término y es adelantado), el modelo aplica una leve reducción de tasa:

| Factor de scoring | Ajuste |
|---|---|
| Sin historial propio | +1,5 % sobre t |
| Referido por cliente en buen estado | -0,2 % |
| CUIT activo > 2 años | -0,3 % |
| **$t'_3$ resultante** | $4{,}5\% + 1{,}5\% - 0{,}2\% - 0{,}3\% = 5{,}5\%$ **TEM** |

**Operación.**

| Ítem | Detalle |
|------|---------|
| Producto | Fideos 500g — 10 paquetes |
| Precio | $1.000 ARS/u |
| Subtotal | $10.000 ARS |
| Logística | `envio` — entrega inmediata a su local |
| Fee | $0 (logística directa, sin depósito) |
| **Total del adelanto** | **$10.000 ARS** |
| Tasa (t'₃) | 5,5 % TEM |
| **Monto a repagar** | $10.000 × 1,055 = **$10.550 ARS** |

**Ejecución del buffer.** StockYa tiene en su depósito los 100 paquetes de fideos de Almacén #2. Entrega 10 paquetes de ese stock a Almacén #3 y emite una orden de recompra al mayorista para reponer el depósito. Almacén #2 no se ve afectado — su stock en depósito sigue siendo 100 unidades.

Este es el primer caso de reutilización del buffer: el pago adelantado de #2 financia la liquidez de #3 sin que StockYa deba inmovilizar capital adicional.

---

### Semana 3-4 — Almacén #4: Cuenta activa, sin adelanto inmediato

Almacén #4 completa el onboarding pero no solicita un adelanto en la primera semana. Su existencia en la plataforma ya es útil: aparece en el grafo como nodo referido por #2, y sus metadatos iniciales (CUIT, categoría monotributo, zona geográfica) alimentan el modelo de scoring antes de que genere actividad.

---

## 4. Grafo de confianza (red de referidos)

```
Periodo i=0 (seed)         Periodo i=1              Periodo i=2
                                                     
  [StockYa]                [StockYa]                [StockYa]
     │                        │                        │
     ├─▶ [#1]                 ├─▶ [#1] ──ref──▶ [#2]  ├─▶ [#1]
     │                        │                   │    ├─▶ [#2]──ref──▶[#3]
     │                        │                   └──ref──▶[#3]       [#4]
     │                        │                   └──ref──▶[#4]
```

**Propiedades del grafo:**

| Elemento | Significado |
|---|---|
| Nodo | Almacén/kiosco registrado |
| Arista dirigida A→B | A refirió a B |
| Peso temporal | Días desde la referencia |
| Atributos de nodo | historial_pagos, ticket_promedio, categorías_stock, días_mora_promedio, engagement |

**Uso en scoring.** El modelo de scoring (ver `specs/features/04-ml-scoring.md`) usa este grafo para propagar señales de confianza: si A tiene buen historial y refirió a B, B hereda una fracción del score de A como prior. A medida que B genera su propio historial, el peso del prior se reduce y el comportamiento propio domina.

---

## 5. Estado del sistema al cierre del Mes 1

| Métrica | Valor |
|---|---|
| Clientes activos | 4 (de 10 onboardeados) |
| Adelantos emitidos | 2 (Almacén #1 y #3) |
| Adelantos pagados en término | 1 (#1, si repagó en semana 4) |
| Crédito vivo (creditUsed total) | $50.800 ARS (#3 pendiente) |
| Capital pool utilizado | $50.800 / $1.400.000 = 3,6 % |
| Ingresos por intereses (realizados) | $2.448 ARS (#1) |
| Ingresos por fee de depósito | $2.000 (#2) + $800 (#1) = $2.800 ARS |
| Aristas en el grafo de referidos | 3 (#1→#2, #2→#3, #2→#4) |

El bajo uso del pool (3,6 %) es esperado: el Mes 1 es de validación con clientes conocidos. El objetivo es llegar al Mes 3 con una utilización del 40-60 % y un d observado ≥ 0,80.

---

## 6. Proyección de crecimiento (primeros 4 meses)

| Mes | Clientes (m_i) | Pool utilizado (estimado) | Ingreso bruto mensual (estimado) |
|-----|---------------|--------------------------|----------------------------------|
| 0   | 10 onboardeados / 4 activos | $50.800 | $5.248 |
| 1   | 13 | ~$200.000 | ~$18.000 |
| 2   | 17 | ~$500.000 | ~$45.000 |
| 3   | 22 | ~$900.000 | ~$81.000 |

A partir del Mes 3, el modelo de scoring propio reemplaza los valores por defecto y permite diferenciar tasas con mayor precisión, comprimiendo el riesgo de default por debajo de d = 0,85.

---

## 7. Gaps y supuestos simplificados

| Ítem | Supuesto actual | Lo que falta |
|---|---|---|
| Precio de productos | Fijo durante el ciclo de 30 días (price lock) | Ajuste por inflación: en Argentina, el precio mayorista puede moverse ±5 % en 30 días |
| Logística | Tercerizada a precio fijo | Acuerdos con proveedores de logística aún no negociados |
| Scoring | Reglas manuales + prior de referido | Modelo ML completo (ver `04-ml-scoring.md`) |
| Capital externo | StockYa no necesita fondeo externo en Mes 1 | A partir del Mes 3, la cartera puede superar A y requerir línea de crédito pyme |
| Garantía anti-default | Última unidad retenida en depósito hasta pago completo | Pendiente de implementar en el endpoint de delivery |
