# ZIRC Client

A classic style, programmable IRC client in active development. ZIRC is being built as a desktop first application with a shared IRC engine and a companion terminal interface.

ZIRC connects to existing IRC networks. It does not operate an IRC network or server.

## Implemented prototype capabilities

- Full screen keyboard driven Linux TUI
- IRC registration, channel messages, private conversations and per channel buffers
- Channel member list, topics, NAMES refresh, invites and MODE commands
- TLS with certificate verification and TLS 1.2 minimum
- SOCKS5 routing with proxy side destination DNS and optional username/password authentication
- IRCv3 capability negotiation and SASL PLAIN over TLS
- Random local nickname generation and configurable username and real name
- Identity and route preview before connecting
- Linux Secret Service keyring lookup and explicit secret deletion through secret-tool
- Optional environment variable credential source
- In memory search of the active conversation buffer
- Persistent mIRC style command aliases with positional argument expansion and bounded recursion
- Bounded automatic reconnect using the original configured route, with exponential backoff and no direct fallback
- CJK, combining mark and common joined emoji terminal cell width handling
- Short IRC parser fuzz campaigns in CI
- No telemetry, persistent chat logging, DCC, or automatic link fetching
- Direct versus configured SOCKS5 route status shown in the TUI

This remains an early prototype. It is not security audited and is not feature complete.

## Privacy model

TLS encrypts traffic but does not hide your IP address from the IRC server. To hide your public IP from the IRC server, route through Tor or a trusted proxy. An IRC host cloak generally hides your address from other users, not from network operators.

ZIRC can confirm that a SOCKS5 proxy connection was used. It cannot independently prove that the proxy is Tor. A proxy operator can observe your source IP and destination. Tor does not protect against account correlation, endpoint compromise, global traffic correlation, or identifying information you disclose in chat.

When a proxy is selected, ZIRC fails closed if the proxy connection fails. Automatic reconnect reuses that same configured route. It never falls back to a direct connection.

## Build

Requires Go 1.23 or newer.

Terminal client:

```bash
go test -race ./...
go vet ./...
go build -trimpath -o zirc ./cmd/zirc
```

### Linux desktop client

The desktop client uses Fyne. On Linux, install the native build dependencies first:

```bash
sudo apt update
sudo apt install -y gcc libgl1-mesa-dev xorg-dev
```

Build and launch the desktop application:

```bash
go mod tidy
go build -trimpath -o zirc-desktop ./cmd/zirc-desktop
./zirc-desktop
```

The initial desktop workspace provides a server connection dialog, server and channel navigation, conversation transcript, nickname list, message composer, and common IRC commands. It currently supports one active network connection. Multiple simultaneous networks, the complete scripting language, script editor, popup customization, DCC, plugin management, and the remaining classic client features are still under development.

## Connect

Direct TLS connection. The IRC network can still see your public IP:

```bash
./zirc --server irc.libera.chat:6697 --tls --nick zirc-user
```

Route through a locally running Tor SOCKS5 service, commonly 127.0.0.1:9050:

```bash
./zirc --server irc.libera.chat:6697 --tls --proxy 127.0.0.1:9050 --nick zirc-user
```

This confirms use of a SOCKS5 endpoint only. It does not verify that the endpoint is Tor. The SOCKS5 implementation sends the IRC server hostname to the proxy for name resolution and does not locally resolve that hostname.

SOCKS5 username/password authentication is not encrypted between the client and proxy. Use it only with a trusted local proxy or a separately protected tunnel.

## Credential storage

For Linux desktop environments, ZIRC can retrieve secrets from the Secret Service keyring through secret-tool from the libsecret-tools package.

Store a SASL password interactively:

```bash
sudo apt install libsecret-tools
secret-tool store --label='ZIRC SASL' service zirc-client account sasl:account
```

Connect using the stored secret:

```bash
./zirc --server irc.example.net:6697 --tls --sasl-user account --sasl-password-keyring
```

