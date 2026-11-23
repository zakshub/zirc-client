# ZIRC Client Roadmap

## Milestone 1: Linux terminal foundation

- [x] Go project structure and CLI flags
- [x] Full screen keyboard driven Linux TUI
- [x] IRC registration and basic message exchange
- [x] TLS with certificate verification
- [x] SOCKS5 proxy routing and proxy side destination DNS
- [x] Privacy and threat model documentation
- [x] Protocol parser tests
- [x] GitHub Actions build and test workflow

## Milestone 2: Privacy and identity controls

- [x] CLI identity preview before connecting
- [x] Random nickname generation
- [x] Configurable username and generic real name
- [x] SOCKS5 username/password authentication
- [x] Fail closed when selected proxy is unavailable
- [x] Explicit direct route and proxy route status
- [ ] Dedicated interactive privacy settings screen
- [ ] Tor route profile diagnostics and independent verification limits
- [ ] Automated DNS and connection leak test suite

## Milestone 3: Authentication and local protection

- [x] IRCv3 capability negotiation
- [x] SASL PLAIN authentication over TLS
- [x] Passwords accepted through environment variables rather than command line arguments
- [x] Linux Secret Service keyring lookup and explicit secret deletion\n- [ ] Native macOS and Windows keyring integration
- [ ] Optional encrypted local history
- [ ] Secure local configuration file permissions, if configuration files are introduced

## Milestone 4: Full IRC usability

- [x] Separate channel and private conversation buffers
- [x] Channel member list and NAMES refresh
- [x] Topics, modes and invites
- [ ] Multiple simultaneous IRC network connections
- [ ] Desktop notifications
- [x] Search within active in-memory conversation buffer\n- [ ] Persistent searchable history
- [ ] Terminal display cell width handling for wide and combining glyphs
- [ ] Automatic reconnect and connection diagnostics

## Milestone 5: Release engineering and assurance

- [x] Linux AMD64 and ARM64 release workflow defined
- [x] SHA256 checksum generation defined
- [x] Keyless Sigstore checksum signing defined
- [ ] Execute and verify a real tagged release
- [x] IRC protocol fuzz test harness added to CI seed corpus
- [ ] Sustained fuzz campaigns for IRC and SOCKS5 parsers
- [ ] Expanded network leak testing
- [ ] Independent security review
- [ ] Documented threat model review after security review
