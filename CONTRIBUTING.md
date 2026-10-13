# Contributing

Use small, focused changes that preserve Gowin's simple API. Explain the bug or
use case, compatibility impact and actual validation performed. Public API
removals, changed import paths and behavior changes need maintainer agreement.
Do not expand the `Service` placeholder into a public API without a design review.

The approved toolchain minimum is Go 1.26.0.
CI checks Go 1.26 and 1.27 on Windows and Linux. Use the latest patch version.
Before submitting changes, run:

```sh
go mod tidy
go fmt ./...
go vet ./...
go test ./...
go test -race ./...
go run golang.org/x/vuln/cmd/govulncheck@v1.1.4 -test ./...
```

For decoding changes, run the fuzz target; for performance changes, record the
machine, Go version, command and repeated benchmark results:

```sh
go test -run='^$' -fuzz=FuzzLegacyUTF16 -fuzztime=30s -parallel=2 .
go test -run='^$' -bench=. -benchmem -count=3 .
```

Registry tests must run on Windows. Keep all mutations inside a test-owned random
HKCU key, register cleanup, and avoid HKLM, elevated rights and existing user
settings. Avoid parallel tests when measuring process handle counts. Never delete
regression tests to obtain a passing build.

On macOS or Linux, compile the Windows suite as an additional check:

```sh
GOOS=windows GOARCH=amd64 go vet ./...
GOOS=windows GOARCH=amd64 go test -c -o /tmp/gowin.test.exe .
GOOS=windows GOARCH=amd64 govulncheck -test ./...
```

Install `govulncheck` for the host OS before setting `GOOS`; otherwise `go run`
would build a Windows executable that the host cannot execute. Cross-compilation
is not a substitute for native Windows tests and race detection.

Keep `go.mod` and `go.sum` tidy, add regression tests for fixes, and document new
public functions and compatibility changes in the changelog. Do not add external
dependencies for functionality the standard library already provides. The Windows
registry package from `x/sys` supplies maintained OS bindings absent from portable
standard-library APIs.

Update the CI matrix when Go's supported releases change. Review the pinned
`govulncheck` tool version periodically; it is not a runtime module dependency.
GitHub Actions are pinned to commit SHAs; Dependabot proposes their updates.
Contributions remain under the existing LICENSE.
