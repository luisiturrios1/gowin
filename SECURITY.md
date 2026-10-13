# Security

Gowin executes registry operations with the calling process's Windows privileges.
It does not elevate privileges, run commands, make network requests, log values
or store credentials. Applications must decide which registry keys callers may
access; Gowin is not an authorization boundary. Do not pass untrusted registry
paths into privileged application code without application-level authorization.

Registry input containing NUL is rejected before writes. Each opened key is local
to one call and closed when the call returns. Reads cope with changing sizes and
bound retries; they still allocate memory proportional to the stored value size.
Callers should avoid exposing arbitrary registry reads to untrusted users.

Environment writes affect the process and child processes. Applications must
handle secrets appropriately; Gowin does not provide secret storage or encryption.

CI scans Windows code and tests with govulncheck. A clean scan reflects the current
Go vulnerability database and does not guarantee absence of vulnerabilities.
Use current patches of supported Go versions. See CONTRIBUTING.md for commands.

For a suspected security issue, use GitHub's private vulnerability reporting on
this repository if enabled. Otherwise, contact the maintainer privately before
posting sensitive details publicly. No response-time or support-duration guarantee
is established by this document.
