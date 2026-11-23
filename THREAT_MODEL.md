# ZIRC Threat Model

## Assets

- Public IP address and network metadata
- IRC account credentials and SASL secrets
- Nickname and account identity
- Private messages and channel conversations
- Local process memory and terminal display

## Adversaries

- Passive network observers
- Malicious or compromised IRC servers
- Malicious channel users sending crafted protocol messages
- Malicious proxy operators
- Local users with access to the same machine
- Malicious plugins or scripts, if plugin support is added later

## Mitigations and residual risk

| Threat | Mitigation | Residual risk |
| --- | --- | --- |
| Network eavesdropping | TLS with certificate validation | Traffic metadata remains observable |
| IRC server sees origin IP | User configured SOCKS5 proxy or Tor endpoint | ZIRC does not verify Tor; proxy can see source IP. Automatic reconnect reuses the same configured route and does not fall back to direct |
| DNS leakage | Destination hostname is sent to SOCKS5 proxy | No independent packet capture leak test yet |
| Credential disclosure | Secrets are not command line arguments; SASL PLAIN requires TLS | Environment variables are not a keyring; SOCKS5 auth leg is not encrypted |
| Local chat exposure | No persistent chat history | Messages remain visible in terminal and process memory |
| Malicious protocol input | Bounded line scanner and parser tests | Broad fuzz campaigns and security review remain outstanding |
| Accidental identity disclosure | Generic identity defaults and route status | Users can still disclose identifying information |
| Dangerous file transfer | DCC is not implemented | Future additions could expand attack surface |
| Plugin compromise | No plugin execution | Not applicable to current prototype |

## Non-goals

ZIRC cannot guarantee anonymity against global traffic correlation, endpoint malware, a compromised operating system, account correlation, a malicious proxy, or identifying information users disclose themselves.
