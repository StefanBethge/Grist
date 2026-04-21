# Grist — TUI + Recipe-Runtime für gseq-table

## Was ist Grist?

Grist ist ein interaktives Terminal-Tool zum Bereinigen und Transformieren tabellarischer Daten (CSV, Excel, JSON). Es setzt auf der Go-Library `github.com/stefanbethge/gseq-table` auf und schließt deren fachliche Lücke: **Exploration interaktiv, Produktion deklarativ.**

Typischer Flow:
1. User öffnet messy Datei in der TUI
2. Klickt/navigiert sich durch Bereinigungs-Operationen (trim, cast, rename, filter, validate, …)
3. Exportiert das Ergebnis als **Recipe** (YAML/JSON) + optional generierten Go-Code
4. In Produktion läuft entweder `grist run recipe.yaml` oder der kompilierte Go-Binary

Vorbilder: OpenRefine (JSON-History), Power Query → M, Trifacta Recipes, Postman „Export as Code".

## Fachliche Identität (geerbt von gseq-table)

- **Zielgruppe:** Entwickler, die messy Daten für Migrationen/Importe/ETL aufbereiten
- **Nicht-Zielgruppe:** Data-Science / Analytics (→ pandas/Polars/DuckDB)
- **Daten-Volumen:** In-Memory, typisch 10k–5M Zeilen
- **Philosophie:** String-first, Fehlerakkumulation statt Fail-Fast, Immutable-first

## Architektur

```
[TUI] ─ operationen ─▶ [Recipe YAML/JSON] ─┬─▶ [Go-Code-Generator]
                             ▲             ├─▶ [Runtime-Replay: grist run]
                             └─ zurück lesbar └─▶ [Audit / Doku]
```

**Recipe als Kanon** — nicht direkt Go-Code ausgeben. Das Recipe-Format ist das eigentliche Produkt. TUI ist Editor, Go-Codegen ist Export.

### Beispiel-Recipe

```yaml
source:
  type: csv
  path: in.csv
  sep: ";"
  encoding: utf-8

steps:
  - trim: { columns: [name, email] }
  - rename: { from: "E-Mail", to: email }
  - cast:
      column: geburtsdatum
      type: date
      format: "02.01.2006"
  - drop_if:
      column: email
      matches: "^$"
  - validate:
      schema: schema.yaml

sink:
  type: csv
  path: out.csv
```

### Begründung Recipe-als-Kanon

- Diff-bar im Git (Reviewer sehen Regeln, nicht generierten Code-Noise)
- Reproduzierbar ohne Go-Toolchain (`grist run`)
- Roundtrip: Recipe → TUI → Recipe (iterieren möglich)
- Go-Codegen kann sich weiterentwickeln, ohne Altlasten zu haben
- Zwei Consumer: Fixed-Binary-Fans (Go-Codegen) vs. Dynamic-Replay (`grist run`)

## Implementierungsreihenfolge

**Wichtig:** Nicht mit der TUI anfangen. Diese Reihenfolge sorgt dafür, dass jeder Schritt für sich Wert liefert — falls ein späterer Schritt nie kommt, steht das Projekt trotzdem.

### Phase 1: Recipe-Format + Runtime (1–2 Wochen)
- YAML/JSON-Schema für Recipes definieren
- Jede gseq-table-Operation bekommt eine Recipe-Repräsentation
- `grist run recipe.yaml` als CLI-Subcommand
- **Deliverable:** Deklarative Pipelines per YAML, auch ohne TUI nutzbar

### Phase 2: Go-Code-Generator (3–5 Tage)
- Nimmt Recipe, emittiert lesbaren Go-Code gegen die gseq-table-API
- `grist build recipe.yaml` → `main.go` oder Package
- **Deliverable:** „Recipe → Go-Projekt scaffolden", unabhängig von TUI

### Phase 3: TUI (4–6 Wochen Fokusarbeit)
- Interaktiver Editor, der Recipes erzeugt
- Stack: `charmbracelet/bubbletea` + `bubbles` + `lipgloss`
- **Deliverable:** Fachliche User können ohne Go-Kenntnisse arbeiten

## MVP-Scope der TUI

### Rein im MVP
- Load: CSV / Excel / JSON
- Table-Browser: paginiert, Spalten-Details (Typ-Inferenz, Null-Count, Unique-Count, Samples)
- Ops: rename, drop, trim, cast, filter, split/merge columns, regex-replace
- Schema-Preview: Validierungsergebnisse live
- Error-Log-Ansicht: gefiltert, sortiert, gruppiert
- Export: Recipe (YAML) + generierter Go-Code + transformiertes Output

