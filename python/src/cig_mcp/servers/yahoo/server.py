# Copyright 2026 CIG Engineering
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""Yahoo Finance MCP server.

Exposes one tool, ``get_quote``, that returns the latest price snapshot
for a ticker via the public Yahoo Finance endpoints (through ``yfinance``).
"""

from typing import Any

import yfinance as yf
from mcp.server.fastmcp import FastMCP

mcp: FastMCP = FastMCP("cig-mcp-yahoo")


@mcp.tool()
def get_quote(symbol: str) -> dict[str, Any]:
    """Fetch the latest price snapshot for a stock symbol.

    Args:
        symbol: Ticker symbol such as ``AAPL``, ``MSFT``, or ``BRK-B``.
            Case-insensitive; surrounding whitespace is trimmed.

    Returns:
        A dict containing ``symbol``, ``last_price``, ``previous_close``,
        ``change``, ``change_percent``, ``currency``, ``exchange``, and
        ``market_cap``. Numeric fields are floats; string/None fields may
        be ``None`` if Yahoo does not return them for the symbol.

    Raises:
        ValueError: If ``symbol`` is empty, cannot be resolved, or the
            upstream call fails.
    """
    cleaned = symbol.strip().upper() if symbol else ""
    if not cleaned:
        raise ValueError("symbol must be a non-empty string")

    ticker = yf.Ticker(cleaned)
    try:
        info = ticker.fast_info
        last_price = info.last_price
        previous_close = info.previous_close
    except Exception as exc:
        raise ValueError(f"failed to fetch quote for {cleaned!r}: {exc}") from exc

    if last_price is None or previous_close is None:
        raise ValueError(f"no quote data available for {cleaned!r}")

    try:
        last_price_f = float(last_price)
        previous_close_f = float(previous_close)
    except (TypeError, ValueError) as exc:
        raise ValueError(f"non-numeric quote data for {cleaned!r}: {exc}") from exc

    change = last_price_f - previous_close_f
    change_percent = (change / previous_close_f) * 100.0 if previous_close_f else 0.0

    return {
        "symbol": cleaned,
        "last_price": last_price_f,
        "previous_close": previous_close_f,
        "change": change,
        "change_percent": change_percent,
        "currency": getattr(info, "currency", None),
        "exchange": getattr(info, "exchange", None),
        "market_cap": getattr(info, "market_cap", None),
    }


def main() -> None:
    """Entry point for the ``cig-mcp-yahoo`` console script."""
    mcp.run()
