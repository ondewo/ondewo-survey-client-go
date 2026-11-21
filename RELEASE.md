# Release History

*****************

## Release ONDEWO SURVEY Go Client 2.1.1

### New Features

* [[OND236-120]](https://ondewo.atlassian.net/browse/OND236-120) Keycloak authentication: the new `auth.KeycloakTokenProvider` and `client.Config`'s Keycloak fields (`RefreshToken` last), validated before any network I/O and redacted in every rendering; `NewChannel` attaches the shared provider of the credential set to every call.
* [[OND236-120]](https://ondewo.atlassian.net/browse/OND236-120) Offline-token start-up: with `RefreshToken` set the provider starts with one `grant_type=refresh_token` exchange and never sends a password grant, not even as a fallback. Providers are shared per credential set by a SHA-256 of the credentials.
* [[OND236-120]](https://ondewo.atlassian.net/browse/OND236-120) Tests: the cross-SDK Keycloak test catalogue (T01-T62) with a fake token endpoint and clock.

### Bug Fixes

* [[OND236-120]](https://ondewo.atlassian.net/browse/OND236-120) Keycloak token refresh: a failed refresh no longer ends the refresh. Every failure kind (a 5xx from a proxy, a DNS or connection failure, a restarting Keycloak, a dead offline session) is retried with full-jitter exponential backoff, 5 s doubling up to 300 s and measured from the moment the failure was recorded; the lazy `PerRPCCredentials` path and the background refresher share one lock and that backoff. A success resets the backoff.
* [[OND236-120]](https://ondewo.atlassian.net/browse/OND236-120) A token that is still valid keeps being served when its refresh fails inside the 30 s expiry leeway; only an EXPIRED token raises, and it raises a typed Keycloak authentication error whatever failed underneath (transport error, non-2xx status, a body that is not JSON or not a JSON object, a missing `access_token`, a non-numeric `expires_in`), with the original error kept as its cause. A malformed response changes no provider state. Concurrent callers share one in-flight refresh, so an outage costs one token request, not one per call.
* [[OND236-120]](https://ondewo.atlassian.net/browse/OND236-120) A token response without an `access_token` is reported by its field names only, never by its values (a token response is credential material), and a quoted Keycloak error body is cut at 500 characters.
* [[OND236-120]](https://ondewo.atlassian.net/browse/OND236-120) An error is named by its bare type in logs and errors, as the python SDK does.

### Improvements

* [[OND236-120]](https://ondewo.atlassian.net/browse/OND236-120) The same Keycloak behaviour as every ONDEWO SDK (the python SDKs are the reference, ondewo-survey-client-python 2.1.1 ships the same fixes): the token request is bounded by a 10 s timeout, durations use a monotonic clock, the `tokenExpirationInS` window is armed on both start-up paths, each failed refresh logs one WARNING with the attempt number, the retry delay and how long the current token stays valid, and the first success afterwards one INFO. No credential appears in any log line, error or rendering of the configuration.
* Tracking API Version [2.1.0](https://github.com/ondewo/ondewo-survey-api/releases/tag/2.1.0) ( [Documentation](https://ondewo.github.io/ondewo-survey-api/) ), generated with ondewo-proto-compiler 5.17.0: the generated code is unchanged from 2.1.0

*****************

## Release ONDEWO SURVEY Go Client 2.1.0

### Improvements

* Tracking API Version [2.1.0](https://github.com/ondewo/ondewo-survey-api/releases/tag/2.1.0) ( [Documentation](https://ondewo.github.io/ondewo-survey-api/) )
* Generated with ondewo-proto-compiler 5.17.0 (protoc 35.1, `protoc-gen-go` v1.36.12). New runtime floors: `google.golang.org/grpc` v1.85.0, `google.golang.org/protobuf` v1.36.12, genproto v0.0.0-20261116203732-fd2c7689a808 and `go 1.26.0`.

*****************

## Release ONDEWO SURVEY Go Client 2.0.1

### New Features

* `client.NewChannel(cfg client.Config, opts ...grpc.DialOption)` opens the gRPC connection for
  plaintext, TLS or mutual TLS, following the TLS contract of the python SDKs
  (`ondewo-client-utils` 4.1.1), so one set of certificates works with every ONDEWO client.
  `client.Config` carries `Host`, `Port`, `Insecure`, an optional `Logger` and the PEM **content**
  (never a file path) of the CA (`GrpcCert`, empty = system roots) and of an optional client
  identity (`GrpcClientCert` / `GrpcClientKey`).
* Refused with an error before gRPC sees them: half a client identity, a certificate and key that
  do not form a pair, a `GrpcCert` that holds no PEM certificate (typically a path), and
  `Insecure: true` combined with a client identity. No error message contains a PEM, a key or the
  whole `Config`; they name the field and `host:port`.
* A plaintext connection logs a warning naming `host:port` through `log/slog` (`Config.Logger`, or
  `slog.Default()`); the package never configures logging. Every `fmt` verb and `slog` rendering of
  a `Config` redacts the client key.
* Bare IPv6 literal hosts are bracketed (`::1` -> `[::1]:<port>`); bracketed hosts and hosts with a
  scheme are used as they are. PEMs with CRLF line endings work.
* Connection defaults of the python SDKs that grpc-go exposes: keepalive pings every 30 s with a
  20 s timeout, only while an RPC is active; 5 s maximum reconnect backoff; 2^31-1 byte maximum
  message size in both directions. Every `grpc.DialOption` passed in is applied after them and wins.
  Documented gaps: grpc-go has no `http2.max_pings_without_data` and only one keepalive timeout, and
  no per-method retry policy is configured (only gRPC's transparent retries apply).

### Improvements

* `tests/tls_test.go` runs real TLS and mutual-TLS handshakes against an in-process server with a
  PKI generated per run; the 100% coverage gate now spans `auth/` and `client/`.
* `tests/release_notes_test.go` pins the `RELEASE.md` slice the GitHub release body is built from:
  the Makefile's perl range, the spelling of every release heading, the `*****` separator closing
  every section, and non-empty notes for the current version.
* README: new section "TLS, mutual TLS and certificates" (modes, connection defaults, a test PKI with
  openssl, TLS security notes, troubleshooting).
* The ONDEWO proto compiler submodule is pinned to 5.15.2 (was 5.15.1).

*****************

## Release ONDEWO SURVEY Go Client 2.0.0

### New Features

* Initial release of the ONDEWO SURVEY (Survey) gRPC client for Go. The module
  ships the stubs generated from the [ONDEWO SURVEY API](https://github.com/ondewo/ondewo-survey-api)
  by version 5.15.1 of the
  [ONDEWO Proto Compiler](https://github.com/ondewo/ondewo-proto-compiler): one `*.pb.go` of
  messages and one `*_grpc.pb.go` of service stubs per `.proto` — 4 files from 2 protos, below
  `api/ondewo/survey/` — compiled against the `google.golang.org/protobuf` and
  `google.golang.org/grpc` runtimes pinned by the compiler image. It exposes 2 services — `ondewo.survey.Surveys` and `ondewo.survey.FHIR`.
* The generated stubs are **committed**. A Go module is resolved straight from its version control
  tree — `go get` clones the tag and compiles what is in it — so a module that generated its code at
  build time would publish nothing.
* A hand-written `auth` package carries the ONDEWO credential: `auth.WithBearerToken(token)` is a
  `grpc.DialOption` that sends the Keycloak access token as `authorization: Bearer <token>` on every
  call, and refuses to attach itself to a plaintext connection.
* `make build` reproduces the whole client from the two submodules — proto compiler image, stub
  generation and `go build` — and `make check_build` asserts that every `.proto` of the API
  produced a stub.

### Improvements

* A real test suite under `tests/`, run by `.github/workflows/ci.yml` on go 1.25 and go 1.27. It
  needs no network and no ONDEWO server: gRPC connections are made over an in-memory `bufconn`
  listener, so a generated client stub, a generated server stub and a real HTTP/2 connection are
  exercised in-process. Messages round-trip on the wire, the two generators' views of every service
  are cross-checked against each other, and all 13 generated unary stubs are called for real.
* `make test_coverage` gates the build on the coverage of the **hand-written** packages
  (`COVERAGE_THRESHOLD`, currently 100%) and fails if it ends up measuring no function at all, so a
  deleted package cannot turn the gate into a green no-op. `make test_coverage_generated` reports
  the generated stubs' figure without gating it.
* `make check_stubs` runs first in CI and fails unless `api/` holds both message and service stubs
  alongside `go.mod`/`go.sum` — no later step is allowed to pass by finding nothing to do.
* The module path carries the `/v2` major-version suffix its release number requires
  (`github.com/ondewo/ondewo-survey-client-go/v2`). Go resolves a module straight from its VCS path,
  and from major version 2 on that path has to carry the major
  ([module reference](https://go.dev/ref/mod#major-version-suffixes)): a `2.0.0` tag on a module
  declared without `/v2` is invisible to `go get`.
* The release targets tag each release twice on the same commit: with the ONDEWO release number
  (`2.0.0`) that the rest of the fleet uses, and with the `v`-prefixed spelling (`v2.0.0`) that is
  the only tag shape the Go module resolver accepts. `make publish_go_module` then warms
  `proxy.golang.org` so the new version is immediately installable with `go get`.

*****************
