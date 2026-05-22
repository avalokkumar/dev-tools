# ADR 0006 — netguard private-IP guard

**Status:** Accepted
**Date:** 2026-05-10

## Context

Phase E adds an HTTP request runner (`http_request`) and a DNS lookup
operation (`dns_lookup`). DevForge runs as an MCP server invoked by AI
coding agents. Without an outbound-network guard, an agent could be tricked
(by a prompt-injection attack on a third-party document, etc.) into using
DevForge to probe the host's internal network or cloud metadata endpoints
(169.254.169.254). Classic SSRF accelerant.

## Decision

A new package `internal/netguard` exists for one purpose: classify an IP
address as private/public and refuse to dial it unless the caller passes
`allowPrivate: true` explicitly.

Rules (treated as private):

- IPv4 RFC1918 (10/8, 172.16/12, 192.168/16)
- CGNAT 100.64.0.0/10
- IPv4 link-local 169.254.0.0/16
- IPv6 link-local fe80::/10 and unique-local fc00::/7
- Loopback (127/8 and ::1)
- Multicast and unspecified (0.0.0.0, ::)

`http_request` resolves the host before dialing AND re-checks at dial time
(defense in depth against DNS rebinding between resolve and dial).

## Consequences

- Default behaviour is safe: agents cannot probe internal services without
  the user explicitly setting `allowPrivate: true`.
- Power users on internal networks can opt in per call.
- Engine emits `HTTP.PRIVATE_BLOCKED` / `DNS.PRIVATE_BLOCKED` diagnostics
  with the resolved IPs so callers see WHY they were blocked.
- All blocked decisions are observable in MCP `tools/call` responses.
