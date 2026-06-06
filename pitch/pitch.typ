// =============================================================
//  PITCH DECK — fintech de microcrédito al comercio de barrio
//  Render:  typst compile pitch.typ
//  Marca:   typst compile --input brand=StockYa --input tagline=Reponé pitch.typ
//  Watch:   typst watch pitch.typ
// =============================================================

// ---- Marca (override por variable de entorno / --input) ----
#let brand   = sys.inputs.at("brand",   default: "StockYa")
#let tagline = sys.inputs.at("tagline", default: "Reponé")

// ---- Paleta ----
#let navy  = rgb("#0E2F3C")
#let teal  = rgb("#159C8E")
#let amber = rgb("#B96510")
#let slate = rgb("#33454F")
#let gray  = rgb("#6B7B85")
#let mist  = rgb("#F2F5F6")
#let ice   = rgb("#D5DEE3")
#let green = rgb("#2E9E6B")
#let risk  = rgb("#C0552B")

// ---- Tipografías ----
#let serif = "Georgia"
#let sans  = ("Helvetica Neue", "Arial")

// ---- Página ----
#set page(
  width: 10in, height: 5.625in,
  margin: (left: 0.55in, right: 0.55in, top: 0.42in, bottom: 0.4in),
  fill: white,
  footer: context {
    let n = counter(page).get().first()
    if n == 1 { return }
    set text(font: sans, size: 9pt, fill: gray)
    grid(columns: (1fr, auto),
      align(left, emph[#brand · #tagline]),
      align(right, str(n)),
    )
  },
)
#set text(font: sans, size: 11.5pt, fill: slate, lang: "es")
#set par(leading: 0.62em)

// ---- Componentes ----
#let h1(t) = text(font: serif, weight: "bold", size: 23pt, fill: navy, t)
#let lead(t) = text(font: sans, style: "italic", size: 12.5pt, fill: amber, t)

#let slide(title: none, subtitle: none, body) = {
  pagebreak(weak: true)
  if title != none { h1(title) }
  if subtitle != none { v(3pt); lead(subtitle) }
  v(12pt)
  body
}

#let card(title, body, dark: false, h: auto) = {
  let bg   = if dark { navy } else { mist }
  let tcol = if dark { teal } else { gray }
  let bcol = if dark { ice } else { slate }
  block(fill: bg, radius: 8pt, inset: 16pt, width: 100%, height: h)[
    #if title != [] {
      text(font: serif, weight: "bold", size: 15pt, fill: tcol, title)
      v(7pt)
    }
    #set text(size: 11.5pt, fill: bcol)
    #set list(spacing: 0.9em, marker: text(fill: tcol)[•])
    #body
  ]
}

#let badge(n, color: navy) = box(
  circle(radius: 13pt, fill: color,
    align(center + horizon, text(fill: white, weight: "bold", size: 12pt, str(n)))),
)

#let step(n, title, body, color: navy) = grid(
  columns: (auto, 1fr), gutter: 13pt, align: (top, top),
  badge(n, color: color),
  [
    #text(font: serif, weight: "bold", size: 13.5pt, fill: navy, title)
    #v(2pt)
    #text(size: 11pt, fill: slate, body)
  ],
)

#let statbig(num, label, color: navy) = align(center)[
  #text(font: serif, weight: "bold", size: 30pt, fill: color, num) \
  #text(size: 10.5pt, fill: slate, label)
]

// barra horizontal para el gráfico de crédito/PBI
#let barrow(label, value, maxv, hi: false) = grid(
  columns: (1.3in, 1fr), gutter: 10pt, align: (right + horizon, left + horizon),
  text(size: 10.5pt, fill: slate, label),
  box(width: 100%)[
    #stack(dir: ltr, spacing: 7pt,
      rect(width: value / maxv * 100%, height: 17pt, radius: 2pt,
        fill: if hi { amber } else { navy }),
      text(size: 10pt, weight: "bold", fill: if hi { amber } else { gray })[#value%],
    )
  ],
)

