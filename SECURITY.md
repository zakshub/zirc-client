# Security Policy

## Security guarantees and limitations

ZIRC aims to reduce accidental identity disclosure. It cannot guarantee anonymity or provide protection against a compromised endpoint.

- The IRC server can observe the source IP address unless traffic is routed through Tor or a proxy, or the network provides a cloak.
- TLS protects traffic in transit but does not conceal the client IP from the IRC server.
- A proxy operator can observe connection metadata and the client's source IP.
- A SOCKS5 route does not prove that the endpoint is Tor.
- SOCKS5 username/password authentication is not encrypted between the client and proxy.
- Tor does not protect against endpoint compromise, account correlation, global traffic correlation, or identifying information voluntarily sent in chat.
- IRC operators may retain connection and account metadata under their own policies.

## Implemented protections

1. TLS certificate verification is enabled when TLS is selected.
2. TLS uses a minimum version of TLS 1.2.
3. There is no insecure TLS certificate bypass option.
4. DCC transfers are not implemented.
5. Message logging and telemetry are not implemented.
6. Generic identity defaults avoid requesting a legal name or email address.
7. SOCKS5 proxy mode sends the destination hostname to the proxy for name resolution.
8. A selected proxy failure does not trigger direct connection fallback.
9. SASL PLAIN is rejected unless TLS is enabled.
10. Proxy and SASL passwords are read from named environment variables, not command line arguments.
11. The interface distinguishes a direct route from a configured SOCKS5 route and does not label an unverified proxy as Tor.

## Known limitations

- Secrets are not stored in an OS keyring. Environment variables are not a secure secret vault.
- SOCKS5 username/password authentication does not encrypt credentials on the client to proxy leg.
- The application has not undergone an independent security audit.
- Automated leak testing and broad fuzzing remain incomplete.
- No guarantee is made that all terminal emulators render untrusted Unicode safely or consistently.

## Reporting

Please report suspected vulnerabilities privately to the repository owner through GitHub security advisories.
