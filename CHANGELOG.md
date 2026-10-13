# Changelog

## Unreleased

### Added

- Module metadata using the historical `github.com/luisiturrios/gowin` import path.
- Typed `GetDwordReg` and error-returning `ShellFolders.Lookup` APIs.
- Explicit `ErrUnsupportedPlatform` responses outside Windows; portable environment helpers.
- Unit tests, registry integration tests, handle-leak regression, fuzzing and benchmarks.
- CI for Go 1.26/1.27 on Windows/Linux, race detection, vet and vulnerability scanning.
- Dependabot for Go modules and GitHub Actions; contribution and migration guidance.

### Fixed

- Registry handles leaked by string and DWORD writes.
- Fixed-size registry reads and incorrect byte/UTF-16 size handling.
- Panics from embedded NUL input; invalid writes are rejected before key creation.
- Excessive requested registry permissions and missing error context.
- README descriptions of unimplemented services and process-only environment writes.

### Dependency and compatibility notes

- Replaced `github.com/btcsuite/winsvc/winapi` with `golang.org/x/sys/windows/registry` v0.49.0.
- The approved minimum Go version is 1.26.0, required by that dependency.
- Existing exported names, signatures, constants and import path are preserved.
- Legacy registry conversions and unexpanded shell folder values remain unchanged.
- Error text changes; wrapped causes remain accessible through `errors.Is`/`errors.As`.
- Continuously growing values stop after eight reads with `ErrShortBuffer` wrapped.
- No version tag, commit, push or release has been created for these changes.
