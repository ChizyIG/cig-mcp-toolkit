<div align="center">

# cig-mcp-toolkit

**Model Context Protocol servers and building blocks for quantitative finance.**

Market data. Portfolio analytics. Risk metrics. Backtesting.
Built in Python and Go. Open source. Maintained by CIG Engineering.

</div>

---

## What this is

`cig-mcp-toolkit` is a set of reference [Model Context Protocol](https://modelcontextprotocol.io) servers and reusable components for building AI-integrated investment tooling. It gives you the scaffolding: auth patterns, schema conventions, transport helpers, testing utilities, and a handful of working example servers, so you can wire real-time portfolio, risk, and market data into any MCP-compatible LLM client.

It is not a trading platform, a data feed, or a modeling library. It is the plumbing between LLM applications and the kinds of systems that investment teams actually run.

## Why MCP for finance

LLMs are good at reasoning over structured context. Investment workflows are *made* of structured context: positions, exposures, factor attributions, risk limits, counterparty hierarchies, P&L attributions. MCP standardizes how a model fetches that context, similar to how a web API standardizes how a browser fetches a page, so the same set of servers can be consumed by any MCP-aware client.

The quant-finance space has almost no public MCP tooling yet. This repo is an attempt to change that, and to contribute the patterns we've developed internally at CIG back to the open source community.

## What's included

### Core library (`cig_mcp/`)

A thin layer on top of the official SDKs that handles the concerns every finance MCP server ends up needing:

- **Schema conventions** for positions, prices, returns, risk metrics, and trade events (Pydantic models on the Python side, struct tags on the Go side)
- **Auth middleware** — API key, OAuth 2.1, and signed-request helpers
- **Streaming helpers** for server-sent prices, PnL, and risk updates
- **Rate-limit and circuit-breaker utilities** for wrapping upstream data providers
- **Testing harness** for MCP server contract tests and regression snapshots

### Reference servers (`servers/`)

Each is a self-contained, runnable MCP server you can clone, point at your own data, and extend. All examples use public data only — no proprietary signals, no internal endpoints.

| Server | Purpose | Primitives |
|---|---|---|
| `market_data` | OHLCV + reference data adapter (Yahoo Finance, Alpha Vantage, FRED) | Tools + Resources |
| `portfolio_snapshot` | Position and exposure snapshots from a CSV/Parquet portfolio file | Resources |
| `risk_metrics` | VaR, CVaR, Sharpe, Sortino, max drawdown on a given return stream | Tools |
| `backtest_runner` | Run a strategy spec against historical data and return the equity curve | Tools + Prompts |
| `data_pipeline_monitor` | Surface job status, freshness, and quality checks for an ETL pipeline | Resources |

## License

Licensed under the Apache License, Version 2.0. See [`LICENSE`](./LICENSE) for the full text.
