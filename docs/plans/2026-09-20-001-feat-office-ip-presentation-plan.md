---
title: Configurable Office IP Labels in FindPlace
type: feat
status: active
date: 2026-09-20
origin: docs/brainstorms/office-ip-presentation.md
---

# Configurable Office IP Labels in FindPlace

## Overview

Add a single envconfig field for office IP → label mappings (single IP and CIDR), apply it inside `FindPlace` before the existing `ipip` fallback, and leave everything else alone.

## Problem Frame

`pkg/web/template.go:99` — `FindPlace` does an `ipip.FindCity` lookup and returns city/province. Office traffic can only be resolved to city granularity, so two offices in the same city are indistinguishable on the status page. Operators have no way to declare "these IPs are office X." (See origin doc for full motivation.)

## Requirements Trace

- R1. Config field for `ip:label` mapping, supporting single IP and CIDR entries.
- R2. `FindPlace` returns the configured label on match, skipping `ipip`.
- R3. `FindPlace` falls back to existing `ipip.FindCity` behavior on miss.
- R4. Empty/unset config leaves behavior identical to current.

## Scope Boundaries

- `IsOfficeIP` (the `pkg/web/template.go` TODO stub) stays untouched.
- No match modes beyond single IP and CIDR.
- No config hot reload — same `init()` semantics as the rest of the config.
- No `ipip`-driven label backfill for unconfigured IPs.

### Deferred to Follow-Up Work

- Wiring `IsOfficeIP` to the same config (deferred by brainstorm).
- Hot reload of office mappings (deferred by brainstorm).

## Context & Research

### Relevant Code and Patterns

- `pkg/settings/config.go` — config struct loaded via `envconfig.Process`; existing slice fields (`ManageAddrs []string`) show the envconfig pattern. New field is `map[string]string` to match the brainstorm's chosen `ip:label` syntax.
- `pkg/ipip/private.go` — precedent for a package-level slice of `*net.IPNet` built at init and a simple linear scan (`block.Contains(ip)`). The new lookup will follow the same shape.
- `pkg/web/template.go:99` — `FindPlace` is the entry point. It is exposed to templates as `findPlace` (`pkg/web/template.go:53`) and is currently called from `ui/templates/status.html:41`.
- `ui/templates/status.html:38-42` — the template already has an `isOffice` branch (currently dead because `IsOfficeIP` always returns false). It will not be touched; the live branch today calls `findPlace`, which is what this plan modifies.

### Institutional Learnings

- No `docs/solutions/` content surfaced for this area.

### External References

- None needed. Standard library `net.ParseIP`, `net.ParseCIDR`, `net.IPNet.Contains` cover the matching. envconfig's `map[string]string` decoder uses the `key:value,key:value` syntax already chosen in the brainstorm.

## Key Technical Decisions

- New field is `map[string]string` keyed by IP/CIDR string, valued by label. Matches the format the user gave in the brainstorm (`ip1:Beijing,ip2:Suzhou`).
- Parsing happens once at init; each map entry is classified as either a single `net.IP` or a `*net.IPNet` and stored in a package-level slice in `pkg/web`. Match is a linear scan — office count is small (single digits), so no indexing needed.
- Parsing lives in `pkg/web/office.go` rather than `pkg/settings`, so `pkg/settings` stays a pure config container (matches existing pattern — only `Usage()` is a method today).
- Overlapping CIDR entries resolve to last-wins on the lookup (consistent with `map[string]string` overwrite semantics). Malformed entries are skipped with a startup log warning rather than crashing.

## Open Questions

### Resolved During Planning

- **Where does parse/lookup live?** → `pkg/web/office.go`, next to the consumer (`FindPlace`).
- **How to handle malformed entries?** → skip with log warning, do not crash init.

### Deferred to Implementation

- Exact log format for malformed-entry warning (will follow the existing pattern in `pkg/ipip/ip_find.go`).
- Whether the init parse should run lazily (sync.Once on first lookup) or eagerly at package init — eager is simpler and matches the existing codebase's style; if a future startup-cost concern appears, revisit.

## Implementation Units

### U1. Add office IP config field and lookup helper

**Goal:** Operators can declare office IP → label mappings through one env var; the program parses them once and exposes a lookup function.

**Requirements:** R1, R4

**Dependencies:** None

**Files:**
- Modify: `pkg/settings/config.go`
- Create: `pkg/web/office.go`
- Test: `pkg/web/office_test.go`

**Approach:**
- Add `OfficeIPs map[string]string` to `settings.Config` with envconfig tag `OFFICE_IPS`.
- In `pkg/web/office.go`, define an internal entry type holding either a single `net.IP` or a `*net.IPNet` plus the label.
- On package init (or via `sync.Once` on first call), iterate `settings.Current.OfficeIPs`. For each `ip, label` pair, attempt `net.ParseCIDR(ip)` first; if that succeeds, store as CIDR. Otherwise attempt `net.ParseIP(ip)`; if that succeeds, store as single IP. If both fail, log a warning and skip the entry.
- Export a `LookupOfficeLabel(ip string) (label string, ok bool)` function that scans the parsed slice and returns the first matching label.

