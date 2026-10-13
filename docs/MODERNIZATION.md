# Gowin modernization audit and validation

Audit date: 2026-10-12. Original revision: `2e88a53ef0f9664a83aeb877fa3862d898eb4ca1`.
Checkout: `/Users/luisiturrios/Repos/gowin`; it matched `origin/master` and started
with a clean working tree. No commits, pushes, tags or releases were made.

## Initial scope and public API

The entire initial repository contained four Go files, README.md and LICENSE.
There were no subpackages, generated files, module metadata, tests, benchmarks,
workflows, contribution guide or changelog. No applicable AGENTS.md was found.

- `envvar.go`: `GetEnvVar` and `WriteEnvVar` wrap the process environment.
- `registry.go`: `GetRegRaw`, `GetReg`, `WriteStringReg`, `WriteDwordReg`, `DeleteKey`.
  It uses Windows `syscall` operations and `github.com/btcsuite/winsvc/winapi`.
- `paths.go`: `ShellFolders{Context bool}`, `ALL`, `USER`, and six public folder
  methods (`ProgramFiles`, `AppData`, `Desktop`, `Documents`, `StartMenu`,
  `StartMenuPrograms`) reading specific registry entries.
- `service.go`: exported `Service`, with private fields and five private no-op
  lifecycle methods. It does not expose working service management.
- README used the historical import `github.com/luisiturrios/gowin`. GitHub redirects
  this URL to `github.com/luisiturrios1/gowin`.
- LICENSE contains the existing GNU GPL version 2 text; it remains unchanged.

There is no production goroutine, channel, mutex, subprocess, network client or
persistent global cache to redesign. The API is already small and function-based.

## Prioritized diagnosis before code changes

| Priority | Finding and evidence | Decision |
| --- | --- | --- |
| Critical | No critical exploit or vulnerable dependency was confirmed in the initial inspection. | Do not invent a critical vulnerability from dependency age. |
| High | Both registry writes defer closing an initial handle, overwrite it with the created key's handle, then never close the latter. Deferred arguments are evaluated immediately. | Create the target key directly and close its owned handle. Add a Windows process-handle regression test. |
| High | Reads pass a count of UTF-16 units as a byte size, ignore the actual returned size/type, and scan a fixed MAX_LONG_PATH array. Large values fail; DWORD semantics are undocumented. | Read actual bytes with growth handling, retain the legacy conversion, and add typed DWORD reading. |
| High | `syscall.StringToUTF16Ptr` and `StringToUTF16` panic on embedded NUL input. | Validate before mutation and return errors wrapping EINVAL. |
| High | No module or pinned dependency; default `go test ./...` cannot run. No tests or CI protect Windows behavior. | Add module metadata, integration tests and native Windows CI. |
| Medium | Five functions duplicate the HKEY switch; writes request KEY_ALL_ACCESS. | Centralize root validation and use narrower rights for each operation. |
| Medium | Folder methods discard every error; values can contain unexpanded placeholders. | Add `Lookup` returning an error; preserve legacy methods and expansion behavior. |
| Medium | README claims service management and does not distinguish process-only environment writes. | Correct feature scope; retain Service instead of removing a public type. |
| Low | Empty imports, inconsistent formatting/comments, no contribution/migration/changelog guidance. | Format and document without changing the library's philosophy. |

Age alone was not used to classify code as obsolete. The replacement dependency
provides maintained Windows bindings and handle-owning APIs; no public architecture
rewrite, new service feature or generalized abstraction was introduced.

## Baseline

Before touching project code:

- `go test ./...`: failed because there was no main module.
- `GO111MODULE=off go test ./...`: failed because winsvc/winapi was unavailable
  in GOPATH. The Windows-only source also lacked platform file constraints.
- In an isolated copy at `/tmp/gowin-baseline.5gkQJy`, temporary module metadata
  and winsvc v1.0.0 were added for baseline inspection only.
- `GOOS=windows GOARCH=amd64 go test -c ... .`: passed, reporting no test files.
- `GOOS=windows GOARCH=amd64 go vet ./...`: passed.
- There were no baseline benchmarks or runtime assertions to compare against.

## Official versions and dependency decision

