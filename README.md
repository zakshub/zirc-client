# ZIRC Client

Privacy first, terminal native IRC client for Linux.

ZIRC is an independent IRC client. It connects to existing IRC networks. It does not operate an IRC network or server.

## Privacy model

ZIRC is designed to minimise identity disclosure and make network routing visible and explicit.

Important: TLS encrypts traffic but does not hide your IP address from the IRC server. To hide your public IP from the IRC server, route through Tor or a trusted proxy. An IRC host cloak generally hides your address from other users, not from network operators. ZIRC must never claim an IP is hidden unless the selected proxy route is active.

## Planned protections

- TLS with certificate verification enabled by default
- SOCKS5 and Tor routing with remote DNS resolution
- Explicit connection route indicator
- No telemetry or analytics
- No real name or email collection by default
- No DCC file transfers by default
- CTCP replies configurable and restricted by default
- Message logging disabled by default
- Secrets excluded from plaintext configuration
- Redaction of sensitive data in diagnostic output
- No automatic link previews or remote image fetching
- Clear warnings when connecting without TLS or without an anonymising route

## Connect

Direct TLS connection (encrypted, but the IRC server can still see your public IP):

```bash
./zirc --server irc.libera.chat:6697 --tls --nick zirc-user
```

Route through a locally running Tor SOCKS5 service (commonly 127.0.0.1:9050):

```bash
./zirc --server irc.libera.chat:6697 --tls --proxy 127.0.0.1:9050 --nick zirc-user
```

The SOCKS5 implementation sends the IRC server hostname to the proxy for name resolution. This avoids local DNS lookup of the IRC destination. The proxy itself still knows the destination, and Tor or proxy routing does not protect against account correlation or information you disclose in chat.

Current proxy support is unauthenticated SOCKS5 only. Do not use a proxy you do not trust.

## Build

Requires Go 1.23 or newer.

```bash
go test ./...
go build -o zirc ./cmd/zirc
```

## Current scope

The repository is in its initial implementation stage. The security model and development rules are established before advanced networking features are enabled.

## Security reporting

Do not publish exploitable vulnerabilities in public issues. Report them privately to the repository owner.

## License

MIT.
