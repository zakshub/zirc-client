# AGENTS.md — Operating Rules for `zakshub/zirc-client`

> **Project Operating System:** This repository follows **[zakshub/project-os](https://github.com/zakshub/project-os)** (Private).  
> Follow all 14 steps of `zakshub/project-os/antigravity/OPERATING_INSTRUCTIONS.md` for every development task.

---

## 1. Project Identity & Security Boundaries
- **Purpose:** Classic-style, programmable IRC client in Go built desktop-first with a shared IRC engine and a companion terminal interface (`cmd/zirc`, `cmd/zirc-desktop`, `internal/`).
- **Classification:** Active Substantial Product
- **Default Branch:** `main` (use short-lived feature branches + Pull Requests; never rewrite history)
- **Strict Security & Scope Constraints:**
  - ZIRC connects to existing IRC networks; it does not operate an IRC server.
  - Security, privacy, and threat-model constraints are first-class—never weaken TLS, input sanitization, parser safety, or secret handling.
  - Preserve existing working CI workflows (`.github/workflows/ci.yml` and `.github/workflows/release.yml`).

## 2. Mandatory Reading Order
1. [README.md](./README.md)
2. [ROADMAP.md](./ROADMAP.md) and [docs/CLASSIC_FEATURE_PARITY.md](./docs/CLASSIC_FEATURE_PARITY.md)
3. [SECURITY.md](./SECURITY.md), [PRIVACY.md](./PRIVACY.md), and [THREAT_MODEL.md](./THREAT_MODEL.md)
4. [Makefile](./Makefile) and `.github/workflows/ci.yml`

## 3. Verified Local Development & Test Commands
```bash
# Format, vet, test, and build (matches Makefile and Go CI)
make fmt
make vet
make test
make build

# Full CI-equivalent checks (race detector, parser fuzz seeds, static analysis, desktop & TUI builds)
go test -race -tags ci ./...
go test ./internal/irc -run=^$ -fuzz=FuzzParseLine -fuzztime=5s
go test ./internal/irc -run=^$ -fuzz=FuzzFormatValidate -fuzztime=5s
go vet -tags ci ./...
go build -trimpath -ldflags="-s -w" -o /tmp/zirc ./cmd/zirc
go build -trimpath -ldflags="-s -w" -o /tmp/zirc-desktop ./cmd/zirc-desktop
```

## 4. Design & Visual Intelligence
- **Desktop / TUI layout changes:** Consult `/root/.intelligence/designer` (`zakshub/designer`) before major interface redesigns.

## 5. Verification & Safety Rules
- Never modify application source code solely for documentation or Project OS chores.
- Never commit IRC credentials, server tokens, or private chat logs.
- Never publish releases or tags without Zak's explicit approval and passing CI checks.
- Always report exact test/vet/build outputs before claiming completion.
