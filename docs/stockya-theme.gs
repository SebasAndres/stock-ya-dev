/**
 * StockYa Theme — aplica la paleta de la UI a la presentación
 *
 * Cómo usar:
 * 1. Abrí la presentación en Google Slides
 * 2. Extensions → Apps Script
 * 3. Pegá este código (reemplazá todo) y guardá
 * 4. Ejecutá applyStockYaTheme()
 * 5. Autorizá los permisos cuando te pida
 */

// ── Paleta StockYa (1:1 con los CSS tokens de frontend/index.html) ──────────
const C = {
  green900: '#0a2e1a',
  green800: '#0d3d22',
  green700: '#115c2e',
  green600: '#157a3a',
  green500: '#1a9645',
  green400: '#22b854',
  green300: '#4dd67a',
  green100: '#d6f5e3',
  green50:  '#edfaf3',

  amber500: '#f59e0b',
  amber400: '#fbbf24',
  amber100: '#fef3c7',

  red500:   '#ef4444',
  red100:   '#fee2e2',

  neutral900: '#111418',
  neutral800: '#1c2128',
  neutral700: '#2d3748',
  neutral600: '#4a5568',
  neutral400: '#9aa5b1',
  neutral200: '#e8ecf0',
  neutral100: '#f2f5f7',
  neutral50:  '#f8fafb',
  white:      '#ffffff',
};

// ── Helpers ──────────────────────────────────────────────────────────────────

function hexToRgb(hex) {
  const r = parseInt(hex.slice(1, 3), 16) / 255;
  const g = parseInt(hex.slice(3, 5), 16) / 255;
  const b = parseInt(hex.slice(5, 7), 16) / 255;
  return { red: r, green: g, blue: b };
}

function setShapeFill(shape, hex) {
  try { shape.getFill().setSolidFill(hex); } catch (_) {}
}

function setTextColor(textRange, hex) {
  try { textRange.getTextStyle().setForegroundColor(hex); } catch (_) {}
}

function setFont(textRange, family, size, bold) {
  try {
    const s = textRange.getTextStyle();
    if (family) s.setFontFamily(family);
    if (size)   s.setFontSize(size);
    if (bold !== undefined) s.setBold(bold);
  } catch (_) {}
}

// ── Clasificación de slides ──────────────────────────────────────────────────
// Slide 1: cover → fondo green900
// Slide 2–14: contenido → fondo green800, excepto slides con fondo alternativo
// Si el título tiene un número solo → slide de sección → fondo green700

function classifySlide(slide, idx) {
  if (idx === 0) return 'cover';
  return 'content';
}

// ── Aplicar fondo ────────────────────────────────────────────────────────────

function applyBackground(slide, type) {
  const bg = slide.getBackground();
  if (type === 'cover') {
    bg.setSolidFill(C.green900);
  } else {
    bg.setSolidFill(C.green800);
  }
}

// ── Procesar cada shape de la slide ─────────────────────────────────────────