// nodo del grafo de comunidades
#let node(c) = circle(radius: 7pt, fill: c, stroke: 1pt + white)

// =============================================================
//  SLIDE 1 — Portada
// =============================================================
#page(fill: navy, margin: 0pt, footer: none)[
  #block(inset: (x: 0.7in, top: 0.85in, bottom: 0.7in), width: 100%, height: 100%)[
    #text(font: serif, weight: "bold", size: 48pt, fill: white, brand)
    #v(6pt)
    #box(width: 1.1in, height: 4pt, fill: teal)
    #v(20pt)
    #text(size: 16pt, fill: ice)[
      Microcrédito para que el comercio de barrio reponga su góndola, \
      con un scoring construido sobre los datos de sus proveedores.
    ]
    #v(1fr)
    #text(size: 12pt, fill: teal, weight: "bold")[
      Inclusión financiera #h(10pt)·#h(10pt) Crédito embebido #h(10pt)·#h(10pt) Argentina
    ]
  ]
]

// =============================================================
//  SLIDE 2 — Problema
// =============================================================
#slide(title: "No tienen crédito.")[
  #grid(columns: (1.55fr, 1fr), gutter: 34pt, align: (top, horizon),
    [
      #text(size: 14pt, fill: slate)[
        No pueden acceder a créditos competitivos (aunque estén bancarizados).
      ]
      #v(14pt)
      #text(size: 13pt, fill: slate)[
        *El comercio chico* (almacenes, kioscos) se autofinancia o vive del
        fiado de su proveedor. El sistema bancario es demasiado chico y demasiado
        exigente para atenderlo.
      ]
    ],
    card([], h: 3.4in)[
      #align(center + horizon)[
        #text(font: serif, weight: "bold", size: 40pt, fill: amber)[≈ 2 de 3]
        #v(6pt)
        #text(size: 11.5pt, fill: slate)[comercios no accedió a financiamiento en los últimos seis meses]
        #v(10pt)
        #text(size: 9.5pt, style: "italic", fill: gray)[ICAF 2025]
      ]
    ],
  )
]

// =============================================================
//  SLIDE 3 — Evidencia: crédito / PBI
// =============================================================
#slide(title: "El sistema financiero argentino casi no presta")[
  #grid(columns: (1.4fr, 1fr), gutter: 36pt, align: (top, horizon),
    [
      #text(size: 11.5pt, fill: gray, weight: "bold")[Crédito al sector privado, como % del PBI]
      #v(12pt)
      #barrow("Argentina", 13.4, 128, hi: true)
      #v(7pt)
      #barrow("Brasil", 72, 128)
      #v(7pt)
      #barrow("España", 78, 128)
      #v(7pt)
      #barrow("Suecia", 128, 128)
    ],
    [
      #statbig("43%", "del empleo es informal (69,5% en los comercios de hasta 5 personas)")
      #v(20pt)
      #statbig("18,4 M", "personas acceden a algún crédito formal, aunque el 69% ya tiene cuenta", color: teal)
    ],
  )
]

// =============================================================
//  SLIDE 5 — Solución
// =============================================================
#slide(
  title: "Crédito para reponer, con un scoring que el banco no tiene",
  subtitle: "Es crédito embebido en la cadena de suministro.",
)[
  #text(size: 12.5pt, fill: slate)[
    Damos microcrédito en pesos al comercio chico para que reponga stock, y lo
    aprobamos rápido usando los datos de compra y pago verificados por sus proveedores.
  ]
  #v(16pt)
  #grid(columns: (1fr, 1fr, 1fr), gutter: 16pt,
    card("Sin estar bancarizado", h: 2.2in)[Aprobación en minutos, sin los requisitos del banco.],
    card("Atado a la reposición", h: 2.2in)[El crédito sigue el ritmo en que el comercio vende y restockea.],
    card("Repago sobre su propio riel", h: 2.2in)[Cobramos por el QR / billetera que el comercio ya usa.],
  )
]

