"""Unit tests for the Yahoo Finance MCP server's get_quote tool.

The yfinance library is mocked so these tests run offline.
"""

import asyncio
from types import SimpleNamespace
from unittest.mock import MagicMock, PropertyMock, patch

import pytest

from cig_mcp.servers.yahoo.server import get_quote, mcp


def _info(**fields: object) -> SimpleNamespace:
    base: dict[str, object] = {
        "last_price": None,
        "previous_close": None,
        "currency": None,
        "exchange": None,
        "market_cap": None,
    }
    base.update(fields)
    return SimpleNamespace(**base)


def test_returns_expected_shape():
    ticker = MagicMock(
        fast_info=_info(
            last_price=150.0,
            previous_close=148.0,
            currency="USD",
            exchange="NMS",
            market_cap=2_500_000_000_000,
        )
    )

    with patch("cig_mcp.servers.yahoo.server.yf.Ticker", return_value=ticker):
        result = get_quote("AAPL")

    assert result == {
        "symbol": "AAPL",
        "last_price": 150.0,
        "previous_close": 148.0,
        "change": pytest.approx(2.0),
        "change_percent": pytest.approx((2.0 / 148.0) * 100.0),
        "currency": "USD",
        "exchange": "NMS",
        "market_cap": 2_500_000_000_000,
    }


def test_normalizes_symbol_case_and_whitespace():
    ticker = MagicMock(fast_info=_info(last_price=10.0, previous_close=10.0))

    with patch("cig_mcp.servers.yahoo.server.yf.Ticker", return_value=ticker) as cls:
        result = get_quote("  msft  ")

    cls.assert_called_once_with("MSFT")
    assert result["symbol"] == "MSFT"
    assert result["change"] == 0.0
    assert result["change_percent"] == 0.0


@pytest.mark.parametrize("bad", ["", "   ", "\t\n"])
def test_rejects_empty_symbol(bad):
    with pytest.raises(ValueError, match="non-empty"):
        get_quote(bad)


def test_raises_when_yahoo_returns_no_data():
    ticker = MagicMock(fast_info=_info())

    with (
        patch("cig_mcp.servers.yahoo.server.yf.Ticker", return_value=ticker),
        pytest.raises(ValueError, match="no quote data"),
    ):
        get_quote("ZZZZ")


def test_wraps_upstream_errors():
    ticker = MagicMock()
    type(ticker).fast_info = PropertyMock(side_effect=RuntimeError("yahoo down"))

    with (
        patch("cig_mcp.servers.yahoo.server.yf.Ticker", return_value=ticker),
        pytest.raises(ValueError, match="failed to fetch"),
    ):
        get_quote("AAPL")


def test_tool_is_registered_with_mcp_server():
    tools = asyncio.run(mcp.list_tools())
    assert any(t.name == "get_quote" for t in tools)
