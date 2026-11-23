# ZIRC Threat Model

## Assets

- Public IP address and network metadata
- IRC account credentials and SASL secrets
- Nickname and account identity
- Private messages and channel conversations
- Local configuration and optional history

## Adversaries

- Passive network observers
- Malicious or compromised IRC servers
- Malicious channel users sending crafted protocol messages
- Malicious proxy operators
- Local users with access to the same machine
- Malicious plugins or scripts, if plugin support is added later

## Mitigations

| Threat | Mitigation |
| --- | --- |
| Network eavesdropping | TLS with certificate validation |
| IRC server sees origin IP | Tor or SOCKS5 routing. A server cloak only hides the host from other users, not network operators. |
| DNS leakage | Resolve destination through proxy |
| Credential disclosure | No plaintext secrets; OS keyring integration planned |
| Local chat exposure | Logging off by default; encrypted history only if enabled |
| Malicious protocol input | Bounded parser, strict validation |
| Accidental identity disclosure | Generic identity defaults and privacy review |
| Dangerous file transfer | DCC disabled by default |
| Plugin compromise | No plugin execution in initial releases |

## Non-goals

ZIRC cannot guarantee anonymity against global traffic correlation, endpoint malware, a compromised operating system, account correlation, or identifying information users disclose themselves.