// =============================================================
//  SLIDE 6 — Por qué el crédito de reposición mueve la aguja
// =============================================================
#slide(
  title: "Financiar la reposición mueve la aguja",
  subtitle: "El comercio ya se financia así, el retorno lo justifica y los análogos escalan.",
)[
  #let stat(num, label, src, color: navy) = card([], h: 2.0in)[
    #align(center)[
      #text(font: serif, weight: "bold", size: 27pt, fill: color)[#num]
      #v(6pt)
      #text(size: 10.5pt, fill: slate)[#label]
      #v(5pt)
      #text(size: 8pt, style: "italic", fill: gray)[#src]
    ]
  ]
  #grid(columns: (1fr, 1fr, 1fr), gutter: 16pt,
    stat("88%", "de las empresas en LatAm vende a plazo (crédito de proveedor), pero el comercio chico paga al contado y queda afuera.", "Coface, 2024"),
    stat("~60% anual", "de retorno real al capital en microempresas: supera con holgura la tasa de un crédito de reposición.", "de Mel, McKenzie y Woodruff, QJE 2008", color: teal),
    stat("6–8%", "de quiebres de stock en comercios independientes (el doble que las cadenas); ~40% es venta perdida para siempre.", "Gruen y Corsten; Aastrup-Kotzab", color: amber),
  )
  #v(12pt)
  #card([], dark: true)[
    #text(font: serif, weight: "bold", size: 12.5pt, fill: teal)[Ya funciona a escala: ]
    #text(size: 11.5pt, fill: ice)[Tienda Pago (Perú/México) y Jaza Duka (Unilever + Mastercard, Kenya) dan crédito corto atado a la compra de inventario → *+15–25% en ventas* del comercio, usado 2,8–3,3 veces por mes.]
  ]
]

// =============================================================
//  SLIDE 7 — Motor de scoring
// =============================================================
#slide(
  title: "El motor de scoring: nuestro verdadero diferencial",
  subtitle: "Cada préstamo agranda el grafo y afina el modelo. Es un activo que los grandes no tienen.",
)[
  #grid(columns: (1fr, 1.25fr), gutter: 30pt, align: (horizon, top),
    // --- grafo ---
    card([], dark: true, h: 2.9in)[
      #text(font: serif, weight: "bold", size: 12.5pt, fill: white)[Grafo de comunidades (group lending)]
      #v(6pt)
      #box(width: 100%, height: 1.35in)[
        #place(left + top, dx: 30pt,  dy: 8pt,  node(green))
        #place(left + top, dx: 95pt,  dy: 32pt, node(navy))
        #place(left + top, dx: 60pt,  dy: 70pt, node(green))
        #place(left + top, dx: 150pt, dy: 14pt, node(green))
        #place(left + top, dx: 175pt, dy: 60pt, node(risk))
        #place(left + top, dx: 120pt, dy: 84pt, node(navy))
        #place(left + top, dx: 215pt, dy: 30pt, node(green))
        #place(line(start: (37pt, 15pt),  end: (102pt, 39pt), stroke: 0.8pt + ice))
        #place(line(start: (102pt, 39pt), end: (67pt, 77pt),  stroke: 0.8pt + ice))
        #place(line(start: (102pt, 39pt), end: (157pt, 21pt), stroke: 0.8pt + ice))
        #place(line(start: (157pt, 21pt), end: (182pt, 67pt), stroke: 0.8pt + ice))
        #place(line(start: (127pt, 91pt), end: (182pt, 67pt), stroke: 0.8pt + ice))
        #place(line(start: (157pt, 21pt), end: (222pt, 37pt), stroke: 0.8pt + ice))
      ]
      #set text(size: 9.5pt, fill: ice)
      #grid(columns: (auto, auto, auto), gutter: 12pt, align: horizon,
        [#node(green) #h(3pt) pagan],
        [#node(navy)  #h(3pt) conector],
        [#node(risk)  #h(3pt) en riesgo],
      )
    ],
    // --- tres fuentes ---
    [
      #text(size: 11.5pt, fill: gray, weight: "bold")[Tres fuentes alimentan el score:]
      #v(10pt)
      #stack(spacing: 12pt,
        step(1, "Grafo de comunidades", "Prestamos a grupos. Quién paga y quién conoce a quién arma la red (colateral social."),
        step(2, "Datos del proveedor", "Compra y pago verificados, en tiempo real."),
        step(3, "Aporte voluntario del comercio", "Comparte sus datos para subir su línea, como en Mercado Pago."),
      )
    ],
  )
]

