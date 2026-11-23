# Security Policy

## Security guarantees and limitations

ZIRC aims to reduce accidental identity disclosure. It cannot guarantee anonymity.

- The IRC server can observe the source IP address unless traffic is routed through Tor or a proxy, or the network provides a cloak.
- TLS protects traffic in transit but does not conceal the client IP from the IRC server.
- A proxy operator can observe connection metadata.
- Tor does not protect against endpoint compromise, account correlation, or identifying information voluntarily sent in chat.
- IRC operators may retain connection and account metadata under their own policies.

## Secure defaults

1. TLS certificate verification is mandatory when TLS is enabled.
2. No insecure TLS bypass in normal configuration.
3. No DCC transfers by default.
4. No automatic CTCP responses beyond protocol necessities.
5. Message logging disabled by default.
6. No telemetry.
7. No real name or email defaults that identify the user.
8. Secrets must never be committed or printed to logs.
9. SOCKS proxy mode must resolve hostnames through the proxy to reduce DNS leakage.
10. The UI must accurately disclose the active transport and routing mode.

## Reporting

Please report suspected vulnerabilities privately to the repository owner through GitHub security advisories.
