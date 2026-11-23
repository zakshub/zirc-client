# ZIRC Classic IRC Feature Parity

## Product decision

ZIRC is a modern IRC client that preserves the breadth and configurability of classic mIRC era clients, script packs, addons, bots, and channel utilities.

The first release is terminal only. The graphical interface is explicitly deferred until the terminal client and shared core are stable.

Repository: https://github.com/zakshub/zirc-client

## Existing foundation

The current prototype already provides:

- Full screen keyboard driven terminal interface
- IRC registration, channels, private conversations, per channel buffers
- Member list, topics, NAMES refresh, invites, MODE commands
- TLS with certificate verification
- SOCKS5 routing with proxy side DNS and no direct fallback
- IRCv3 capability negotiation and SASL PLAIN over TLS
- Random nickname generation and identity preview
- Linux Secret Service keyring integration
- In memory search
- Bounded reconnect with exponential backoff
- Unicode terminal cell width handling
- Parser fuzz tests in CI

Do not replace or regress these capabilities.

## Product architecture

Keep the IRC protocol engine independent from the user interface.

Target layers:

1. Core: connection lifecycle, IRC protocol, IRCv3, networks, identities, channels, users, events.
2. Command system: built in commands, aliases, command completion, command history, contextual help.
3. Script runtime: permissioned addon API, event hooks, timers, commands, configuration, sandbox boundaries.
4. Feature modules: channel tools, moderation, notifications, logging, themes, automation, transfers.
5. Terminal UI: first supported interface.
6. GUI: later interface consuming the same core and command APIs.

No script or addon should receive unrestricted access to the operating system by default.

## Classic feature catalogue

### Client and interface

- Multiple network connections and independent network profiles
- Multiple server tabs and channel/private conversation buffers
- Status window and server event stream
- Nick list and channel topic panel
- Split views and configurable layouts where terminal constraints permit
- Tab completion for nicknames, channels and commands
- Configurable key bindings
- Command history and input recall
- Highlighting, timestamps and configurable message formats
- Themes, colour palettes and per network appearance
- Notifications and keyword alerts
- Away status and away messages
- Ignore lists and configurable filtering
- User and channel information panels
- Help browser and command reference

### IRC commands and protocol

- Standard IRC commands including JOIN, PART, QUIT, NICK, PRIVMSG, NOTICE, TOPIC, MODE, KICK, INVITE, WHO, WHOIS, WHOWAS, NAMES, LIST, MOTD, AWAY and PING
- IRCv3 capability negotiation
- SASL authentication
- Account and extended join metadata when supported
- Message tags, server time, batch and multiline support where negotiated
- Network specific command profiles
- Raw command console with dangerous command confirmation
- Server feature discovery and graceful fallback

### mIRC style scripting and customization

- User defined aliases
- Custom commands
- Event handlers for connect, disconnect, join, part, quit, nick change, messages, notices, kicks, mode changes and topic changes
- Timers and scheduled actions
- Variables and scoped state
- Reusable script modules
- Script enable, disable and reload
- Configuration files and import/export
- Script dependency metadata
- Script permissions and capability declarations
- Debug console and script error reporting
- Script package manifest and versioning

### Popups, menus and interaction

- Context menus for users, channels and messages
- Configurable nick list actions
- Channel list actions
- URL and text context actions
- Custom menu definitions from trusted addons
- Keyboard accessible alternatives to pointer menus

### Channel administration

- Operator and voice status tools
- Topic management
- Invite, kick and ban workflows
- Ban list inspection and editing
- Channel mode helper
- User lookup and WHOIS shortcuts
- Nick tracking and notifications
- Moderation presets with explicit confirmation
- Audit trail of actions initiated by the local client
- Safe rate limits for repeated actions

### Bots and automation

- Local automation rules
- Auto join and auto identify workflows
- Auto response rules with opt in controls
- Away responder
- Highlight and notification rules
- Channel presence monitoring
- NickServ and ChanServ helper workflows where network services support them
- External bot connection profiles
- Bot command interface with permission boundaries

