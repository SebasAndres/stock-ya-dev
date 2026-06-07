# StockYa

<img src="demo.png">

**Microcréditos productivos en stock para comercios de barrio.** Los almacenes y kioscos argentinos adelantan mercadería contra una línea de crédito y pagan a 30 días — sin ir al banco, sin historial crediticio formal, sin garantías materiales.

---

## El problema

El 45,3% de la población adulta argentina no tiene acceso a crédito formal de ningún tipo. 9,2 millones de personas trabajan en la informalidad (INDEC). Los que sí consiguen crédito vía fintechs lo pagan caro: en MercadoPago el costo financiero total puede llegar al 600% anual.

El mercado informal argentino mueve más de US$150.000 millones al año y opera prácticamente sin acceso al crédito productivo.

El nudo del problema es circular: para tener historial crediticio hay que haber pedido un crédito. Para que te den un crédito necesitás historial. **Los sectores más vulnerables están atrapados en ese bucle.**

---

## La solución

StockYa rompe ese bucle con tres ejes:

### 1. Microcréditos atados al stock

El crédito se ejecuta como compra inmediata de mercadería a través de nuestros partners mayoristas (Coto, Diarco, Makro). El dinero nunca pasa por el comerciante: va directo a la mercadería. Eso elimina el riesgo de fraude o desvío de fondos y garantiza el destino productivo del crédito.

```
[Almacén] ── pide adelanto ──▶ [StockYa]
                                     │
                     compra al mayorista (Coto, Diarco, Makro…)
                                     │
              ┌──────────────────────┴──────────────────────┐
              ▼                                              ▼
      envio: entrega directa                   deposito: retiro desde
      al local (sin fee)                       sucursal StockYa (+2%)
```

La garantía anti-default es la mercadería misma: se entrega el stock completo menos la última unidad, que se libera cuando el pago está completo. Sin garantías externas, sin avales.

### 2. ADN digital y social

Construimos un perfil crediticio a partir de lo que el sistema tradicional ignora: la ubicación del negocio, la red de proveedores, los contactos del comerciante, sus referidos y su documentación (facturas, recibos, remitos). Toda esa información forma el **ADN Social** — mientras más completo, mejores tasas.

El ADN Social se representa como un grafo de relaciones entre comerciantes conectados por ubicación, referencias y patrones de compra. Con cada transacción el grafo se vuelve más denso y más predictivo. El buen comportamiento de pago se propaga: un comerciante nuevo rodeado de buenos pagadores hereda esa confianza, lo que resuelve el problema de cold start y permite ofrecer mejor tasa desde el primer día. En sentido inverso, si una zona empieza a deteriorar sus pagos, el grafo lo detecta antes de que se traduzca en defaults individuales.

### 3. Group Lending

Reemplazamos el colateral de activos materiales por presión social. El mismo mecanismo usado en Bangladesh logró tasas de repago superiores al 97%. La red de pares es la garantía.

---

## Onboarding

El proceso es simple: el comerciante se registra, establece la ubicación de su negocio y el sistema valida sus datos con RENAPER y el Banco Central. A partir de ahí puede compartir contactos, seleccionar locales amigos, subir documentación y construir su ADN Social para acceder a mejores condiciones.

Para escalar la red inicial, arrancamos desde comunidades ya formadas: iglesias, comedores comunitarios, clubes de barrio y organizaciones sociales que ya tienen la información que necesitamos.

---

## Modelo de negocio

**Motor financiero:** capturamos el diferencial entre la tasa a la que conseguimos capital y la tasa a la que prestamos. Un scoring más preciso maximiza ese spread: al conocer mejor el riesgo real, prestamos con más confianza, reducimos la mora y ofrecemos mejores tasas sin resignar margen.

**Motor de datos:** el ADN digital construye información que ningún banco ni bureau crediticio tiene sobre este segmento. A escala, ese grafo es un activo con valor propio:
- Puede licenciarse como infraestructura de scoring a otras instituciones
- Puede monetizarse como datos agregados de consumo popular para distribuidores y marcas
- Se convierte en una barrera de entrada que ningún competidor puede replicar sin haber construido la misma red en el mismo territorio

---

## Fundraising

Buscamos **US$500.000** para un runway de 12 a 18 meses: construir el MVP, validar el scoring social, originar los primeros microcréditos productivos, medir mora/repago y demostrar tracción en comunidades piloto.

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
just dev
# Servidor en http://localhost:8080
```

El frontend se sirve estático desde el mismo proceso en `/`.

---

## Probar la app

Abrí `http://localhost:8080` en el navegador. El flujo completo es:

### 1. Registro

Completá el formulario con nombre del negocio, CUIT, teléfono, dirección y tipo (almacén o kiosco). Al enviar, el sistema corre una validación contra RENAPER y el Banco Central: si el CUIT está inhabilitado o tiene deudas en el BCRA, el acceso se bloquea o se muestra una advertencia con el límite de crédito reducido. Con un CUIT limpio, se asigna una línea de crédito inicial y entrás al dashboard.

> Para testear el camino feliz usá el CUIT `20-12345678-9`. Para ver el bloqueo por deuda BCRA usá `20-00000000-0`.

### 2. Dashboard

Muestra el crédito disponible, los adelantos activos y el historial. Desde acá accedés a las tres secciones principales por la barra inferior:

- **Adelantar** — pedir nuevo stock a crédito
- **Mi Depósito** — mercadería guardada en sucursales StockYa
- **Mi Stock** — adelantos activos y deuda pendiente

### 3. Pedir un adelanto (3 pasos)

**Paso 1 — Catálogo:** filtrá por mayorista (Coto, Diarco, Makro, Supermercado Chino) y sumá productos al carrito con los botones `+`/`−`. El total se actualiza en tiempo real en la barra superior.

**Paso 2 — Logística:** elegí cómo recibir la mercadería:
- *Envío directo* — llega al local, sin costo extra
- *Depósito en sucursal* — quedá guardada en una sucursal StockYa hasta que la retirés; tiene un fee del 2%

**Paso 3 — Confirmación:** revisá el resumen con el total a repagar (capital + fee + interés a 30 días) y confirmá. El adelanto queda registrado y el crédito usado se descuenta del disponible.

### 4. Pedir entrega desde el depósito

Si elegiste logística *depósito*, el stock queda retenido en la sucursal. Desde **Mi Depósito** podés ver las unidades guardadas por adelanto y tocar **Pedir envío** para elegir la sucursal de retiro y despachar la entrega.

### 5. Pagar un adelanto

Desde **Mi Stock** aparecen los adelantos activos con su estado (`al día`, `en riesgo`, `vencido`) y el monto a pagar. Tocá **Pagar adelanto** en cualquiera de ellos para registrar el repago: el crédito se libera y queda disponible para el próximo adelanto.

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
| Registro de negocios + validación RENAPER/BCRA | Implementado |
| Catálogo de productos y mayoristas | Implementado |
| Ciclo completo de adelantos (crear / pagar / pedir entrega) | Implementado |
| Logística dual (`envio` / `deposito`) | Implementado |
| ADN Social — subida de documentos y scoring | Implementado |
| Programa de referidos + impacto en scoring | Implementado |
| ML scoring y forecasting | Stub — spec en `specs/features/04-ml-scoring.md` |
| Grafo de relaciones entre comerciantes | Pendiente |
| Persistencia en base de datos | Pendiente (hoy: in-memory) |
| Enforcement del límite de crédito | Pendiente (known gap en `specs/domain.md`) |
| Aging de status (`ok → warn → danger`) | Pendiente |
