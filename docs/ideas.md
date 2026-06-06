# Ideas

## Idea 1 - StockYa: Microcréditos a almacenes para mantener stock

**Problema**:
Los almacenes y pymes minoristas tienen un problema de cash flow cíclico: necesitan reponer stock antes de haber cobrado las ventas anteriores. Sin historial crediticio formal, quedan excluidos del sistema financiero tradicional.

**Solución**:
Plataforma B2B con crédito embebido + suscripción. El almacén hace su pedido a mayoristas a través de la app, que otorga un microcrédito atado a esa compra específica (BNPL B2B). El historial de comportamiento real (compras, pagos, suscripción) actúa como scoring propietario.

### Flujo

~~~
  [Almacén (pyme)]
       |
       | 1. Se suscribe + pide crédito para compra
       v
  +-----------+
  |  StockYa  |  <-- scoring basado en comportamiento
  +-----------+
       |                        |
       | 2. Orden de compra      | 3. Crédito desembolsado
       v                        v
  [Mayorista]            [Almacén recibe stock]
       |
       | 4. Despacho logístico (opcional)
       v
  [Depósito / Almacén]
~~~

### Modelo de suscripción

| Tier        | Precio/mes | Línea de crédito | Features                                      |
|-------------|------------|-----------------|-----------------------------------------------|
| Básico      | Gratis     | Sin crédito      | Marketplace mayoristas, precios en tiempo real |
| Crecimiento | $X         | Hasta $50k       | + Crédito, analytics de stock básico           |
| Pro         | $2X        | Hasta $200k      | + Mejor tasa, predicción de demanda            |
| Logística   | $3X        | Hasta $500k      | + Gestión de depósito terciarizada             |

### Revenue streams

1. **Suscripciones** — ingreso predecible
2. **Spread sobre créditos** — escala con volumen
3. **Comisión mayoristas** — ellos pagan por el volumen que les traemos
4. **Logística** — margen de servicio

### Diferenciador clave

> "No somos una fintech que le presta a almacenes. Somos la infraestructura financiera y logística de la cadena corta de distribución minorista."

---

## Análisis: Cold Start de microcréditos

El problema del cold start es que sin historial no hay scoring, y sin crédito no hay historial.

### Estrategia de cold start por capas

**Capa 1 — Señales de identidad y contexto (día 0)**
- CUIT/CUIL activo en AFIP (categoría monotributo o responsable inscripto)
- Antigüedad del local (habilitación municipal)
- Zona geográfica (índice de densidad comercial, nivel socioeconómico del barrio)
- Tipo de local (almacén, kiosco, verdulería) → ticket promedio esperable

**Capa 2 — Señales de comportamiento externas**
- Mercado Pago / Naranja X / Boa: si tiene historial de cobros con QR, se puede estimar facturación real
- Proveedores ya conocidos: si declara que compra en Makro/Diarco/Yaguar, se puede pedir un estado de cuenta
- AFIP: facturación declarada del último año (Monotributo categoría = límite de facturación)

**Capa 3 — Crédito escalonado con fricción intencional**

```
Semana 1-2:   Crédito mínimo ($5k-$10k), plazo 15 días, 1 producto
              → Si paga en término: sube automáticamente
Mes 1:        Crédito $20k-$50k, plazo 30 días, catálogo completo
Mes 3+:       Línea basada en comportamiento real, scoring propio
```

La fricción intencional (empezar chico) filtra a los malos pagadores antes de exponerse a montos grandes.

**Capa 4 — Scoring propio (t > 3 meses)**
Variables propietarias:
- Frecuencia de recompra
- Ticket promedio y varianza
- Pago en término (días de mora promedio)
- Categorías compradas (diversificación de stock)
- Engagement con la app (predictivo de continuidad)

### Límites del crédito en cold start

| Señal disponible               | Crédito máximo sugerido |
|-------------------------------|------------------------|
| Solo CUIT activo               | $5.000 - $10.000        |
| + Facturación AFIP             | Hasta 10% facturación mensual declarada |
| + Historial QR (Mercado Pago)  | Hasta 15% facturación estimada real     |
| + 1 mes de comportamiento      | Modelo propio          |

---

## Análisis: Compras en mayoristas en Argentina

### Principales mayoristas y su operatoria

| Mayorista        | Modalidad de compra       | Mínimo de compra | Forma de pago        |
|-----------------|--------------------------|-----------------|----------------------|
| **Makro**        | Autoservicio + e-commerce | ~$50k-$100k     | Contado / débito / crédito propio |
| **Diarco**       | Autoservicio              | Sin mínimo fijo  | Contado / débito     |
| **Yaguar**       | Autoservicio + reparto    | Por categoría   | Contado              |
| **Maxiconsumo**  | Autoservicio              | Sin mínimo fijo  | Contado / débito     |
| **Vital**        | Reparto + preventa        | Por zona        | Contado / cuenta corriente para clientes fidelizados |
| **Distribuidores de marcas** (Arcor, Molinos, etc.) | Preventa con preventista | Por zona | Cuenta corriente a 30/60 días |

### Fricciones actuales (oportunidades para StockYa)

