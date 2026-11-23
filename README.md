# ZIRC Client

Privacy first, terminal native IRC client for Linux.

ZIRC is an independent IRC client. It connects to existing IRC networks. It does not operate an IRC network or server.

## Privacy model

ZIRC is designed to minimise identity disclosure and make network routing visible and explicit.

Important: TLS encrypts traffic but does not hide your IP address from the IRC server. IP masking requires a network supplied cloak, Tor, or a correctly configured proxy. ZIRC must never claim an IP is hidden unless the selected route is actually active.

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
