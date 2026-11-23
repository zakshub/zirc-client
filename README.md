# ZIRC Client

Privacy first, terminal native IRC client for Linux.

ZIRC is an independent IRC client. It connects to existing IRC networks. It does not operate an IRC network or server.

## Current prototype

- Full screen keyboard driven Linux TUI
- IRC registration, channel messages, private conversations and per-channel buffers
- Channel member list, topics, NAMES refresh, invites and MODE commands
- TLS with certificate verification enabled by default
- SOCKS5 proxy routing with proxy side destination DNS and optional username/password authentication
- IRCv3 capability negotiation and SASL PLAIN authentication over TLS
- Random local nickname generation and configurable username/real name
- Identity and route preview before connection
- Proxy credentials and SASL passwords read from environment variables, never command line arguments
- No telemetry, message logging, DCC, or automatic link fetching
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

Route through a locally running Tor SOCKS5 service, commonly 127.0.0.1:9050. ZIRC can confirm the SOCKS5 route, but cannot independently prove that the proxy is Tor:

```bash
./zirc --server irc.libera.chat:6697 --tls --proxy 127.0.0.1:9050 --nick zirc-user
```

The SOCKS5 implementation sends the IRC server hostname to the proxy for name resolution. It does not locally resolve the IRC server hostname. If the selected proxy is unavailable, ZIRC fails the connection rather than falling back to a direct connection.

SOCKS5 supports no-auth and RFC 1929 username/password modes. Proxy passwords are read from a named environment variable, not command line arguments. SOCKS5 username/password authentication is not encrypted between the client and proxy, so use it only with a trusted local proxy or a separately protected tunnel. A proxy route is not proof that the endpoint is Tor. Do not use a proxy you do not trust.

Identity preview without connecting:

```bash
./zirc --server irc.libera.chat:6697 --random-nick --identity-preview --tls --proxy 127.0.0.1:9050
```

Use SOCKS5 authentication without exposing the password in shell history or process arguments:

```bash
export ZIRC_PROXY_PASSWORD='your-proxy-password'
./zirc --server irc.example.net:6697 --proxy 127.0.0.1:9050 --proxy-user account --proxy-password-env ZIRC_PROXY_PASSWORD
```

SASL PLAIN authentication is available on TLS connections. Store the password in an environment variable and pass only its variable name:

```bash
export ZIRC_SASL_PASSWORD='your-account-password'
./zirc --server irc.example.net:6697 --tls --sasl-user account --sasl-password-env ZIRC_SASL_PASSWORD
```

Environment variables are not a substitute for an OS keyring and may be visible to privileged local processes. Keyring support remains planned.

## TUI controls

- Type a message and press Enter to send it to the active channel.
- Tab cycles through the status view and known channel labels. Each channel and private conversation has its own transcript buffer.
- `/join #channel` joins a channel.
- `/nick nickname` changes your nickname.
- `/msg nickname message` or `/query nickname` opens a private conversation.
- `/names` refreshes the active channel member list.
- `/topic [text]` reads or changes the active channel topic.
- `/invite nickname` invites a user to the active channel.
- `/mode [flags]` sends a channel MODE command.
- `/randomnick` requests a cryptographically random nickname.
- `/privacy` displays the current identity and route summary.
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