// =============================================================
//  SLIDE 8 — Onboarding
// =============================================================
#slide(title: "Onboarding")[
  #v(4pt)
  #grid(columns: (1fr, 1fr), rows: (1fr, 1fr), gutter: 18pt,
    step(1, "Llamada (agentic) de 15 minutos", "Conocemos el comercio y entendemos las necesidades del negocio en específico."),
    step(2, "Profiling", "Se genera el perfil del comercio junto con sus comercios amigos. Se agrega información existente de proveedores."),
    step(3, "Se aprueba el crédito", "Se ofrece un primer crédito para reponer stock. El monto máximo depende del score."),
    step(4, "Cada repago mejora el perfil", "Las interacciones y repagos afinan el scoring y bajan el costo de originar."),
  )
]

// =============================================================
//  SLIDE 9 — Diferencial competitivo
// =============================================================
#slide(
  title: "Vemos a quien nadie más ve",
  subtitle: "No competimos en pagos ni PSP, solo en crédito para el comercio no bancarizado.",
)[
  #grid(columns: (1fr, 1fr), gutter: 22pt,
    card("El banco y Mercado Pago", h: 3.2in)[
      - Scorean solo con su propia huella digital.
      - El comercio cash-heavy es invisible para ellos.
      - El crédito queda gateado por la facturación digital.
    ],
    card("Nosotros", dark: true, h: 3.2in)[
      - Tenemos el grafo de comunidades: quién paga y quién conoce a quién.
      - Lo sumamos a los datos del proveedor y a lo que el comercio aporta.
      - Un foso de datos que los grandes no tienen.
    ],
  )
]

// =============================================================
//  SLIDE 10 — Más datos + crédito que no se desvía (mora)
// =============================================================
#slide(
  title: "Más datos, y crédito que no se desvía",
  subtitle: "Dos palancas que bajan la mora: información que el resto no ve, y crédito atado a la mercadería.",
)[
  #grid(columns: (1fr, 1fr), gutter: 22pt,
    card("Ventaja de información", h: 2.95in)[
      - Tres fuentes que el banco y MP no tienen: compra y pago verificados por el proveedor, el grafo de quién paga a quién, y el aporte voluntario del comercio.
      - Scoreamos al comercio cash-heavy, invisible para la huella digital.
      - Cada préstamo agranda el grafo: el modelo mejora con cada operación (flywheel).
    ],
    card("Crédito productivo, menos mora", dark: true, h: 2.95in)[
      - Prestamos para reponer stock (en fase 2, en mercadería). El crédito entra al negocio, no a la caja.
      - El efectivo se desvía a consumo; el capital en especie se queda y genera la venta que repaga.
      - El destino es verificable: financiamos la góndola, no un gasto que no podemos ver.
    ],
  )
  #v(8pt)
  #text(size: 9pt, style: "italic", fill: gray)[
    Evidencia: Fafchamps, McKenzie, Quinn y Woodruff (2014), “Microenterprise growth and the flypaper effect”, _Journal of Development Economics_ — solo el capital en especie hizo crecer la ganancia; el efectivo, no.
  ]
]

