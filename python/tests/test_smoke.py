"""Smoke tests ensuring the package imports and exposes its version."""

import cig_mcp


def test_version_is_nonempty() -> None:
    assert isinstance(cig_mcp.__version__, str)
    assert cig_mcp.__version__


def test_version_is_semver() -> None:
    parts = cig_mcp.__version__.split(".")
    assert len(parts) == 3, f"expected MAJOR.MINOR.PATCH, got {cig_mcp.__version__!r}"
    assert all(p.isdigit() for p in parts), cig_mcp.__version__
