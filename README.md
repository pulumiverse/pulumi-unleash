# Unleash Pulumi Provider

A [Pulumi](https://www.pulumi.com) provider for [Unleash](https://www.getunleash.io), bridged from
the official [`Unleash/terraform-provider-unleash`](https://github.com/Unleash/terraform-provider-unleash)
using the [Pulumi Terraform Bridge](https://github.com/pulumi/pulumi-terraform-bridge).

Because it bridges the upstream Terraform provider directly, resources and data sources come from the
same schema rather than a parallel reimplementation, so the two stay in lockstep as the upstream
provider evolves.

## Installation

```bash
npm install @pulumiverse/unleash
pip install pulumiverse-unleash
dotnet add package Pulumiverse.Unleash
go get github.com/pulumiverse/pulumi-unleash/sdk/v3
```

The plugin binary is resolved from this repository's GitHub releases via the `PluginDownloadURL`
baked into the schema, so `pulumi up` fetches it with no extra configuration.

The provider is not in the Pulumi Registry, so a *manual* `pulumi plugin install` needs the server
spelled out. Installing the SDK is the normal path and needs none of this:

```bash
pulumi plugin install resource unleash 3.0.0 --server github://api.github.com/pulumiverse
```

## Building from source

Prerequisites: Go 1.25+, [`pulumictl`](https://github.com/pulumi/pulumictl),
[`mise`](https://mise.jdx.dev), the Pulumi CLI, Node.js, Yarn, Python 3, and the .NET SDK.

```bash
pulumi plugin install converter terraform   # required - schema generation uses pulumiConvert
make tfgen        # generate the Pulumi schema from the upstream Terraform provider
make provider     # build bin/pulumi-resource-unleash
make build_sdks   # generate and build the dotnet, go, nodejs and python SDKs
```

## Configuration

| Option | Environment variable | Notes |
| --- | --- | --- |
| `baseUrl` | `UNLEASH_URL` | Everything before `/api` |
| `authorization` | `UNLEASH_AUTH_TOKEN` | Admin API token; marked secret |
| `maxConcurrentRequests` | `UNLEASH_MAX_CONCURRENT_REQUESTS` | Upstream default is 2 |

The upstream Terraform provider also accepts a bare `AUTH_TOKEN` as a fallback. **This provider
deliberately does not**, because an unscoped `AUTH_TOKEN` sitting in a shell or CI environment would
silently authenticate against Unleash. Use `UNLEASH_AUTH_TOKEN`.

## Resource coverage and the Enterprise gap

All 13 resources and 8 data sources from the upstream provider are bridged and present in the
generated schema.

Automated end-to-end tests cover only what an **open-source** Unleash server can actually do. Unleash
OSS ships a fixed `default` project plus `development`/`production` environments, and its community
build registers no write routes for most administrative APIs — creating projects, environments,
users, groups, roles or service accounts, and configuring OIDC/SAML all require an Enterprise
licence. The upstream provider's own acceptance tests set `UNLEASH_LICENSE` for this reason.

| Resource | Live test coverage |
| --- | --- |
| `ApiToken` | Real create/destroy against OSS Unleash in CI |
| Everything else | Compile and schema generation only |

The `ApiToken` test (`examples/oss-ts`, driven by `TestOssTs` in
`examples/examples_nodejs_test.go`) starts its own Unleash + Postgres containers via
[testcontainers-go](https://github.com/testcontainers/testcontainers-go), so it runs identically locally
and in CI with no setup step. The root [`docker-compose.yml`](./docker-compose.yml) is still there
as a convenience for running the example by hand (`docker compose up` + `pulumi up`) — it is
optional and the test does not depend on it.

Contributions widening this are welcome, but note that a licence secret is not available to pull
requests from forks, so fork CI cannot exercise Enterprise resources either.

## Versioning

The major version tracks the upstream Terraform provider's major version, matching the convention
other bridged community providers use. This provider's 3.x line bridges
`terraform-provider-unleash` 3.x.

## License

Apache-2.0 — see [LICENSE](./LICENSE).