// =============================================================
//  SLIDE 11 — Evidencia de mora
// =============================================================
#slide(
  title: "La mora es del instrumento, no del segmento",
  subtitle: "El segmento no es incobrable: el crédito que hoy le llega está mal diseñado.",
)[
  #grid(columns: (1fr, 1fr), gutter: 22pt,
    card("El no bancarizado hoy paga mal")[
      - En Argentina, los deudores *sin banco* tienen *36,1%* de irregularidad; las entidades no financieras, *21–44%* de mora, frente a *~1,6%* de los bancos.
      - No es la persona, es el instrumento: hoy reciben *efectivo caro y sin control de destino*.
    ],
    card("Cambiá el instrumento, baja la mora", dark: true)[
      - *En especie / destino verificado:* el efectivo se desvía a consumo; la mercadería genera la venta que repaga (Fafchamps, 2014). Crédito etiquetado ≈ *100% de repago* (Augsburg, 2023).
      - *Grupo + datos propios:* Grameen sostuvo *98–99%*; Mercado Crédito, scoreando con datos transaccionales y sin bureau, tiene *6,7%* de mora a 15–90 días.
    ],
  )
  #v(8pt)
  #text(size: 9pt, style: "italic", fill: gray)[
    Fuentes: BCRA / Ámbito (feb 2026); Fafchamps, McKenzie, Quinn y Woodruff, _JDE_ 2014 (RCT Ghana, n=793); Augsburg et al., _JDE_ 2023 (RCT India); Giné y Karlan, _JDE_ 2014; Grameen Bank; MercadoLibre Q2 2025.
  ]
]

// =============================================================
//  SLIDE 12 — Modelo de negocio
// =============================================================
#slide(title: "Cómo ganamos plata")[
  #grid(columns: (1fr, 1fr), gutter: 24pt, align: (horizon, top),
    [
      #card([], dark: true)[
        #text(font: serif, weight: "bold", size: 14pt, fill: teal)[Ingreso: spread de tasa]
        #v(8pt)
        #align(center, text(size: 13pt, fill: ice)[
          tasa que cobramos #h(6pt) − #h(6pt) costo de fondeo
        ])
      ]
      #v(10pt)
      #text(size: 11pt, fill: slate)[
        El capital concesional del organismo internacional baja nuestro costo de
        fondeo: ensancha el margen o nos deja bajar la tasa al comercio.
      ]
    ],
    [
      #text(size: 11.5pt, fill: gray, weight: "bold")[Capital en capas (blended)]
      #v(10pt)
      #stack(spacing: 10pt,
        card("Tramo senior (inversor privado)")[retorno ajustado por riesgo],
        card("Primera pérdida + asistencia técnica")[organismo internacional (BID / Banco Mundial)],
      )
      #v(8pt)
      #text(size: 10.5pt, style: "italic", fill: teal)[
        El multilateral absorbe la primera pérdida → el privado entra de forma segura.
      ]
    ],
  )
]

// =============================================================
//  SLIDE 13 — Go to market
// =============================================================
#slide(title: "Por dónde entramos")[
  #v(4pt)
  #stack(spacing: 16pt,
    step(1, "Alianza con distribuidoras medianas",
      "Las que no van a montar su propia fintech. Una sola integración nos da la cartera de cientos de comercios."),
    step(2, "Crédito graduado",
      "Empezamos con montos chicos y crecemos con la conducta de pago. Bajo riesgo desde el día uno."),
    step(3, "Prestamos a comunidades (group lending)",
      "El préstamo en grupo construye el grafo: quién paga y quién conoce a quién. La data sale del crédito, no de encuestas caras."),
  )
]

