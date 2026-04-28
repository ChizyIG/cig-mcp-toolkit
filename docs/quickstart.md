# Quickstart

This guide walks you through installing `cig-mcp-toolkit` from source and running the reference servers against an MCP-aware client (e.g., Claude Desktop).

> **Status:** v0.1 ships two reference servers:
>
> - **Yahoo Finance** (Python) — `get_quote` for real-time-ish equity/crypto/FX prices.
> - **SEC EDGAR** (Go) — `lookup_company` and `list_filings` for SEC-registered issuers.

## Prerequisites

- Python 3.11+ and [`uv`](https://docs.astral.sh/uv/) — for the Yahoo server.
- Go 1.25+ — for the EDGAR server.
- An MCP client — [Claude Desktop](https://claude.ai/download) is the easiest.

## Clone the repo

```bash
git clone https://github.com/ChizyIG/cig-mcp-toolkit.git
cd cig-mcp-toolkit
```

## Yahoo Finance server (Python)

### Install

```bash
cd python
uv sync
```

This creates `.venv/` with `cig-mcp` and its dependencies (including `yfinance`).

### Run standalone

```bash
uv run cig-mcp-yahoo
# or
uv run python -m cig_mcp.servers.yahoo
```

The process waits for MCP frames on stdin and writes responses to stdout. You will not see output until a client connects.

### Wire into Claude Desktop

Edit `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) and add `mcpServers` as a top-level sibling of `preferences`:

```json
{
  "mcpServers": {
    "cig-yahoo": {
      "command": "uv",
      "args": [
        "--directory",
        "/absolute/path/to/cig-mcp-toolkit/python",
        "run",
        "cig-mcp-yahoo"
      ]
    }
  }
}
```

Fully quit Claude Desktop (Cmd+Q) and relaunch. Ask: *"What's AAPL trading at?"* Claude calls `get_quote("AAPL")` and returns `last_price`, `previous_close`, `change`, `change_percent`, `currency`, `exchange`, `market_cap`.

### Limits

- Yahoo Finance is a public, unauthenticated endpoint — rate limits, occasional 5xx errors, and shape drift are upstream concerns we do not control.
- `fast_info` returns a price snapshot, not real-time tick data.
- For research and prototyping. Do not use it as the data source for production trading.

## SEC EDGAR server (Go)

### Install

```bash
cd go
go install ./cmd/cig-mcp-edgar
```

This puts the `cig-mcp-edgar` binary in `$(go env GOBIN)` (or `$(go env GOPATH)/bin` if `GOBIN` isn't set). Make sure that directory is on your `PATH` — if `which cig-mcp-edgar` comes up empty, add it:

```bash
# zsh
echo 'export PATH="$(go env GOPATH)/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

# bash
echo 'export PATH="$(go env GOPATH)/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

### Run standalone

```bash
cig-mcp-edgar
```

Like the Yahoo server, it speaks MCP over stdio and stays silent until a client connects.

### User-Agent (recommended)

SEC requires a real contact in the User-Agent header for automated traffic. The default UA identifies this project; if you fork or run at any volume, set your own:

```bash
export EDGAR_USER_AGENT="my-research-tool/1.0 (you@example.com)"
```

### Wire into Claude Desktop

```json
{
  "mcpServers": {
    "cig-edgar": {
      "command": "cig-mcp-edgar",
      "env": {
        "EDGAR_USER_AGENT": "your-name (your-email@example.com)"
      }
    }
  }
}
```

Or, if `cig-mcp-edgar` isn't on Claude Desktop's `PATH`, use the absolute path:

```json
{
  "mcpServers": {
    "cig-edgar": {
      "command": "/Users/you/go/bin/cig-mcp-edgar"
    }
  }
}
```

Restart Claude Desktop. Two new tools appear: `lookup_company` and `list_filings`. Ask: *"List Apple's last five 10-Q filings"* — Claude calls `list_filings(ticker="AAPL", form="10-Q", limit=5)` and gets accession numbers and direct URLs to the filings.

### Limits

- Backed by the public SEC EDGAR endpoints (`data.sec.gov/submissions/`, `www.sec.gov/files/company_tickers.json`). SEC rate-limits abusive clients; respect their [fair access guidelines](https://www.sec.gov/os/accessing-edgar-data).
- The ticker→CIK map is loaded once per server lifetime. Restart the server if SEC adds new issuers you care about.
- Only equity/registrant tickers in `company_tickers.json` are recognized.

## Next steps

- See [CONTRIBUTING.md](../CONTRIBUTING.md) to add a new tool or server.
- File bugs or feature requests at [github.com/ChizyIG/cig-mcp-toolkit/issues](https://github.com/ChizyIG/cig-mcp-toolkit/issues).
