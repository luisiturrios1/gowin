# Gowin

A small Go library for working with the Windows registry, shell folder locations,
and process environment variables.

[API documentation](https://pkg.go.dev/github.com/luisiturrios/gowin) ·
[Examples](examples/basic) · [Changelog](CHANGELOG.md) ·
[Contributing](CONTRIBUTING.md)

## Contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Quick start](#quick-start)
- [Usage](#usage)
- [API documentation](#api-documentation)
- [Compatibility and limitations](#compatibility-and-limitations)
- [Development](#development)
- [Contributing](#contributing)
- [Security](#security)
- [Changelog](#changelog)
- [License](#license)
- [Support](#support)

## Features

- Read registry values and write strings or DWORDs.
- Create registry keys through writes and delete individual subkeys.
- Read current-user or shared shell folder locations.
- Read and update process environment variables.

Gowin keeps a small, function-based API. Windows service management is not
implemented; `Service` remains a legacy placeholder.

## Requirements

- **Go 1.26.0 or newer.** CI is configured for Go 1.26 and 1.27.
- **Windows** for registry and shell folder operations. Supported Windows versions
  follow [Go's Windows support policy](https://go.dev/wiki/Windows).

Use the latest security patch for your Go release. Environment helpers work on
every platform supported by Go. On other operating systems, registry operations
return `gowin.ErrUnsupportedPlatform`.

## Installation

Inside your Go module, run:

```sh
go get github.com/luisiturrios/gowin@latest
```

Import the library:

```go
import "github.com/luisiturrios/gowin"
```

The historical import path is retained even though the repository is now hosted
at `github.com/luisiturrios1/gowin`.

**The modernization changes in this checkout are unreleased.** `@latest` fetches
the published repository state. To try the local changes in another module, see
[Using a local checkout](#using-a-local-checkout).

## Quick start

Read the current user's Documents registry value and an environment variable:

```go
package main

import (
	"fmt"
	"log"

	"github.com/luisiturrios/gowin"
)

func main() {
	folders := gowin.ShellFolders{Context: gowin.USER}
	path, err := folders.Lookup("Documents")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Documents (registry value):", path)
	fmt.Println("USERPROFILE:", gowin.GetEnvVar("USERPROFILE"))
}
```

From this repository, run the [complete example](examples/basic/main.go) on Windows:

```sh
go run ./examples/basic
```

This example only reads values and does not modify the registry.

## Usage

The following snippets belong inside functions; snippets containing `return err`
require a function that returns `error`.

### Shell folders

Use `gowin.USER` for the current user or `gowin.ALL` for shared folders:

```go
folders := gowin.ShellFolders{Context: gowin.ALL}
path, err := folders.Lookup("AppData")
if err != nil {
	return err
}
fmt.Println(path)
```

`Lookup` accepts `ProgramFiles`, `AppData`, `Desktop`, `Documents`, `StartMenu`,
and `StartMenuPrograms`. `ProgramFiles` ignores the selected context.

The original methods, such as `folders.AppData()` and `folders.Documents()`,
remain available. They discard read errors and return an empty string on failure.
Use `Lookup` when you need to handle errors.

### Windows registry

Create an application key, read back a string and a DWORD, then remove the key:

```go
const path = `Software\MyApplication\GowinExample`
if err := gowin.WriteStringReg("HKCU", path, "Greeting", "Hello, Windows!"); err != nil {
	return err
}
if err := gowin.WriteDwordReg("HKCU", path, "Count", 42); err != nil {
	return err
}
greeting, err := gowin.GetReg("HKCU", path, "Greeting")
if err != nil {
	return err
}
count, err := gowin.GetDwordReg("HKCU", path, "Count")
if err != nil {
	return err
}
fmt.Println(greeting, count)
return gowin.DeleteKey("HKCU", `Software\MyApplication`, "GowinExample")
```

Use a key owned by your application: deleting a key also deletes its values.
`DeleteKey` deletes a subkey, not an individual value, and does not recursively
delete child keys.

Supported hive names are `HKLM`, `HKCU`, `HKCC`, `HKCR`, and `HKU` (case-sensitive).
An empty value name selects the default value. Writes create missing keys. Access
follows Windows ACLs; HKCU normally needs no administrator rights, while protected
keys may require elevation.

Use `GetDwordReg` for numeric DWORD values. `GetReg` interprets values as UTF-16
and stops at the first NUL; it does not format integers as decimal strings.
`GetRegRaw` returns the legacy `[]uint16` representation, retaining embedded NUL
units and stripping trailing zero units.

### Error handling

Windows errors include operation and key context. Use `errors.Is` and `errors.As`
to inspect their underlying causes:

```go
// Requires imports: "errors" and "golang.org/x/sys/windows/registry".
_, err := gowin.GetReg("HKCU", `Software\MyApplication`, "Missing")
if errors.Is(err, registry.ErrNotExist) {
	// The key or value does not exist.
}
```

Registry paths, names, and written strings containing NUL return an error instead
of panicking. Avoid matching error text or comparing wrapped errors directly.

### Environment variables

```go
if err := gowin.WriteEnvVar("MY_APP_MODE", "development"); err != nil {
	return err
}
fmt.Println(gowin.GetEnvVar("MY_APP_MODE"))
```

These helpers wrap `os.Getenv` and `os.Setenv`. Writes affect this process and
subsequently created child processes; they do not persist Windows user or machine
environment variables or update other running processes.

An empty result cannot distinguish an unset variable from an empty one. Use
`os.LookupEnv` when that distinction matters.

## API documentation

See the [published API reference](https://pkg.go.dev/github.com/luisiturrios/gowin)
for indexed versions. For the API in this checkout, run:

```sh
go doc .
```

## Compatibility and limitations

- `REG_EXPAND_SZ` placeholders such as `%USERPROFILE%` remain unexpanded,
  including shell folder values.
- Folder locations use the original registry entries rather than the Known Folder
  API. Values can be missing or depend on the Windows user profile.
- Registry view follows the process architecture. There is no explicit WOW64 view
  selector or transaction API.
- Reads retry growing values up to eight times, then return an error wrapping
  `registry.ErrShortBuffer`. Operations are synchronous and have no cancellation API.
- Read-only folder lookups may run concurrently. Callers must synchronize changes
  to `ShellFolders.Context`.
- `Service` has no public lifecycle operations and does not manage services.

See the [migration guide](MIGRATION.md) for the Go minimum, module import path,
and compatibility changes.

## Development

### Running checks

From the repository root:

```sh
go mod tidy
go fmt ./...
go vet ./...
go test ./...
go test -race ./...
```

Run fuzzing, benchmarks, and vulnerability analysis:

```sh
go test -run='^$' -fuzz=FuzzLegacyUTF16 -fuzztime=15s -parallel=2 .
go test -run='^$' -bench=. -benchmem .
go run golang.org/x/vuln/cmd/govulncheck@v1.1.4 -test ./...
```

Run on Windows to exercise real registry operations. Tests use isolated random
HKCU keys and clean them up; no administrator privileges are required. On other
platforms, tests cover environment helpers, decoding, validation, folder mappings,
and unsupported-platform behavior. Cross-compilation does not validate Windows
runtime behavior.

CI is configured for Windows and Linux tests, formatting, `go vet`, race detection,
fuzzing, a Windows registry benchmark, and Windows vulnerability scanning.
Dependabot checks Go modules and GitHub Actions weekly. See the
[audit and validation report](docs/MODERNIZATION.md) for measured results and
outstanding native Windows validation.

### Using a local checkout

To test unreleased changes from another Go module, replace the dependency with
your local Gowin checkout:

```sh
go mod edit -require=github.com/luisiturrios/gowin@v0.0.0
go mod edit -replace=github.com/luisiturrios/gowin=/absolute/path/to/gowin
go mod tidy
```

## Contributing

Contributions are welcome. Read [CONTRIBUTING.md](CONTRIBUTING.md) for development
checks, regression-test expectations, and public API compatibility guidelines.

## Security

See [SECURITY.md](SECURITY.md) for security considerations and guidance on reporting
suspected vulnerabilities.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for unreleased changes and release notes.

## License

Gowin uses the existing [GNU General Public License version 2](LICENSE).

## Support

To support development, the original donation contact is luisiturrios@me.com.

<!-- Modified 2026-10-12: conventional Go library README structure, usage and maintenance documentation. -->