For SOCKS5, store the secret under account proxy:username and run with --proxy-user username --proxy-password-keyring.

Remove stored secrets:

```bash
./zirc --forget-sasl-password --sasl-user account
./zirc --forget-proxy-password --proxy-user username
```

Keyring integration currently targets Linux Secret Service. Environment variables remain available as an alternative, but they are not a secure vault and may be visible to privileged local processes.

## TUI controls

- Type a message and press Enter to send it to the active channel.
- Tab completes commands, aliases, nicknames and channels when input is present. With an empty input line, Tab cycles through the status view and known channel labels.
- Ctrl+P and Ctrl+N move backward and forward through the in memory command input history.
- `/join #channel` joins a channel.
- `/nick nickname` changes your nickname.
- `/msg nickname message` or `/query nickname` opens a private conversation.
- `/names` refreshes the active channel member list.
- `/topic [text]` reads or changes the active channel topic.
- `/invite nickname` invites a user to the active channel.
- `/mode [flags]` sends a channel MODE command.
- Moderation: `/kick nickname [reason]`, `/ban mask`, `/unban mask`, `/op nickname`, `/deop nickname`, `/voice nickname`, `/devoice nickname`.
- `/randomnick` requests a cryptographically random nickname.
- `/search words` searches the active in memory conversation buffer. Results are not saved to disk.
- `/whois nickname` requests user information from the server.
- `/who nickname-or-mask` requests a WHO lookup.
- `/list [channel]` requests the channel list or a specific channel.
- `/away [message]` sets or clears your away status.
- `/notice target message` sends an IRC NOTICE.
- `/me action` sends a CTCP ACTION to the active channel.
- `/alias name /command $1-` creates a persistent command alias saved under `~/.config/zirc/aliases.json` with private file permissions. `$1` through `$9` insert positional arguments, `$1-` inserts the remaining arguments, `$*` inserts all arguments, and `$` inserts a literal dollar sign.
- `/aliases` lists aliases and `/unalias name` removes one. Alias expansion is limited to five nested levels and cannot execute shell commands or access local files.
- `/kick nickname [reason]`, `/ban mask`, `/unban mask`, `/op nickname`, `/deop nickname`, `/voice nickname` and `/devoice nickname` send standard IRC moderation commands. Server permissions still apply.
- Ctrl+P and Ctrl+N recall previous and next entered commands. Tab completes command names, aliases, nicknames and channels when the input line is not empty. With an empty input line, Tab switches views.
- `/clear` clears the active buffer from the current session memory.
- `/privacy` displays the current identity and route summary.
- `/part` leaves the active channel.
- `/help` lists commands.
- Ctrl+C or `/quit` exits.

## Not yet implemented

- Multiple simultaneous IRC networks
- Complete mIRC style scripting language, remote event handlers, popup menus, identifiers, variables, timers and script editor
- Script loading, unloading, debugging and compatibility import
- Custom dialogs, custom windows, toolbar and popup customization
- Desktop notifications and sound event configuration
- DCC chat and file transfer
- CTCP utility suite
- Addon and plugin manager
- Dedicated interactive privacy settings screen
- Independent Tor route verification and automated DNS or connection leak tests
- Optional encrypted persistent history
- Native macOS and Windows keyring integration
- More extensive connection diagnostics and user controlled retry settings
- Sustained fuzz campaigns and dedicated SOCKS5 handshake fuzzing
- Verified tagged release execution
- Independent security audit

## Security

Read [SECURITY.md](SECURITY.md), [PRIVACY.md](PRIVACY.md), and [THREAT_MODEL.md](THREAT_MODEL.md). Do not use this prototype for high risk communications until a qualified independent security review has been completed.

## License

MIT.

---

## Project Operating System

For public contributor guidance, local verification commands, and AI assistant rules, see [AGENTS.md](./AGENTS.md) and the repository documentation.  
Internal maintainer governance and portfolio tracking use the owner's private `zakshub/project-os` repository and private `Zakshub Portfolio` board (owner access only).
