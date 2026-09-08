# Changelog

## [0.4.0] - 2026-SEP-08

### Changed

- **Breaking:** `HttpGet`, `HttpPost`, `HttpPut`, `HttpDelete`, and `HttpPatch` take a last `ErrorParserFunc` argument. Pass `nil` to keep v0.3.0 behavior (`*ApiError` unmarshaled from `{"message": "..."}`). Exchange and INTX SDKs must add `nil` at existing call sites.
- `ApiResponse.Error` is now `error` so a custom parser can return a product-specific type.

### Added

- `ErrorParserFunc`: optional hook for unexpected HTTP statuses. Transport failures (invalid URL, `Do`, body read) remain `*ApiError`.
- `headersFunc` is nil-safe.

### Security

- Bump `golang.org/x/net` to v0.55.0 ([CVE-2026-25680](https://github.com/advisories/GHSA-5cv4-jp36-h3mw) / Dependabot alert #1).

## [0.3.0] - 2026-05-29

### Changed

- Relocated module from `github.com/coinbase-samples/core-go` to `github.com/coinbase/core-go`.
- No intentional API changes; equivalent to v0.2.3 on the previous module path.
