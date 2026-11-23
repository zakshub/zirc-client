# ZIRC Client Roadmap

## Milestone 1: Linux terminal foundation

- [x] Go project structure and CLI flags
- [x] Full-screen keyboard-driven Linux TUI
- [x] IRC registration and basic message exchange
- [x] TLS with certificate verification
- [x] SOCKS5 proxy routing
- [x] Proxy-side destination hostname resolution
- [x] Privacy and threat model documentation
- [x] Protocol parser tests
- [x] SOCKS5 remote DNS test
- [x] GitHub Actions build and test workflow

## Milestone 2: Privacy and identity controls

- [ ] Dedicated privacy settings screen
- [ ] Dedicated interactive privacy settings screen
- [x] CLI identity preview before connecting
- [x] Random locally generated nick option
- [x] Configurable username and realname minimisation
- [ ] Explicit Tor route profile and independent route verification
- [x] SOCKS5 username/password authentication
- [ ] Expanded connection leak and DNS behaviour test suite
- [x] Prevent accidental direct fallback when proxy is selected

## Milestone 3: Authentication and local protection

- [x] IRCv3 capability negotiation
- [x] SASL PLAIN authentication over TLS
- [ ] OS keyring integration for secrets
- [x] Proxy and SASL secrets are accepted through environment variables rather than command line arguments
- [ ] OS keyring and secret lifecycle controls
- [ ] Optional encrypted local history
- [ ] Secure config file permissions

## Milestone 4: Full IRC usability

- [ ] Multiple network connections
- [x] Channel and user list panels
- [x] Topics, modes and invites
- [ ] Desktop notifications
- [ ] Searchable in-memory history
- [x] Unicode rune input and rune-safe clipping
- [ ] Terminal display-cell width handling for wide glyphs
- [ ] Reconnect policy and connection diagnostics

## Milestone 5: Release engineering

- [ ] Linux packages and reproducible release builds
- [ ] Signed release checksums
- [ ] Security review
- [ ] Fuzz testing of IRC and SOCKS5 parsers
- [ ] Documented threat model review