### Explizit NICHT im MVP
- Multi-Table / Joins (zu komplexe UX für MVP)
- SQL-Sources (erst wenn `gseq-table/sql` existiert)
- Custom-Go-Code-Snippets in der TUI (gegen Editoren kann man nicht konkurrieren)
- Plugin-System
- Visual-Pipeline-Graph
- Remote-Collaboration

## Tech-Stack

- **Sprache:** Go 1.23+
- **Core-Dependency:** `github.com/stefanbethge/gseq-table`
- **TUI:** `github.com/charmbracelet/bubbletea` + `bubbles` + `lipgloss`
- **Recipe-Serialisierung:** `gopkg.in/yaml.v3` (primär YAML, JSON als Alternative)
- **Codegen:** Go-Templates (`text/template`) + `go/format` für Formatierung
- **Testing:** Standard-Library + `testscript` für CLI-Tests

## Projektstruktur (Vorschlag)

```
grist/
├── cmd/grist/          # CLI-Entrypoint
├── recipe/             # Recipe-Format + Parser + Validator
├── runtime/            # Recipe-Executor (grist run)
├── codegen/            # Go-Code-Generator (grist build)
├── tui/                # Bubbletea-App
├── ops/                # Op-Definitionen (gemeinsam genutzt)
├── examples/           # Beispiel-Recipes + generierter Output
├── go.mod
├── go.sum
├── README.md
└── CLAUDE.md
```

## Richtlinien für die Entwicklung

### Code-Qualität (geerbt von gseq-table)
- gofmt-konform, `go vet` clean, `go test -race ./...` clean
- Type Assertions immer mit `ok`-Check
- Resource-Cleanup via `defer`
- Exportierte Typen/Funktionen dokumentiert

### API-Stil
- Result-basierte Fehlerbehandlung (kompatibel mit gseq-Muster)
- Functional Options für erweiterbare Konfiguration
- Keine Panics in normalem Kontrollfluss

### Recipe-Stabilität
- Recipes sind Nutzer-Daten → Versionierung wichtig (`version: 1` im Header)
- Breaking Changes nur mit Migrations-Pfad
- Unbekannte Felder → Warnung, nicht Fehler (Forward-Compat)

### TUI-Prinzipien
- Keyboard-first, Maus optional
- Jede Operation sofort sichtbar (Preview)
- Undo/Redo auf Recipe-Ebene (letzter Step entfernen ≠ Daten-Undo)
- SSH- und tmux-kompatibel

## Risiken / bewusste Trade-offs

1. **UX-Aufwand wird unterschätzt.** 80% der TUI-Arbeit ist „fühlt es sich richtig an?" — rechne mit 2–3× der initialen Schätzung.
2. **Solo-Maintainer-Risiko.** Library + TUI = doppelte Last. Deshalb Recipe-first-Strategie: kein verlorener Schritt.
3. **Scope-Creep-Magnet.** Disziplin beim MVP ist Pflicht. Im Zweifel: nein.
4. **Terminal-Kompat.** Windows-Terminals, SSH, tmux, Fonts — jede Kombination ein potentieller Bug.

## Non-Goals (explizit nicht gewollt)

- Out-of-Memory / Streaming-Engine (→ DuckDB existiert, ist unschlagbar)
- Eigenes Query-Language
- Distributed Processing
- Analytics / BI / Reporting
- Realtime / Event-Processing
- Visuelles Drag-&-Drop im Browser (wir bauen Terminal, nicht Web-UI)

## Verwandte Repos

- `github.com/stefanbethge/gseq-table` — Core-Library (Basis, separates Projekt, separate Release-Kadenz)
- `github.com/stefanbethge/gseq` — Utility-Library (Result, Option, etc.)

## Namens-Herkunft

„Grist" wie in *grist for the mill* — Rohmaterial, das durchs Mahlen Wert bekommt. Passt zur fachlichen Identität: messy Input → verwertbares Output.

## Offene Fragen (für zukünftige Entscheidungen)

- [ ] Recipe-Format: YAML primär oder JSON primär?
- [ ] Codegen: eine Datei oder Package-Struktur?
- [ ] Recipe-Versionierung: Semver im Header oder nur `version: N`?
- [ ] Lizenz: MIT wie gseq-table, oder anders?
- [ ] Modul-Pfad: `github.com/stefanbethge/grist`?
- [ ] Beziehung zu gseq-table-Releases: loose coupling per Semver oder Monorepo-Abhängigkeit?