function processShape(shape, slideType) {
  const shapeType = shape.getShapeType
    ? shape.getShapeType().toString()
    : '';

  // Shapes de acento / decorativas (rectángulos sin texto o con texto corto de número)
  const fill = shape.getFill ? shape.getFill() : null;
  const hasText = shape.getText && shape.getText().asString().trim().length > 0;
  const text = hasText ? shape.getText().asString().trim() : '';

  // Detectar si es un shape numérico grande (número de slide: "2", "3", etc.)
  const isSlideNumber = /^\d{1,2}$/.test(text) && text.length <= 2;

  // Detectar porcentaje stat prominente (e.g. "88%", "≈60% anual", "2 de 3")
  const isStat = /^[\d≈][^a-z]{0,20}$/.test(text) && text.length < 20;

  if (fill) {
    try {
      const fillType = fill.getType().toString();
      // Solo tocamos shapes con relleno sólido (no imágenes, no gradientes)
      if (fillType === 'SOLID') {
        const currentHex = fill.getSolidFill().getColor().asHexString().toLowerCase();
        // Si el fondo es blanco/claro → lo convertimos a un card oscuro
        if (['#ffffff', '#f8f8f8', '#eeeeee', '#e0e0e0', '#dddddd'].includes(currentHex)) {
          setShapeFill(shape, C.neutral800);
        }
        // Si ya es un verde o color de marca → lo mapeamos al equivalente de la paleta
        if (currentHex.startsWith('#1') || currentHex.startsWith('#2') || currentHex.startsWith('#0')) {
          // Dejar los colores oscuros existentes; solo normalizamos al verde 500
        }
        // Shapes de acento que deberían ser green400
        if (['#4caf50','#43a047','#2e7d32','#1b5e20'].includes(currentHex)) {
          setShapeFill(shape, C.green400);
        }
      }
    } catch (_) {}
  }

  // Procesar texto
  if (!hasText) return;

  try {
    const paragraphs = shape.getText().getParagraphs();
    paragraphs.forEach((para, pIdx) => {
      const paraText = para.getRange().asString().trim();

      // Número de slide al pie → gris suave
      if (isSlideNumber && paragraphs.length === 1) {
        setTextColor(para.getRange(), C.neutral400);
        setFont(para.getRange(), 'DM Sans', null, false);
        return;
      }

      // Stats grandes → amber para contraste
      if (isStat && pIdx === 0 && !isSlideNumber) {
        setTextColor(para.getRange(), C.amber400);
        setFont(para.getRange(), 'Syne', null, true);
        return;
      }

      // Detectar si es título principal (primer párrafo de shapes grandes, en bold)
      const style = para.getRange().getTextStyle();
      let fontSize = 12;
      try { fontSize = style.getFontSize() || 12; } catch (_) {}
      const isBold = (() => { try { return style.isBold(); } catch (_) { return false; } })();

      if (fontSize >= 24 || (isBold && fontSize >= 18)) {
        // Título principal → blanco puro con Syne
        setTextColor(para.getRange(), C.white);
        setFont(para.getRange(), 'Syne', null, true);
      } else if (fontSize >= 14) {
        // Subtítulo → green300
        setTextColor(para.getRange(), C.green300);
        setFont(para.getRange(), 'DM Sans', null, null);
      } else {
        // Body / descriptivo → blanco con opacidad simulada (no hay alpha en Apps Script)
        setTextColor(para.getRange(), C.neutral200);
        setFont(para.getRange(), 'DM Sans', null, null);
      }

      // Override: textos tipo "StockYa · Reponé" al pie → neutral400
      if (paraText.includes('StockYa ·') || paraText.includes('StockYa·')) {
        setTextColor(para.getRange(), C.neutral400);
        setFont(para.getRange(), 'DM Sans', null, false);
      }

      // Override: fuentes de datos (ICAF, Coface, etc.) → neutral400 italic
      if (paraText.match(/ICAF|Coface|Gruen|Fafchamps|Augsburg|QJE|2024|2025|2023|2022/)) {
        setTextColor(para.getRange(), C.neutral400);
      }

      // Override: palabras clave en verde brillante dentro del texto
      // (No podemos hacer word-level sin iterar runs, hacemos el párrafo entero si es corto)
    });
  } catch (_) {}
}

// ── Líneas y conectores → green600 ──────────────────────────────────────────

function processLine(line) {
  try {
    line.getLine().getLineFill().setSolidFill(C.green600);
  } catch (_) {}
}

// ── Cover slide: tratamiento especial ───────────────────────────────────────

function styleCoverSlide(slide) {
  const shapes = slide.getShapes();
  shapes.forEach(shape => {
    if (!shape.getText) return;
    const text = shape.getText().asString().trim();
    if (!text) return;

    const style = shape.getText().getTextStyle();
    let fontSize = 12;
    try { fontSize = style.getFontSize() || 12; } catch (_) {}

    if (text.toLowerCase().includes('stockya')) {
      // Logo text → white + Syne bold
      setTextColor(shape.getText(), C.white);
      setFont(shape.getText(), 'Syne', null, true);
    } else if (fontSize >= 20) {
      setTextColor(shape.getText(), C.white);
      setFont(shape.getText(), 'Syne', null, true);
    } else {
      // Tagline / subtítulo de portada → green300
      setTextColor(shape.getText(), C.green300);
      setFont(shape.getText(), 'DM Sans', null, false);
    }
  });
}

// ── Función principal ────────────────────────────────────────────────────────

function applyStockYaTheme() {
  const pres = SlidesApp.getActivePresentation();
  const slides = pres.getSlides();

  slides.forEach((slide, idx) => {
    const type = classifySlide(slide, idx);
    applyBackground(slide, type);

    if (type === 'cover') {
      styleCoverSlide(slide);
    }

    // Shapes genéricas
    slide.getShapes().forEach(shape => {
      if (type !== 'cover') processShape(shape, type);
    });

    // Tablas
    slide.getTables().forEach(table => {
      for (let r = 0; r < table.getNumRows(); r++) {
        for (let col = 0; col < table.getNumColumns(); col++) {
          const cell = table.getCell(r, col);
          try {
            if (r === 0) {
              // Fila de encabezado → fondo green700, texto white
              cell.getFill().setSolidFill(C.green700);
              setTextColor(cell.getText(), C.white);
              setFont(cell.getText(), 'Syne', null, true);
            } else {
              // Filas alternadas
              const bg = r % 2 === 0 ? C.neutral800 : C.neutral900;
              cell.getFill().setSolidFill(bg);
              setTextColor(cell.getText(), C.neutral200);
              setFont(cell.getText(), 'DM Sans', null, false);
            }
          } catch (_) {}
        }
      }
    });
  });

  // Forzar guardado
  pres.saveAndClose();

  // Necesitamos reabrir la presentación luego de saveAndClose si se quiere seguir editando
  Logger.log('✅ Tema StockYa aplicado correctamente a ' + slides.length + ' slides.');
  SpreadsheetApp.getUi && Browser.msgBox('✅ Tema StockYa aplicado a ' + slides.length + ' slides.');
}
