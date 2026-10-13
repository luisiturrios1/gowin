# Modernization migration notes

## Approved toolchain requirement

This modernization uses `golang.org/x/sys` v0.49.0, which requires Go 1.26.0.
The old repository had no `go.mod` or documented minimum, so older-toolchain users
are affected by this newly declared requirement. The maintainer approved Go
1.26.0 and the current `x/sys` dependency on 2026-10-12. Use Go 1.26.9 or Go
1.27.2 (or newer supported security patches).

If older compilers must remain supported, choose an explicit minimum with the
maintainer and reevaluate the dependency, implementation and CI before release.
Do not just lower the `go` directive: that cannot override a dependency's minimum.

## Import path

Continue importing `github.com/luisiturrios/gowin`. GitHub redirects this historical
path to `luisiturrios1/gowin`; the module keeps the historical name to avoid forcing
consumers to change imports. Do not replace it with `github.com/luisiturrios1/gowin`
in `require` directives, because that would disagree with the declared module path.

Consumers already using the newer repository spelling without modules should
use the historical module import when adopting this candidate. Supporting both
module identities is not introduced here. Test module resolution when publishing
the first tag; the candidate has only been validated with a local module replacement.

## Existing calls

No exported function or method was removed or had its signature changed.
`ALL`, `USER`, `ShellFolders.Context` and the `Service` placeholder remain.
The six original folder methods still discard errors; adopt `Lookup` where errors
need to be handled. Numeric DWORD reads can use the new `GetDwordReg` instead of
decoding the legacy UTF-16 representation.

`GetReg` still stops at the first NUL, and `GetRegRaw` still strips trailing zero
units. Binary values retain the legacy interpretation, including zero-padding of
an odd final byte. `REG_EXPAND_SZ` and folder values remain unexpanded.

Registry errors now add context with `%w`. Replace direct equality checks and
error-string comparisons with `errors.Is`; `errors.As` can recover Windows errno.
Inputs containing NUL now return `syscall.EINVAL` wrapped instead of panicking.
Repeated growth during a read eventually returns `registry.ErrShortBuffer` wrapped.

## Version recommendation

There are no existing tags in the audited repository. Start with **v0.1.0** after
passing native Windows CI. Keep this unreleased
until those gates pass. Avoid a v2 module path: the preserved public API does not
justify a major-version migration. If the maintainer establishes a pre-existing
compatibility contract requiring old toolchains or the alternate import spelling,
resolve that contract explicitly before choosing a release version.