**Patterns to follow:**
- `pkg/ipip/private.go` — package-level slice + linear `Contains` scan.
- `pkg/ipip/ip_find.go` — log-and-skip pattern when external data is missing/invalid.
- `pkg/web/template.go:27-32` — `init()` reading from `settings.Current`.

**Test scenarios:**
- Covers AE1. Single-IP entry `1.2.3.4:Beijing` parses and `LookupOfficeLabel("1.2.3.4")` returns `("Beijing", true)`.
- Covers AE2. CIDR entry `10.0.0.0/24:OfficeNet` parses and `LookupOfficeLabel("10.0.0.5")` returns `("OfficeNet", true)`.
- Lookup against an unconfigured IP returns `( "", false)`.
- Malformed entry (`not-an-ip:Foo`) is skipped without panic and without surfacing in the lookup.
- IPv6 single-IP entry parses and matches.
- Last-wins on overlapping CIDR entries (`10.0.0.0/16:Outer` and `10.0.0.0/24:Inner` → `10.0.0.5` returns `"Inner"`).

**Verification:**
- `go build ./...` succeeds.
- `go vet ./...` clean.
- New unit tests pass; `pkg/status` test still passes.

### U2. Wire FindPlace to the office IP lookup

**Goal:** `FindPlace` returns the configured label on hit, falls through to existing `ipip` logic on miss.

**Requirements:** R2, R3, R4

**Dependencies:** U1

**Files:**
- Modify: `pkg/web/template.go`
- Create: `pkg/web/template_test.go`

**Approach:**
- At the top of `FindPlace(ip string)`, call `LookupOfficeLabel(ip)`. If `ok`, return `label` immediately.
- Otherwise the existing `ipip.FindCity` path runs unchanged, including the `city == ""` → `"[未知地区]"` branch.

**Patterns to follow:**
- `pkg/web/template.go:99-108` — current `FindPlace` shape (early-return then fall-through).

**Test scenarios:**
- Covers AE1, AE3. With `settings.Current.OfficeIPs` set to `1.2.3.4:Beijing`, `FindPlace("1.2.3.4")` returns `"Beijing"` and `FindPlace("8.8.8.8")` returns the `ipip` result (or `"[未知地区]"` when `IPIP_DATX_PATH` is unset, which is the env in `go test`).
- Covers AE4. With `settings.Current.OfficeIPs` unset, behavior is byte-identical to the pre-change function.
- Covers AE2. CIDR entry set, `FindPlace("10.0.0.5")` returns `"OfficeNet"`.
- Order check: when the lookup hits, `ipip.FindCity` is never invoked (verifiable via a test that does not load the ipip data file and would otherwise return `"[未知地区]"`).

**Verification:**
- All four AEs from the origin doc are exercised by automated tests.
- `go test ./...` green.
- Manual smoke: status page renders office connections with the configured label.

## System-Wide Impact

- **Interaction graph:** Only `findPlace` template helper changes. No middleware, no routes, no other template helpers touched.
- **Error propagation:** Malformed config entries are logged at startup, not propagated to requests. A bad entry silently skips — operators should monitor startup logs.
- **State lifecycle risks:** None. Parsed entries are immutable after init.
- **API surface parity:** No public CLI/HTTP API changes. `FindPlace` keeps its signature.
- **Integration coverage:** The order-of-operations test (lookup hit never invokes `ipip`) is the only behavior visible only when both layers cooperate, and the AE tests exercise it.

## Risks & Dependencies

- Operators must understand envconfig's `map[string]string` syntax (`OFFICE_IPS=1.2.3.4:Beijing,10.0.0.0/24:OfficeNet`) — same syntax as every other `map[string]string` env var would be in this codebase, but worth flagging in the PR description.
- Startup log noise on malformed entries is the only operator-visible failure mode; acceptable trade-off for not crashing on bad config.

## Documentation / Operational Notes

- `.env.example` should gain one commented line showing the `OFFICE_IPS` syntax (single IP, CIDR, mixed).
- PR description should call out the envconfig `map[string]string` syntax explicitly so operators don't try `OFFICE_IPS=1.2.3.4 Beijing,5.6.7.8 Suzhou`.

## Sources & References

- **Origin document:** `docs/brainstorms/office-ip-presentation.md`
- **Existing IP-list pattern:** `pkg/ipip/private.go`
- **Config field pattern:** `pkg/settings/config.go`
- **Consumer to modify:** `pkg/web/template.go` (`FindPlace`)