[Go release history](https://go.dev/doc/devel/release) and the
[official download feed](https://go.dev/dl/?mode=json) identify Go 1.27.2 as the
current stable release and Go 1.26.9 as the other supported release line at the
audit date. Go supports each release until two newer major releases exist.

`go list -m -json golang.org/x/sys@latest` returned **v0.49.0**, dated 2026-10-06,
with GoVersion **1.26.0**. `go mod download -json` verified its module hashes and
upstream origin (`https://go.googlesource.com/sys`, tag `refs/tags/v0.49.0`, commit
`b6b8557e8e2974925aec7a136ba0e323eedc5f43`). This replaces winsvc/winapi;
`go list -m` reported winsvc's latest release as v1.0.0 from 2015, but no specific
vulnerability is claimed for it.

The maintainer explicitly approved **Go 1.26.0 and current x/sys** during this
work. This excludes older compilers; the requirement is justified by the maintained
OS dependency, rather than by an unnecessary language feature. Go 1.27 is tested,
but not required. go.sum contains only the two entries for x/sys v0.49.0.

Tooling uses govulncheck v1.1.4 and actionlint v1.7.12; these are not library
runtime dependencies. GitHub Actions use checkout v7.0.1 and setup-go v7.0.0,
verified against upstream release tags and pinned to their commit SHAs.

## Changes and affected files

| Area | Files | Result |
| --- | --- | --- |
| Module | go.mod, go.sum | Reproducible dependency resolution; historical module identity preserved. |
| Registry | registry.go, registry_windows.go, registry_other.go | Shared validation, maintained Windows bindings, safe byte decoding, owned handles, wrapping, eight-attempt growth bound, typed DWORD API, explicit platform error. |
| Folders | paths.go | One mapping implementation and additive `ShellFolders.Lookup`; original six methods retained. |
| Environment and package docs | envvar.go, doc.go | Accurate process semantics and portable wrappers. |
| Services | service.go | Format and accurate placeholder documentation; no public functionality removed or invented. |
| Tests | envvar_test.go, paths_test.go, registry_test.go, registry_other_test.go, registry_windows_test.go | Table tests, concurrency checks, fuzzing, registry integration and handle-leak regression. |
| Example | examples/basic/main.go | Complete, read-only quick start. |
| Automation | .github/workflows/ci.yml, .github/dependabot.yml | Windows/Linux Go matrix, formatting, vet, tests, race, fuzzing, cross-builds, registry benchmark, Windows govulncheck and weekly dependency updates. |
| Documentation | README.md, CONTRIBUTING.md, CHANGELOG.md, MIGRATION.md, SECURITY.md, docs/MODERNIZATION.md | Actual features, compatibility decisions, maintenance and measured validation. |

No interface hierarchy, generics, context parameter or external helper dependency
was added. These local, synchronous Windows calls do not provide a cancellable
operation; merely adding context would not cancel the underlying API. The boolean
`ShellFolders.Context` remains for source compatibility. Standard-library
`errors`, `encoding/binary`, `unicode/utf16`, `slices` (in tests), `testing.B.Loop`
and `sync.WaitGroup.Go` are used where appropriate.

## Security and architecture results

- Every opened/created registry handle has one local owner and a deferred close.
  Predefined hive handles are not closed or shared as mutable state.
- NUL validation occurs before writes create a key. Error messages include
  operation/key/name context and do not include written values.
- Reads use QUERY_VALUE; writes use SET_VALUE instead of ALL_ACCESS. Deletion
  opens the parent with ENUMERATE_SUB_KEYS. The Windows API checks deletion rights
  on the child independently; see [RegDeleteKeyW](https://learn.microsoft.com/en-us/windows/win32/api/winreg/nf-winreg-regdeletekeyw).
- Production Gowin no longer imports unsafe. A Windows-only regression test uses
  one pointer for the GetProcessHandleCount API to detect the original leak.
- Buffer growth retries are bounded, including concurrent modification. Allocations
  remain proportional to stored data; arbitrary privileged registry access must
  still be authorized by the application.
- There are no new production goroutines, mutable caches or cancellation leaks.
  Environment access relies on the standard library's synchronization. Callers
  still must synchronize changes to ShellFolders.Context.
- govulncheck's macOS and Windows source scans (including tests) reported
  **No vulnerabilities found**. This is not an exhaustive security proof.

## Validation performed

Host: macOS darwin/arm64, Apple M4 Pro. Latest toolchain: Go 1.27.2.

| Command/check | Actual result |
| --- | --- |
| go mod tidy; go fmt ./... | Passed; module files reproducible and sources formatted. |
| go vet ./... | Passed with Go 1.26.0, 1.26.9 and 1.27.2. |
| go test ./... | Passed with Go 1.26.0, 1.26.9 and 1.27.2 on macOS. |
| go test -race ./... | Passed with Go 1.26.9 and 1.27.2 on macOS. |
| go test -cover ./... | Library: 87.0% for the macOS build; executable example: 0.0%. This excludes Windows-only files. |
| go test -run='^$' -fuzz=FuzzLegacyUTF16 -fuzztime=15s -parallel=2 . | Passed: 3,027,818 executions, no failures. |
| go test -run='^$' -bench=. -benchmem -count=3 . | Passed; results below. |
| GOOS=windows GOARCH=amd64 go vet ./... | Passed with Go 1.26.9 and 1.27.2. |
| Windows amd64 test-binary compilation | Passed with Go 1.26.0, 1.26.9 and 1.27.2. |
| Windows 386 and arm64 builds and test-binary compilation | Passed with Go 1.27.2. |
| govulncheck -test ./... | Passed for macOS and Windows/amd64 targets using Go 1.27.2. |
| go mod verify | All modules verified. |
| actionlint v1.7.12 | Passed for .github/workflows/ci.yml. |
| README snippets as external consumer | All five examples pass Windows build and vet with a local module replacement. |
| git diff --check | Passed. |

No Windows executable was run locally. Registry integration, native Windows race
checks, handle counts and BenchmarkRegistryRead are compiled and configured in CI,
**but their runtime results remain pending**. GitHub Actions has not been dispatched
for this unpublished patch. actionlint validates workflow structure; it cannot
prove the remote jobs pass. Linux execution is also configured but was not run on
this macOS host.

### Benchmark results

Three runs, Go 1.27.2, darwin/arm64, Apple M4 Pro:

| Operation | ns/op range | B/op | allocs/op |
| --- | ---: | ---: | ---: |
| GetEnvVar | 14.31–14.42 | 0 | 0 |
| Legacy UTF-16, 32 units | 22.05–22.75 | 64 | 1 |
| Legacy UTF-16, 4,096 units | 2,444–2,470 | 8,192 | 1 |
| Legacy UTF-16, 100,000 units | 58,077–58,671 | 204,800 | 1 |

These establish a reproducible reference. There is no measured before/after
performance improvement claim. Native registry performance must be measured on Windows.

## Compatibility and release recommendation

The existing exported names, signatures, constants, struct fields and historical
module identity remain. GetRegRaw still returns UTF-16 units for every registry
type; odd binary data is zero-padded and trailing zero units are removed. GetReg
still stops at the first NUL. No folder expansion, Known Folder migration, public
service operation or recursive deletion was introduced.

Compatibility changes are documented: approved Go 1.26 minimum; contextual error
text and wrapping; NUL input returns an error rather than a panic; large values
now read successfully; rapidly growing values eventually return a wrapped short-
buffer error. Error equality/string matching should migrate to errors.Is/As.
Adding a go.mod also means consumers using the alternate `luisiturrios1` spelling
need to adopt the historical declared module identity. No forced canonical-path
rename was made.

There were no existing release tags. Recommend **v0.1.0**, after native Windows CI
passes and the maintainer reviews the patch. A v2 import path is unnecessary for
this compatible public API. The code is prepared for review; production readiness
cannot be certified without the native Windows runtime validation.

## Pending work and roadmap

1. Run the committed candidate in native Windows CI before release; review test,
   race, handle-regression, vulnerability and registry benchmark results. No push
   or commit was authorized for this session.
2. Verify resolution of the preserved historical module path from a clean external
   consumer once the first tag exists, without a local replacement.
3. Consider optional typed registry reads (binary, QWORD, multi-string) and explicit
   registry view selection only with concrete consumer use cases.
4. Design an additive Known Folder API if callers need fully resolved paths;
   preserve legacy getters and use Windows runtime tests for any behavior change.
5. Decide whether service management belongs in scope; review lifecycle, privilege,
   timeouts and cancellation before exposing any new Service API.
6. Keep supported Go lines, SHA-pinned actions, govulncheck and dependencies current.
   Do not add performance claims without repeatable before/after measurements.
