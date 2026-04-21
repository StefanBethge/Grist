# Grist

Grist is a terminal tool for cleaning and transforming tabular data (CSV,
Excel, JSON). It is built on top of
[`gseq-table`](https://github.com/stefanbethge/gseq-table) and closes its
usability gap: **explore interactively, run declaratively.**

The canonical artefact is a **recipe** (YAML). Recipes are produced
interactively in the TUI, executed by `grist run`, or compiled to Go source
by `grist build`.

> Status: early scaffolding. See `CLAUDE.md` for the project plan and phase
> breakdown.

## Install

```sh
go install github.com/stefanbethge/grist/cmd/grist@latest
```

## Usage

```sh
grist run   recipe.yaml   # apply a recipe to its declared source
grist build recipe.yaml   # generate Go source from a recipe
grist tui                 # interactive editor (Phase 3, not implemented yet)
```

## Development

```sh
make check   # gofmt, go vet, go test -race
make build   # build ./grist
```

Requires Go 1.23 or newer.

## License

MIT — see [LICENSE](./LICENSE).
