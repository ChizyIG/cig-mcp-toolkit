# Contributing to cig-mcp-toolkit

Thanks for your interest in contributing. This project is maintained by CIG Engineering and is open to contributions from the broader quant/finance and MCP communities.

> **Project status:** `cig-mcp-toolkit` is pre-v0.1. The Python (`python/`) and Go (`go/`) trees are still being scaffolded into a real project structure; the core library (`cig_mcp/`) and reference servers (`servers/`) described in the README are being built out. The workflows below describe the target steady state — some commands become usable as the corresponding modules land.

## Code of Conduct

By participating, you agree to abide by our [Code of Conduct](./CODE_OF_CONDUCT.md). Report unacceptable behavior to `legal@chizyig.com`.

## Ways to contribute

- **Report a bug** — use the [Bug report](./.github/ISSUE_TEMPLATE/bug_report.yml) issue template.
- **Propose a feature or a new reference server** — use the [Feature request](./.github/ISSUE_TEMPLATE/feature_request.yml) issue template. For larger designs (a new server, a schema change), open a discussion issue first so we can align on scope before code gets written.
- **Report a security issue** — see [SECURITY.md](./SECURITY.md). Do *not* file public issues for vulnerabilities.
- **Send a pull request** — see the workflow below.

## Development setup

### Prerequisites

- Python 3.11+ managed via [`uv`](https://docs.astral.sh/uv/)
- Go 1.22+
- Git

### Python (`python/`)

```bash
cd python
uv sync                  # install dependencies into .venv
uv run pytest            # run tests
uv run ruff check .      # lint
uv run ruff format .     # format
uv run mypy .            # type-check
```

### Go (`go/`)

```bash
cd go
go mod download
go test ./...
go vet ./...
golangci-lint run
```

## Making a change

1. **Fork** the repository and create a topic branch from `main`:
   ```bash
   git checkout -b feat/my-change
   ```
2. **Make focused commits.** We follow [Conventional Commits](https://www.conventionalcommits.org/):
   - `feat:` new functionality
   - `fix:` bug fix
   - `docs:` documentation only
   - `refactor:` code change that is neither a fix nor a feature
   - `test:` adding or updating tests
   - `chore:` tooling, CI, build, dependencies
   - Use `!` or a `BREAKING CHANGE:` footer for breaking changes.
3. **Add tests.** Every new feature or bug fix should come with a test that fails before your change and passes after. MCP servers must have contract tests that exercise the tools/resources they expose.
4. **Update docs.** If you change a public API, a schema, or a server's behavior, update the relevant docs under `docs/` in the same PR.
5. **Run the full local check** before pushing:
   ```bash
   # in python/
   uv run ruff check . && uv run mypy . && uv run pytest
   # in go/
   go vet ./... && golangci-lint run && go test ./...
   ```

## Getting help

If you're stuck on setup or have a question that doesn't fit an issue template, email `support.dev@chizyig.com` or open a thread in GitHub Discussions. For bug reports and feature requests, please use the issue templates so details don't get lost.

## Pull request process

- Keep PRs small and focused. One logical change per PR.
- Fill out the PR template.
- Link the issue your PR closes (`Closes #123`) so it's auto-closed on merge.
- Expect a review turnaround of a few business days. Drive-by PRs that reshape broad swaths of the codebase without prior discussion may be closed with a pointer to open a design issue first.
- Squash-and-merge is the default. Your commit subjects should be meaningful on their own — the squashed title will follow Conventional Commits.

## What we won't merge

- Proprietary data, internal model weights, or anything that looks like a production trading signal. This repo is explicitly open-source plumbing; the edge stays out.
- Dependencies on non-public data feeds. Reference servers must work against public or user-supplied data sources.
- Code without tests, unless it's pure docs/config.

## License

This project is licensed under Apache License 2.0 (see [LICENSE](./LICENSE)). By submitting a contribution you agree that your contribution is licensed under the same terms, as described in Section 5 of the license.
