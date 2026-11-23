# ZIRC Privacy Design

## Data minimisation

ZIRC should collect and transmit only protocol fields required to connect and communicate.

Default identity values:
- Nickname: user supplied or generated locally
- Username: defaults to the selected nickname unless configured
- Real name: generic value, never a legal name by default
- Email: not requested

## Network routing

Implemented routing modes:
- Direct TCP
- TLS direct
- SOCKS5 proxy with proxy side destination hostname resolution
- SOCKS5 username/password authentication

A Tor SOCKS5 endpoint can be configured by the user. ZIRC does not independently verify that a SOCKS5 endpoint is Tor. A network host cloak can reduce exposure to other IRC users, but it does not hide the source IP from IRC network operators.

The connection status distinguishes encryption from routing. A TLS connection is not an anonymised connection. If the selected proxy is unavailable, ZIRC fails rather than falling back to a direct route. Automatic reconnect attempts reuse the same route.

## Credentials and local data

- No chat history is persisted to disk.
- Linux Secret Service keyring lookup is supported through secret-tool.
- CLI commands can remove stored proxy and SASL secrets.
- Environment variables remain supported as an alternative credential source and may be visible to privileged local processes.
- No credentials are written to a ZIRC config file.
- No telemetry or crash uploads are implemented.
- No message logging is implemented.
- Native macOS and Windows keyring integration is not implemented.

## Protocol attack surface

- DCC is not implemented.
- No automatic URL fetching, previews, scripts or remote content requests.
- IRC messages and server numerics are treated as untrusted input.
- IRC command output is bounded to protocol line limits.
- IRC line parser and formatter fuzz harnesses run short campaigns in CI.
- Dedicated SOCKS5 fuzzing, longer fuzz campaigns and independent review remain outstanding.