// =============================================================
//  SLIDE 14 — Competencia
// =============================================================
#slide(
  title: "El panorama competitivo",
  subtitle: "Nuestra cuña: el comercio chico, informal y de efectivo que los grandes subatienden, y la agregación de varios proveedores que un solo CPG no hace.",
)[
  #grid(columns: (1fr, 1fr, 1fr), gutter: 16pt,
    card("Mercado Pago / Mercado Crédito", h: 2.2in)[Presta, pero solo a quien ya genera huella digital propia.],
    card("Tarjeta Naranja", h: 2.2in)[Probó el modelo: crédito al comercio desde afuera del banco.],
    card("Embedded players (ML, Galicia, ICBC)", h: 2.2in)[Crédito embebido en ventas, pero no en el comercio cash-heavy.],
  )
]

// =============================================================
//  SLIDE 16 — Fundeable
// =============================================================
#slide(title: "Por qué un organismo internacional financia esto")[
  #let metric(t) = card([], dark: true, h: 0.75in)[
    #align(center + horizon, text(font: serif, weight: "bold", size: 15pt, fill: teal, t))
  ]
  #text(size: 11.5pt, fill: gray, weight: "bold")[Las métricas que un BID o Banco Mundial compra:]
  #v(10pt)
  #grid(columns: (1fr, 1fr, 1fr), gutter: 16pt,
    metric("Inclusión"), metric("Formalización"), metric("Trazabilidad"),
  )
  #v(16pt)
  #text(size: 11.5pt, fill: gray, weight: "bold")[Y el modelo ya tiene antecedentes:]
  #v(10pt)
  #grid(columns: (1fr, 1fr, 1fr), gutter: 16pt,
    card("BID Lab → Wayni Móvil", h: 1.4in)[ya financió micropréstamos a comercios en Argentina],
    card("Tarjeta Naranja", h: 1.4in)[crédito al comercio nacido fuera del banco],
    card("Microfinanzas", h: 1.4in)[Nobel de la Paz 2006 y de Economía 2019 detrás],
  )
]

// =============================================================
//  SLIDE 17 — Riesgos
// =============================================================
#slide(title: "Lo que puede salir mal (y cómo lo manejamos)")[
  #let riskrow(problem, fix) = grid(
    columns: (1fr, auto, 1.1fr), gutter: 14pt, align: (right + horizon, horizon, left + horizon),
    text(size: 12pt, weight: "bold", fill: risk, problem),
    text(size: 13pt, fill: gray)[→],
    text(size: 12pt, fill: slate, fix),
  )
  #v(4pt)
  #stack(spacing: 13pt,
    riskrow("Dependencia del distribuidor", "Varias distribuidoras medianas, nunca una sola."),
    riskrow("Visión parcial del comercio", "Arrancamos por el proveedor dominante y agregamos."),
    riskrow("Sector en crisis, mora alta", "Préstamo-sonda chico + crédito graduado."),
    riskrow("Sobreendeudamiento: Andhra Pradesh 2010 (repago 95% → 1%)", "Límite de exposición por grafo, crédito graduado y cobranza blanda; sin metas de colocación agresivas."),
    riskrow("Macro e inflación", "Plazos cortos con repricing rápido."),
  )
]

// =============================================================
//  SLIDE 18 — Roadmap & ask
// =============================================================
#slide(title: "Hacia dónde vamos")[
  #grid(columns: (1fr, 1fr, 1fr), gutter: 16pt,
    card("Fase 1", h: 1.95in)[Crédito en pesos para reponer, scoring con data de proveedor, piloto con 1-2 distribuidoras.],
    card("Fase 2", h: 1.95in)[Crédito en mercadería y warehouse opcional para quien no tiene espacio.],
    card("Fase 3", h: 1.95in)[Agregación multi-proveedor y scoring propio que reemplaza al juicio humano.],
  )
  #v(14pt)
  #block(fill: navy, radius: 8pt, inset: 14pt, width: 100%)[
    #text(font: serif, weight: "bold", size: 13pt, fill: teal)[El ask: ]
    #text(size: 12pt, fill: ice)[
      asistencia técnica + primera pérdida del organismo internacional, tramo senior
      privado, y alianzas con distribuidoras para el piloto.
    ]
  ]
]