### Addons and utilities

- Calculator and unit conversion commands
- Dictionary and text formatting helpers
- URL detection and safe link opening
- Clipboard integration when explicitly enabled
- Search within current session
- Optional encrypted local history
- Export selected conversation data
- Channel and nickname watch lists
- Network diagnostics
- Connection and proxy diagnostics
- Unicode and text conversion helpers

### File transfer and media

- Explicitly initiated file transfer support only after a separate security design
- No automatic downloads
- No automatic execution of received files
- Clear transfer consent, size limits and destination controls
- Optional integration with external media tools in a later release
- No embedded media playback in the terminal MVP

DCC and other peer to peer transfer features are not part of the initial MVP. They require separate threat modelling because they can expose IP addresses and accept untrusted files.

## Security boundaries

ZIRC must prioritize privacy, but must never claim that TLS alone hides the user's IP address.

Required safeguards:

- TLS certificate verification by default
- SOCKS5 proxy routing with remote DNS and fail closed behavior
- Never silently fall back from proxy to direct connection
- Clear route and identity preview before connecting
- Random nickname option
- Minimize real name and username disclosure
- Do not collect telemetry
- Do not fetch remote URLs automatically
- Do not automatically download or execute files
- Secrets must not be stored in plaintext configuration
- Keyring integration where available
- Script permissions default to deny
- Scripts cannot access files, shell, network destinations or secrets unless a narrowly scoped permission is granted
- Resource limits for scripts, timers and event handlers
- Flood and abuse prevention limits
- Confirm destructive moderation actions
- Redact secrets from logs and diagnostics
- Security tests for proxy failure, DNS routing and accidental direct connections

An IRC server can still observe the connecting IP when the client connects directly. A proxy operator can observe the user's source IP. A server cloak generally hides the address from ordinary users, not from network operators.

## Explicit exclusions

ZIRC will not ship automated flooding, mass cloning, credential theft, stealth surveillance, malware execution, destructive channel takeover or harassment features.

Where historical scripts offered aggressive functions, ZIRC should provide safe administrative equivalents such as rate limited moderation, connection diagnostics, controlled test messages and defensive detection.

## Terminal MVP acceptance criteria

The first public MVP is complete only when:

- Linux terminal client builds from a documented clean checkout
- Existing TLS, SOCKS5, SASL, keyring and fail closed behavior pass regression tests
- Multiple IRC networks can be connected and switched independently
- Core IRCv3 capabilities are negotiated correctly
- Built in commands have discoverable help and tab completion
- Aliases and event hooks work with a restricted script API
- Script permissions are visible and default deny
- Themes and key bindings are configurable
- Channel moderation tools include confirmations and rate limits
- No persistent message logging occurs unless the user explicitly enables it
- Tests cover parser, reconnect, proxy routing, permissions and command execution
- README documents limitations honestly

## Delivery sequence

### M0: Baseline and architecture
Audit existing packages, tests, protocol coverage and security boundaries. Preserve the working prototype.

### M1: Multi network core
Independent network sessions, server profiles, network switching and per network state.

### M2: Command experience
Command registry, help, completion, history, aliases and contextual actions.

### M3: Classic interaction
Events, notifications, ignore lists, away tools, nick tracking, themes and key bindings.

### M4: Script and addon foundation
Restricted script runtime, permissions, timers, event subscriptions, package manifests and diagnostics.

### M5: Administration and utilities
Channel moderation workflows, service helpers, network diagnostics and safe automation.

### M6: Hardening and release
Regression tests, fuzzing, proxy and DNS leak tests, packaging, release documentation and independent security review.

### G1: GUI
Begin only after the terminal core APIs are stable. Reuse the same protocol engine, command registry, script API and settings model.

## Release principle

Do not attempt to implement every historical feature in one unreviewed change. Build a stable core, add feature modules in small testable milestones, and preserve a traceable feature parity checklist.
