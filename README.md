<div align="center">

# cig-mcp-toolkit

**Model Context Protocol servers and building blocks for quantitative finance.**

Market data. Portfolio analytics. Risk metrics. Backtesting.
Built in Python and Go. Open source. Maintained by CIG Engineering.

[![CI](https://github.com/ChizyIG/cig-mcp-toolkit/actions/workflows/ci.yml/badge.svg)](https://github.com/ChizyIG/cig-mcp-toolkit/actions/workflows/ci.yml)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](./LICENSE)

</div>

---

> **Status:** v0.1 — early. Two reference servers ship today (Yahoo Finance, SEC EDGAR). The wider scope below is the roadmap, not what's installed.

## What this is

`cig-mcp-toolkit` is a set of reference [Model Context Protocol](https://modelcontextprotocol.io) servers and reusable components for building AI-integrated investment tooling. It gives you the scaffolding: auth patterns, schema conventions, transport helpers, testing utilities, and a handful of working example servers, so you can wire real-time portfolio, risk, and market data into any MCP-compatible LLM client.

It is not a trading platform, a data feed, or a modeling library. It is the plumbing between LLM applications and the kinds of systems that investment teams actually run.

## Why MCP for finance

LLMs are good at reasoning over structured context. Investment workflows are *made* of structured context: positions, exposures, factor attributions, risk limits, counterparty hierarchies, P&L attributions. MCP standardizes how a model fetches that context, similar to how a web API standardizes how a browser fetches a page, so the same set of servers can be consumed by any MCP-aware client.

The quant-finance space has almost no public MCP tooling yet. This repo is an attempt to change that, and to contribute the patterns we've developed internally at CIG back to the open source community.

## What's in v0.1

| Server | Language | Tools | Data source |
|---|---|---|---|
| `cig-mcp-yahoo` | Python | `get_quote` | Yahoo Finance (via `yfinance`) — equities, ETFs, crypto, FX, indices |
| `cig-mcp-edgar` | Go | `lookup_company`, `list_filings` | SEC EDGAR — CIK lookup, recent filings (10-K, 10-Q, 8-K, etc.) |

See [`docs/quickstart.md`](./docs/quickstart.md) for install + Claude Desktop wiring.

## Quick taste

After installing (see quickstart), point an MCP client at one of the servers and ask:

> *What's NVDA trading at right now?* &nbsp;→ `cig-mcp-yahoo` calls `get_quote("NVDA")`
>
> *List Apple's last five 10-Q filings.* &nbsp;→ `cig-mcp-edgar` calls `list_filings(ticker="AAPL", form="10-Q", limit=5)`

## Roadmap

The longer-term scope. Items below are *not* shipped yet — track progress in [issues](https://github.com/ChizyIG/cig-mcp-toolkit/issues) or open one to propose a new server.

### Core library (`cig_mcp/`)

A thin layer on top of the official SDKs that handles the concerns every finance MCP server ends up needing:

- **Schema conventions** for positions, prices, returns, risk metrics, and trade events (Pydantic models on the Python side, struct tags on the Go side)
- **Auth middleware** — API key, OAuth 2.1, and signed-request helpers
- **Streaming helpers** for server-sent prices, PnL, and risk updates
- **Rate-limit and circuit-breaker utilities** for wrapping upstream data providers
- **Testing harness** for MCP server contract tests and regression snapshots

### Planned reference servers

| Server | Purpose | Primitives |
|---|---|---|
| `market_data` | OHLCV + reference data across multiple providers (Alpha Vantage, FRED) | Tools + Resources |
| `portfolio_snapshot` | Position and exposure snapshots from a CSV/Parquet portfolio file | Resources |
| `risk_metrics` | VaR, CVaR, Sharpe, Sortino, max drawdown on a given return stream | Tools |
| `backtest_runner` | Run a strategy spec against historical data and return the equity curve | Tools + Prompts |
| `data_pipeline_monitor` | Surface job status, freshness, and quality checks for an ETL pipeline | Resources |

All reference servers use public or user-supplied data only — no proprietary signals, no internal endpoints.

## Contributing

See [`CONTRIBUTING.md`](./CONTRIBUTING.md) for development setup, the PR process, and what we will and won't merge.

Security issues: please follow [`SECURITY.md`](./SECURITY.md) — do not file them as public issues.

## License

Licensed under the Apache License, Version 2.0. See [`LICENSE`](./LICENSE) for the full text.
