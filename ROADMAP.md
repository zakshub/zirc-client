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


## Milestone 6: Desktop application foundation

- [x] Add a separate desktop application entry point
- [x] Introduce a reusable desktop workspace package
- [x] Build a classic IRC layout with network and channel navigation
- [x] Add server connection dialog with TLS selection
- [x] Reuse the existing IRC protocol client for desktop connections
- [x] Add transcript, topic, nickname list and message composer
- [x] Add common IRC commands to the desktop composer
- [x] Add Linux desktop build dependencies and CI build
- [ ] Add multiple simultaneous network sessions
- [ ] Persist server profiles and auto join settings
- [ ] Add SASL and SOCKS5 settings to the desktop connection dialog
- [ ] Add reconnect controls and connection diagnostics
- [ ] Add GUI integration tests

## Milestone 7: Classic scripting environment

- [ ] Define the ZIRC script language and compatibility goals
- [ ] Implement aliases and command expansion
- [ ] Implement remote event handlers
- [ ] Implement variables, identifiers and scoped state
- [ ] Implement timers and scheduled actions
- [ ] Implement popup menu definitions
- [ ] Implement script loading, unloading and dependency handling
- [ ] Implement a script editor with syntax highlighting and diagnostics
- [ ] Add a permission manifest and sandbox for scripts
- [ ] Add compatibility import and migration tools for historical script packs

## Milestone 8: Classic IRC feature suite

- [ ] CTCP tools and configurable responses
- [ ] DCC chat and file transfer with explicit user approval
- [ ] Ignore lists and configurable user protection
- [ ] Channel and nickname tracking
- [ ] Custom dialogs, windows, toolbars and popup menus
- [ ] Sound and desktop notification events
- [ ] Bot and addon framework with permission controls
- [ ] Theme and layout customization

## Milestone 9: Release and platform support

- [ ] Linux desktop package
- [ ] Windows desktop package
- [ ] Portable configuration export and import
- [ ] Script compatibility test corpus
- [ ] Full end to end IRC network test suite
- [ ] Independent security review
