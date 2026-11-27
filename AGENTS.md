# AGENTS.md — Operating Rules for `zakshub/zirc-client`

> **Public Contributors vs. Internal Governance:** This is a **public, MIT-licensed** repository. Public contributors and external tools should follow this `AGENTS.md`, `README.md`, `SECURITY.md`, `PRIVACY.md`, and `THREAT_MODEL.md`. Internal maintainer governance and portfolio tracking follow the owner's **private** `zakshub/project-os` repository (`antigravity/OPERATING_INSTRUCTIONS.md`) and **private** `Zakshub Portfolio` board (owner access only).

---

## 1. Project Identity & Non-Negotiable Security Invariants
- **Purpose:** Classic-style, programmable IRC client in Go built desktop-first with a shared IRC engine and a companion terminal interface (`cmd/zirc`, `cmd/zirc-desktop`, `internal/`).
- **Classification:** Active Substantial Product (early prototype — not yet independently security audited)
- **Default Branch:** `main` (use short-lived feature branches + Pull Requests; never rewrite history)
- **Non-Negotiable Security & Privacy Invariants (per [SECURITY.md](./SECURITY.md) & [PRIVACY.md](./PRIVACY.md)):**
  - ZIRC connects to existing IRC networks; it does not operate an IRC server.
  - **Fail-closed SOCKS5 proxy routing:** When a SOCKS5 proxy is configured, resolve destination hostnames through the proxy and fail closed if the proxy is unavailable.
  - **No direct connection fallback:** Never fall back to a direct connection when a selected proxy fails (including during automatic reconnect).
  - **TLS 1.2 minimum & no certificate bypass:** Enforce TLS 1.2+ with certificate verification whenever TLS is enabled; never add an insecure certificate verification bypass option.
  - **SASL PLAIN only over TLS:** Reject SASL PLAIN authentication unless TLS is active.
  - **No credentials through CLI arguments:** Never accept passwords or secrets via command-line flags/arguments (use Linux Secret Service keyring / `secret-tool` or optional environment variable fallback).
  - **No unrequested telemetry or persistent chat logging:** Do not add telemetry, unrequested persistent chat logging, DCC transfers, or automatic link fetching.
  - **Preserve existing CI workflows:** Leave `.github/workflows/ci.yml` and `.github/workflows/release.yml` intact unless a dedicated CI/release task explicitly approves changes.

## 2. Mandatory Reading Order
1. [README.md](./README.md)
2. [ROADMAP.md](./ROADMAP.md) and [docs/CLASSIC_FEATURE_PARITY.md](./docs/CLASSIC_FEATURE_PARITY.md)
3. [SECURITY.md](./SECURITY.md), [PRIVACY.md](./PRIVACY.md), and [THREAT_MODEL.md](./THREAT_MODEL.md)
4. [Makefile](./Makefile), `.github/workflows/ci.yml`, and `.github/workflows/release.yml`

## 3. Verified Local Development & Test Commands

### Linux Desktop Prerequisites (Fyne Native Dependencies)
Building or testing the desktop client (`./cmd/zirc-desktop`) on Linux requires native CGo/OpenGL/X11 packages (matching `README.md` and `.github/workflows/ci.yml`):
```bash
sudo apt update
sudo apt install -y gcc libgl1-mesa-dev xorg-dev
```

### `Makefile` Shortcuts (and Exact Side Effects)
Note that `Makefile` targets are local convenience shortcuts—they do **not** run the full CI suite:
```bash
make fmt    # Runs `gofmt -w cmd internal` — WARNING: mutates Go files in place
make vet    # Runs `go vet ./...` (without `-tags ci`)
make test   # Runs `go test ./...` (without `-race`, `-tags ci`, or fuzz campaigns)
make build  # Builds ONLY the terminal binary `./zirc` (`./cmd/zirc`) — does NOT build `./cmd/zirc-desktop`
make clean  # Deletes the local `./zirc` binary
```

### Full Go CI-Equivalent Verification Commands
```bash
go test -race -tags ci ./...
go test ./internal/irc -run=^$ -fuzz=FuzzParseLine -fuzztime=5s
go test ./internal/irc -run=^$ -fuzz=FuzzFormatValidate -fuzztime=5s
go vet -tags ci ./...
go build -trimpath -ldflags="-s -w" -o /tmp/zirc ./cmd/zirc
go build -trimpath -ldflags="-s -w" -o /tmp/zirc-desktop ./cmd/zirc-desktop
```

## 4. Design Intelligence
- **Desktop / TUI interface changes:** Consult canonical UX/UI intelligence (`zakshub/designer`, via your configured local intelligence workspace when available) before major layout or interaction redesigns.

## 5. Verification & Safety Rules
- Never modify Go application source code solely for documentation or Project OS chores.
- Never commit IRC credentials, server tokens, or private chat logs.
- Never disclose exploitable security vulnerabilities in public GitHub issues—follow [SECURITY.md](./SECURITY.md).
- Never publish releases or tags without Zak's explicit approval and passing CI checks.
- Always report exact test/vet/fuzz/build outputs before claiming completion.
