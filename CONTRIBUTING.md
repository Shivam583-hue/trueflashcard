# Contributing

Thanks for your interest in improving True Flashcards! This guide covers how to get set up,
the conventions the codebase follows, and what a good pull request looks like.

## Getting set up

Follow the **Getting started** section in the [README](README.md) — clone, copy the
`.env.example` files, `make install`, `make migrate`, `make run`. You'll want a Postgres
database (a free Neon project is plenty).

## Ways to contribute

- **Bugs** — open an issue with steps to reproduce, what you expected, and what happened.
- **Features / changes** — for anything non-trivial, open an issue first so we can agree on
  the approach before you invest time.
- **Docs** — fixes and clarifications are always welcome.

## Development workflow

1. Fork the repo and create a branch off `main` (e.g. `feat/deck-shuffle` or `fix/import-trim`).
2. Make your change, keeping it focused — one logical change per PR.
3. Run the checks below until they're green.
4. Commit with a short, imperative message (`Add deck shuffle mode`, not `added stuff`).
5. Open a pull request describing **what** changed and **why**.

### Checks before you push

**Backend** (`server/`):

```bash
go build ./...
go vet ./...
go test ./...
```

**Frontend** (`web/`):

```bash
pnpm lint
pnpm tsc --noEmit
pnpm build
```

## Code conventions

These keep the codebase consistent — please follow them:

- **Self-documenting code, no comments.** Express intent through naming and structure rather
  than inline comments. The rare exception is a short note explaining a genuinely non-obvious
  *why*.
- **Production-style structure.** Clear package/module boundaries, no god files. New backend
  logic belongs in the right `internal/` package; new UI in a focused component.
- **Frontend aesthetic.** Minimalist, dark, gray gradients. Deliberate spacing and restrained,
  purposeful motion — animate `transform`/`opacity` only, use the shared easing tokens, and
  honor `prefers-reduced-motion`. Avoid generic "AI slop" UI.
- **Security at the boundary.** All SQL goes through sqlc (parameterized — never interpolate
  user input). Validate and sanitize every input at the gRPC/Connect boundary.

## Working with generated code

Some code is generated — edit the source, not the output, then regenerate:

- **Protobuf / RPC contract** — edit `server/proto/**/*.proto`, then `make proto` (regenerates
  Go, Connect, and the TypeScript client).
- **Database queries** — edit `server/db/queries/*.sql`, then run `sqlc generate` in `server/`.
- **Migrations** — add a new pair of `server/db/migrations/NNNN_name.{up,down}.sql` files;
  never edit a migration that's already been applied.

Commit the regenerated files alongside your change.

## Tests

- Add or update backend tests for new service logic (`server/internal/**/*_test.go`). The
  existing tests use a stubbed querier so they run without a database — follow that pattern.
- For the frontend, make sure `lint`, `tsc`, and `build` all pass.

## Pull request checklist

- [ ] Focused, single-purpose change
- [ ] Backend: `go build`, `go vet`, `go test` pass
- [ ] Frontend: `pnpm lint`, `pnpm tsc --noEmit`, `pnpm build` pass
- [ ] Generated code regenerated and committed (if the proto/queries changed)
- [ ] Clear PR description (what + why)

By contributing, you agree that your contributions are licensed under the project's
[MIT License](LICENSE).
