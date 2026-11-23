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
- [ ] Identity preview before connecting
- [ ] Random locally generated nick option
- [ ] Configurable username and realname minimisation
- [ ] Explicit Tor route profile and route status
- [ ] SOCKS5 username/password authentication
- [ ] Connection leak and DNS behaviour test suite
- [ ] Prevent accidental direct fallback when proxy is selected

## Milestone 3: Authentication and local protection

- [ ] IRCv3 capability negotiation
- [ ] SASL authentication
- [ ] OS keyring integration for secrets
- [ ] No secrets in command history or process arguments
- [ ] Optional encrypted local history
- [ ] Secure config file permissions

## Milestone 4: Full IRC usability

- [ ] Multiple network connections
- [ ] Channel and user list panels
- [ ] Topics, modes, invites and notifications
- [ ] Searchable history
- [ ] Unicode-aware terminal editing and rendering
- [ ] Reconnect policy and connection diagnostics

## Milestone 5: Release engineering

- [ ] Linux packages and reproducible release builds
- [ ] Signed release checksums
- [ ] Security review
- [ ] Fuzz testing of IRC and SOCKS5 parsers
- [ ] Documented threat model review
