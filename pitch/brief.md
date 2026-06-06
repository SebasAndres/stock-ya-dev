# Reponé — Brief del proyecto (hand-off)

> Nombre **provisorio**. Fintech de microcrédito para el comercio chico no bancarizado en Argentina.
> Este documento resume el estado del proyecto y las decisiones tomadas hasta ahora. Sirve como contexto para retomar el trabajo (p. ej. guardándolo como `CLAUDE.md` en el directorio del proyecto).

---

## 1. Qué es el proyecto (en una línea)

Damos microcrédito al comercio chico de barrio (almacenes, kioscos) que está fuera del sistema bancario, para que reponga stock, usando un motor de scoring que los grandes no tienen. **No competimos en pagos ni como PSP** — solo construimos herramientas de crédito para el no bancarizado.

---

## 2. Decisiones tomadas (lo que ya está definido)

- **Segmento:** comercios chicos / almacenes / kioscos de barrio, cash-heavy, al margen del crédito bancario.
- **Producto principal:** microcrédito **en pesos para reponer stock**.
- **Producto opcional / fase 2:** crédito en **mercadería** + **warehousing tercerizado** (no lo operamos nosotros) para quien no tiene espacio. Quedó como opcional, no como núcleo del v1.
- **Ingreso:** **spread de tasa de interés** (tasa que cobramos − costo de fondeo).
- **Capital:** **mixto / blended.** Organismos internacionales (BID / Banco Mundial), a quienes les vendemos inclusión y formalización; y sector privado, a quien le vendemos rentabilidad (lado fintech).
- **Originación + scoring (el pivote clave):** dejamos atrás el trabajo de campo caro de cooperativas/iglesias *como mecanismo principal de recolección de datos*, y nos paramos sobre tres fuentes:
  1. **Grafo de comunidades (group lending):** prestamos a grupos; el grafo de quién paga y quién conoce a quién **sale del propio crédito**, no de encuestas caras.
  2. **Datos del proveedor:** historial de compra y pago **verificado por el distribuidor** (crédito embebido).
  3. **Aporte voluntario del comercio:** comparte sus datos para subir su línea (modelo Mercado Pago).
- **El diferencial / foso:** el **grafo de comunidades + el modelo de scoring** para el comercio no bancarizado. Es lo que los grandes no tienen y se vuelve más valioso con cada préstamo (flywheel).
- **Dato que NO usamos:** "facturas en negro" auto-declaradas por el comercio. Se descartaron por ser no verificables y por riesgo de AML. Lo que usamos es el registro del **proveedor** (verificable de su lado) + consentimiento explícito del comercio.
- **Posicionamiento competitivo:** no peleamos el riel de pagos (lo usamos para el repago). Atacamos al comercio cash-heavy e informal que MP y los bancos no pueden scorear.
- **Marco regulatorio:** inscribirnos como **Proveedor No Financiero de Crédito (PNFC)** — permite prestar sin ser banco. Somos **sujeto obligado ante la UIF** → programa KYC/AML desde el día uno. La **formalización es la propuesta de valor**, no un costo.
- **Entregable hecho:** deck de pitch de 15 slides, "Reponé", sin el embudo de mercado (`Repone-pitch.pptx`).

---

## 3. El modelo de negocio (resumen tipo canvas)

- **Cliente:** comercio chico de barrio, cash-heavy, fuera del crédito formal. Sub-segmentos: bancarizados (más info) y no bancarizados (scoreados por proveedor + grafo + aporte voluntario).
- **Propuesta de valor (comercio):** crédito rápido sin estar bancarizado, atado al ciclo de reposición, repago sobre el riel (QR/MP) que ya usa.
- **Propuesta de valor (fondeador):** inclusión + formalización + trazabilidad, todo medible.
- **Socios clave:** distribuidoras medianas (datos + canal de originación), banco/PSP regulado (perímetro legal), organismo internacional (capital), comunidades/grupos (grafo), MP/QR como riel de cobro.
- **Actividades clave:** originación (distribuidor + grupos), scoring, otorgamiento y cobranza, compliance KYC/AML, construcción del dataset/grafo.
- **Recursos clave:** capital prestable, motor de scoring + grafo, registro PNFC + programa UIF, plataforma.
- **Costos:** costo de fondeo, mora/incobrables, originación, tech + compliance, fees de 3PL (solo si se activa warehousing).
- **Ingresos:** spread de tasa (principal). Eventual markup de mercadería en fase 2.

### Estructura de capital (blended)
- **Tramo de primera pérdida + asistencia técnica:** organismo internacional. Absorbe los primeros incobrables y subsidia la originación.
- **Tramo senior:** inversor privado, con retorno ajustado por riesgo, desriesgado por el tramo de primera pérdida.

---

## 4. Datos del problema y del mercado (para el pitch, con fuentes)

**Problema central:** en Argentina ya está resuelta la bancarización de pagos, pero no la de crédito. El comercio chico tiene billetera pero no crédito; se autofinancia o depende del fiado del proveedor.

- Empleo informal: **43%** de la población ocupada (~9,2 M), Q4 2025 (INDEC). En empresas de hasta 5 trabajadores trepa al **69,5%**; el comercio está entre los sectores más afectados.
- Crédito al sector privado ≈ **13,4% del PBI** en Argentina, vs Brasil 72%, España 78%, Suecia 128%.
- **69%** de los adultos tiene cuenta (dic 2024, BCRA), pero solo **18,4 M** de personas accede a algún crédito formal (≈ la mitad de los adultos queda afuera).
- Menos del **15%** de las empresas consigue crédito bancario; las dos fuentes principales de financiamiento de la pyme son el capital propio y la **red de proveedores y clientes** (UIA).
- **66%+** de las empresas no pudo acceder a financiamiento por al menos 6 meses (ICAF 2025).
- El BCRA valida el modelo: los PNFC ya rivalizan con los bancos en préstamos personales, y los datos de pagos electrónicos sirven para construir perfil crediticio de quienes no tienen historial (trabajadores informales).

