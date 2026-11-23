# ZIRC Client

Privacy first, terminal native IRC client for Linux.

ZIRC is an independent IRC client. It connects to existing IRC networks. It does not operate an IRC network or server.

## Current prototype

- Full screen keyboard driven Linux TUI
- IRC registration, channels and basic messages
- TLS with certificate verification enabled by default
- SOCKS5 proxy routing with proxy side destination DNS
- No telemetry, message logging, DCC, or automatic link fetching
- Generic real name and username defaults
- Direct versus proxy route status shown in the TUI

This is an early prototype, not yet a security audited or feature complete IRC client.

## Privacy model

TLS encrypts traffic but does not hide your IP address from the IRC server. To hide your public IP from the IRC server, route through Tor or a trusted proxy. An IRC host cloak generally hides your address from other users, not from network operators. ZIRC must never claim an IP is hidden unless the selected proxy route is active.

A proxy operator can observe your source IP and destination. Tor changes the trust model but does not protect against account correlation, endpoint compromise, or identifying information you disclose in chat.

## Build

Requires Go 1.23 or newer.

```bash
go test ./...
go vet ./...
go build -o zirc ./cmd/zirc
```

## Connect

Direct TLS connection. The IRC network can still see your public IP:

```bash
./zirc --server irc.libera.chat:6697 --tls --nick zirc-user
```

Route through a locally running Tor SOCKS5 service, commonly 127.0.0.1:9050:

```bash
./zirc --server irc.libera.chat:6697 --tls --proxy 127.0.0.1:9050 --nick zirc-user
```

The SOCKS5 implementation sends the IRC server hostname to the proxy for name resolution. It does not locally resolve the IRC server hostname. If the selected proxy is unavailable, ZIRC fails the connection rather than falling back to a direct connection.

Current proxy support is unauthenticated SOCKS5 only. Do not use a proxy you do not trust.

## TUI controls

- Type a message and press Enter to send it to the active channel.
- Tab cycles through the status view, channels and private conversations.
- `/join #channel` joins a channel.
- `/nick nickname` changes your nickname.
- `/msg nickname message` sends a private message.
- `/part` leaves the active channel.
- `/help` lists commands.
- Ctrl+C or `/quit` exits.

## Implemented privacy protections

- TLS certificate validation
- SOCKS5 proxy route with remote hostname resolution
- No direct connection fallback when proxy connection fails
- Generic real name and username defaults
- No telemetry
- No message logging
- No DCC implementation
- No automatic link previews or remote content requests
- Explicit direct route warning

## Planned privacy work

- Dedicated privacy settings screen and identity preview
- Random nickname generation
- SOCKS5 authentication
- Tor route verification and diagnostics
- IRCv3 and SASL authentication
- OS keyring integration
- Optional encrypted local history
- Fuzzing and independent security review

## Security

Read [SECURITY.md](SECURITY.md), [PRIVACY.md](PRIVACY.md), and [THREAT_MODEL.md](THREAT_MODEL.md). This prototype has not undergone an independent security audit. Do not use it for high risk communications until the remaining security work and review are complete.

## License

MIT.
