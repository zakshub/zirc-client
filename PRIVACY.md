# ZIRC Privacy Design

## Data minimisation

ZIRC should collect and transmit only protocol fields required to connect and communicate.

Default identity values:
- Nickname: user supplied or generated locally
- Username: non-identifying local default
- Real name: generic value, never a legal name by default
- Email: not requested unless a network feature explicitly requires it

## Network routing

Supported routing modes planned:
- Direct TCP
- TLS direct
- SOCKS5 proxy with remote DNS
- Tor SOCKS5 endpoint

A network host cloak can reduce exposure to other IRC users, but it does not hide the source IP from IRC network operators.

The connection status must distinguish encryption from IP masking. A TLS connection is not an anonymised connection.

## Local data

- No chat history unless the user explicitly enables it.
- No credentials in plaintext config files.
- No telemetry or crash uploads.
- Diagnostic logs redact hostnames, nicknames, account names and message contents where possible.
- Configuration file permissions should be restricted to the current user.

## Protocol attack surface

- DCC is disabled by default.
- CTCP handling is conservative and configurable.
- No automatic URL fetching, previews, scripts or remote content.
- IRC messages and server numerics are treated as untrusted input.
- Input and output are bounded to prevent memory exhaustion.
