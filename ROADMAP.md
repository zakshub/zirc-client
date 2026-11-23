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
- [x] Linux Secret Service keyring lookup and explicit secret deletion
- [x] Environment variable credential source as an alternative
- [ ] Native macOS and Windows keyring integration
- [ ] Optional encrypted persistent history

## Milestone 4: Full IRC usability

- [x] Separate channel and private conversation buffers
- [x] Channel member list and NAMES refresh
- [x] Topics, modes and invites
- [x] Search within active in memory conversation buffer
- [x] Bounded automatic reconnect with exponential backoff and same route enforcement
- [x] Terminal display cell width handling for wide glyphs, combining marks and common joined emoji sequences
- [x] Persistent command aliases with bounded positional expansion
- [x] IRC moderation commands: kick, ban, unban, op, deop, voice and devoice
- [x] In memory command history and keyboard recall
- [x] Tab completion for commands, aliases, nicknames and channels
- [ ] Multiple simultaneous IRC network connections
- [ ] Desktop notifications
- [ ] Persistent searchable history
- [ ] Dedicated interactive privacy settings screen
- [ ] More detailed connection diagnostics and user controlled retry policy

## Milestone 5: Release engineering and assurance

- [x] Linux AMD64 and ARM64 release workflow defined
- [x] SHA256 checksum generation defined
- [x] Keyless Sigstore checksum signing defined
- [x] IRC protocol fuzz test harness added to CI seed corpus
- [x] Short CI fuzz campaigns for IRC line parsing and formatting
- [ ] Dedicated SOCKS5 handshake fuzzing and longer fuzz campaigns
- [ ] Expanded network leak testing
- [ ] Execute and verify a real tagged release
- [ ] Independent security review
- [ ] Documented threat model review after security review