1. **Pago en efectivo o débito**: la mayoría exige pago inmediato, bloqueando a quien no tiene liquidez
2. **Mínimos de compra**: un almacén chico no siempre llega al mínimo solo → posibilidad de compras agrupadas
3. **Sin precio online**: los precios cambian semanalmente (inflación), el almacenero no puede comparar fácil
4. **Logística propia**: el almacenero tiene que ir en remís o pagar flete, perdiendo medio día
5. **Sin historial como cliente**: los distribuidores no dan cuenta corriente a clientes nuevos o pequeños

### Cómo encaja StockYa en este ecosistema

- **Agregar demanda**: agrupar pedidos de varios almacenes para superar mínimos y negociar precio
- **Ser el comprador formal**: StockYa compra al mayorista (con su propia cuenta) y le vende al almacén con crédito → resuelve el problema de pago inmediato
- **Comparador de precios en tiempo real**: scraping / acuerdos con mayoristas para mostrar precios actualizados
- **Logística compartida**: un solo flete para varios almacenes del mismo barrio (milk run)

### Modelo legal/operativo sugerido

StockYa actúa como **revendedor** (compra y revende) o como **agente de compra con financiamiento**:
- Opción A (más simple): StockYa compra, factura al almacén, cobra con interés → margen + spread
- Opción B (más liviana): StockYa es intermediario, el mayorista factura directamente, StockYa garantiza el pago → comisión + interés sobre el crédito garantizado

---

## Análisis: StockYa como buffer de inventario (dos perfiles de cliente)

### El insight

No todos los almacenes tienen el mismo problema. Hay dos perfiles bien distintos que pueden coexistir en la plataforma:

| Perfil | Situación | Necesidad |
|--------|-----------|-----------|
| **Comprador adelantado** | Tiene liquidez | Quiere precio de volumen sin comprar volumen él solo, y sin perder medio día yendo al mayorista |
| **Comprador financiado** | Sin liquidez | Necesita el stock ahora, paga en 30/60 días |

Ambos perfiles se benefician del mismo mecanismo: **StockYa agrega demanda y actúa como buffer de inventario**.

### Cómo funciona el buffer

~~~
[Comprador adelantado] --paga $X por adelantado-->  +------------+
                                                    |            |
[Comprador financiado] --solicita crédito $Y------> | StockYa    |
                                                    | (buffer)   |
                                                    +------------+
                                                         |
                                          compra $X+$Y en bulk al mayorista
                                          (descuento por volumen agregado)
                                                         |
                                    +---------+----------+
                                    |                    |
                             entrega inmediata     entrega inmediata
                             al adelantado         al financiado
                             (ya pagó)             (paga en 30/60d)
~~~

### Por qué esto es mejor que un modelo puro de crédito

1. **El stock ya comprado es garantía implícita**: el crédito no es un préstamo en efectivo, es mercadería ya adquirida. El riesgo de default no destruye el capital, destruye el stock → el almacén necesita la mercadería para vender, su incentivo a pagar es directo.

2. **El adelantado financia al financiado**: los pagos por adelantado reducen el capital propio que StockYa necesita inmovilizar para los créditos → menor costo financiero, mayor margen.

3. **Descuento por volumen real**: cuanta más demanda agregada (adelantados + financiados), mayor poder de negociación ante el mayorista → precio más bajo → más margen o precio más competitivo para los almacenes.

4. **Rotación predecible**: si los adelantados hacen pedidos recurrentes (suscripción semanal/mensual), StockYa puede planificar compras al mayorista con anticipación → elimina el riesgo de sobrestock.

### Modelo de precio diferenciado por perfil

| Acción | Adelantado | Financiado |
|--------|-----------|------------|
| Pago | Por adelantado (antes de entrega) | A 30/60 días post-entrega |
| Precio del producto | Precio mayorista + margen mínimo | Precio mayorista + margen + interés |
| Suscripción | Tier básico/crecimiento | Tier crecimiento/pro |
| Beneficio extra | Prioridad en stock, precio fijo garantizado | Acceso al capital, sin necesitar liquidez |

### Gestión del buffer: cuánto stock mantener

El buffer no necesita ser un depósito propio desde el día 1. Estrategia por etapas:

```
MVP (0-3 meses):   Sin stock físico. StockYa compra bajo pedido confirmado.
                   Adelantados: pagan, StockYa compra al día siguiente, entrega en 48hs.
                   Financiados: mismo flujo, entrega en 48hs, crédito activado al confirmar compra.

Etapa 2 (3-6m):   Mini-buffer de productos de alta rotación (arroz, aceite, azúcar, yerba).
                   Permite entrega en el día para esos SKUs.

Etapa 3 (6m+):    Buffer dinámico basado en predicción de demanda de la red de almacenes.
```

### Riesgo principal del modelo buffer

- **Inflación**: si el stock queda parado más de 1 semana en Argentina, el precio de reposición sube y el margen se erosiona.
  - Mitigación: buffer solo en productos de alta rotación, precios actualizados semanalmente, los adelantados "bloquean" precio al momento del pago.
- **Sobrestock por cancelaciones**: el adelantado pagó, no puede cancelar sin penalidad.
  - Mitigación: política clara de no-devolución en compras prepagadas (se puede ofrecer crédito para próxima compra en lugar de reembolso).
