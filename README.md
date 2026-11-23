# ZIRC Client

Privacy conscious, terminal native IRC client for Linux.

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
- Proxy and SASL passwords accepted through environment variables, not command line arguments
- No telemetry, message logging, DCC, or automatic link fetching
- Direct versus configured SOCKS5 route status shown in the TUI
- IRC parser fuzz seed tests included in CI

This remains an early prototype. It is not security audited and is not feature complete.

## Privacy model

TLS encrypts traffic but does not hide your IP address from the IRC server. To hide your public IP from the IRC server, route through Tor or a trusted proxy. An IRC host cloak generally hides your address from other users, not from network operators.

ZIRC can confirm that a SOCKS5 proxy connection was used. It cannot independently prove that the proxy is Tor. A proxy operator can observe your source IP and destination. Tor does not protect against account correlation, endpoint compromise, global traffic correlation, or identifying information you disclose in chat.

When a proxy is selected, ZIRC fails closed if the proxy connection fails. It does not fall back to a direct connection.

## Build

Requires Go 1.23 or newer.

```bash
go test -race ./...
go vet ./...
go build -trimpath -o zirc ./cmd/zirc
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

This confirms use of a SOCKS5 endpoint only. It does not verify that the endpoint is Tor. The SOCKS5 implementation sends the IRC server hostname to the proxy for name resolution and does not locally resolve that hostname.

SOCKS5 supports no authentication and RFC 1929 username/password authentication. The latter is not encrypted between the client and proxy. Use it only with a trusted local proxy or a separately protected tunnel.

Identity preview without connecting:

```bash
./zirc --server irc.libera.chat:6697 --random-nick --identity-preview --tls --proxy 127.0.0.1:9050
```

Proxy authentication:

```bash
export ZIRC_PROXY_PASSWORD='your-proxy-password'
./zirc --server irc.example.net:6697 --proxy 127.0.0.1:9050 --proxy-user account --proxy-password-env ZIRC_PROXY_PASSWORD
```

SASL PLAIN authentication:

```bash
export ZIRC_SASL_PASSWORD='your-account-password'
./zirc --server irc.example.net:6697 --tls --sasl-user account --sasl-password-env ZIRC_SASL_PASSWORD
```

Environment variables are not a substitute for an OS keyring. They may be visible to privileged local processes and can persist in a shell session. Avoid placing literal secrets in shell history.

## TUI controls

- Type a message and press Enter to send it to the active channel.
- Tab cycles through the status view and known channel labels.
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

## Not yet implemented

- OS keyring integration and secret lifecycle controls
- Optional encrypted local history and searchable history
- Multiple simultaneous IRC networks
- Automatic reconnect and connection diagnostics
- Desktop notifications
- Full Unicode terminal display cell width handling
- Dedicated interactive privacy settings screen
- Independent Tor route verification
- Expanded automated DNS and connection leak testing
- Broad protocol fuzz campaigns and independent security review
- Verified tagged release execution

## Security

Read [SECURITY.md](SECURITY.md), [PRIVACY.md](PRIVACY.md), and [THREAT_MODEL.md](THREAT_MODEL.md). Do not use this prototype for high risk communications until a qualified independent security review has been completed.

## License

MIT.