**Tamaño del segmento (órdenes de magnitud, fuentes distintas):**
- Kioscos formales: **~96.000** (cayeron desde ~112.000 en un año; por primera vez debajo de 100k).
- Comercios de cercanía (autoservicios + almacenes + kioscos): **~350.000** según panel privado.
- Registro MiPyME: **~1,8 M** de inscriptos; ~1/4 es comercio (~400k unidades comerciales formales). El **84%** son microempresas (<10 empleados).

---

## 5. Precedentes y competencia

**Precedentes que validan el modelo:**
- **BID Lab → Wayni Móvil:** ya financió micropréstamos a comercios en Argentina; conectado a una red de ~245.000 kioscos.
- **Tarjeta Naranja:** crédito al comercio nacido **fuera del banco** (arrancó como la financiación de la tienda deportiva Salto 96 en Córdoba; Banco Galicia entró con el 49% en 1995 para escalar). Modelo relacional, ellos absorbían el riesgo de mora.
- **Microfinanzas / group lending:** Grameen (Nobel de la Paz 2006) y el enfoque experimental de Banerjee-Duflo-Kremer (Nobel de Economía 2019).

**Competencia:**
- **Mercado Pago / Mercado Crédito:** presta, pero solo a quien ya genera huella digital propia (scoring dominado por préstamos previos y datos transaccionales internos; Dinero Plus exige piso de facturación vía MP/ML). El comercio cash-heavy le es invisible.
- **Embedded finance (Mercado Libre, Galicia, ICBC; plataformas CPG):** crédito embebido en ventas, pero no en el comercio cash-heavy ni con agregación multi-proveedor.

---

## 6. Riesgos y guardarraíles (tener en cuenta)

- **AML:** como sujeto obligado, no se puede ser ingenuo con el origen de fondos. La formalización es el argumento, no la excusa.
- **Trampa de "captación":** un warehouse que paga interés sobre depósitos = intermediación financiera = requiere licencia bancaria. Si se activa el warehousing, estructurarlo como **custodia/warrant**, nunca como depósito que rinde.
- **Spread en contexto argentino:** el fondeo barato (concesional) baja el costo de fondos y ensancha el spread. Riesgos: inflación → usar **plazos cortos con repricing rápido**; techo de tasa impuesto por el mandato del fondeador (no se puede gougear); unit economics ajustadas → se arreglan con el flywheel de datos y el repago recurrente.
- **Crédito embebido (data de proveedor):** dependencia/concentración del distribuidor → varias distribuidoras medianas, no una; visión parcial (un proveedor ve solo lo suyo) → arrancar por el dominante y agregar; riesgo de desintermediación por CPG/ML; **selección hacia el comercio semi-formal** (el 100% informal no está en ningún sistema de proveedor todavía).
- **Group lending / "caso India":** doble filo. India tiene el modelo de grupos más grande (SHG/JLG) pero también la crisis de Andhra Pradesh (2010): sobreendeudamiento y cobranza coercitiva. Tomar el mecanismo con salvaguardas: **préstamo-sonda chico, crédito graduado, cobranza blanda primero (legal al final), nunca coerción**.
- **Grafo:** tiene arranque en frío (al inicio la red es chica; apoyarse más en data de proveedor + crédito graduado). El valor compone con el volumen.
- **Privacidad:** el grafo de relaciones + el aporte voluntario exigen **consentimiento explícito (Ley 25.326)**, visible en el onboarding.
- **Passthrough del fee logístico:** trasladar el costo de warehousing al comercio chico (el más sensible al precio) puede matar la adopción; validar disposición a pagar o que lo subsidie el fondeador. (Aplica solo si se activa la fase 2.)

---

## 7. Qué quedó pausado o descartado

- **Tokens de productos:** quedan como ledger interno de inventario en fase 2, no como feature de cara al cliente ni en el pitch.
- **Warehouse que paga interés (plazo fijo en mercadería):** descartado por la trampa de captación; warehousing solo como logística/custodia opcional.
- **Cooperativas/iglesias como mecanismo principal de recolección de datos:** degradado; la comunidad vuelve vía group lending (el grafo sale del crédito, no de encuestas).
- **"Facturas en negro" como fuente de scoring:** descartado (AML + no verificable).
- **El embudo de mercado (TAM/SAM/SOM):** se pidió dejarlo fuera del deck.

---

## 8. Entregables

- `Repone-pitch.pptx` — deck de 15 slides: portada → problema → evidencia (gráfico crédito/PBI) → insight (el proveedor es el banco) → solución → cómo funciona → motor de scoring (grafo) → diferencial/foso → modelo de negocio → go-to-market → competencia → regulación → impacto/fundeable → riesgos → roadmap & ask.

---

## 9. Próximos pasos / preguntas abiertas

- Definir el **beachhead** concreto: 1-2 distribuidoras o zonas, y cuántos comercios se pueden originar en 12 meses.
- Modelo de **unit economics** del spread (costo de fondeo, mora tolerable, opex por préstamo, cartera de breakeven).
- Slide/diseño técnico del **grafo y el score**: features, señales de red, peso de cada fuente.
- Estructura legal/societaria del **convenio con un PSP/banco** regulado.
- Definir qué tramo de plata del BID se pide: grant/asistencia técnica (plataforma + piloto) vs. deuda/garantía (cartera).
