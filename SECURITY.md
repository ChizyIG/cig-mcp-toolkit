# Security Policy

We take the security of `cig-mcp-toolkit` seriously. This document explains how to report vulnerabilities and what to expect after you do.

## Supported versions

The project is pre-1.0. Only the latest `main` and the most recent tagged release on the `0.x` line receive security fixes. Once `1.0.0` ships, this table will be updated with a formal support window.

| Version | Supported |
|---------|-----------|
| `main`  | Yes       |
| Latest `0.x` tag | Yes |
| Older `0.x` tags | No |

## Reporting a vulnerability

**Please do not open a public GitHub issue for a security vulnerability.** Public disclosure before a fix is available puts users at risk.

Use one of these private channels instead:

1. **Preferred — GitHub Security Advisories.** Go to the repository's [Security tab](https://github.com/ChizyIG/cig-mcp-toolkit/security/advisories/new) and open a private advisory draft. This gives us a place to collaborate with you on a fix.
2. **Email.** Send a report to `legal@chizyig.com`. If you want the message encrypted, request our PGP key in a low-detail first email and we'll reply with the key fingerprint.

### What to include

The more of the following you can provide, the faster we can triage:

- A description of the vulnerability and the component it affects (server name, package, file).
- Version or commit hash where you reproduced it.
- Steps to reproduce, minimal PoC code, or a test case.
- Impact assessment — what an attacker could do, under what assumptions.
- Your name and contact info (for the credit section of the eventual advisory, if you want credit).

## What to expect

- **Acknowledgement** within 3 business days of your report.
- **Initial triage** — severity assessment and next steps — within 7 business days.
- **Fix target** — 30 days from triage for High/Critical, 90 days for Low/Medium. We'll communicate if an issue needs longer.
- **Coordinated disclosure** — we'll agree on a disclosure date with you, request a CVE when applicable, and credit you in the advisory unless you prefer to remain anonymous.

## Scope

In scope:

- Code in this repository (Python and Go sources, schemas, CI configuration).
- Reference MCP servers under `servers/` when run as documented.
- Dependency chains we control (our pinned versions, not upstream projects themselves).

Out of scope:

- Vulnerabilities in third-party data providers (Yahoo Finance, Alpha Vantage, FRED, etc.) — report those to the respective vendor.
- Social engineering of maintainers or contributors.
- Denial-of-service attacks requiring massive unthrottled traffic.
- Configuration issues in deployments you control (e.g., you exposed a local-only server to the public internet).

## Safe harbor

If you make a good-faith effort to comply with this policy during security research, we will not pursue or support legal action against you. We consider security research conducted under this policy to be authorized and will work with you to understand and resolve the issue quickly.
