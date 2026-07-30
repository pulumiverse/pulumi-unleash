// Copyright 2016-2024, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package unleash

import (
	"path"

	// Allow embedding bridge-metadata.json in the provider.
	_ "embed"

	unleashshim "github.com/Unleash/terraform-provider-unleash/shim"

	pfbridge "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfbridge"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge/tokens"

	"github.com/hagaym1/pulumi-unleash/provider/v3/pkg/version"
)

// all of the token components used below.
const (
	// This variable controls the default name of the package in the package
	// registries for nodejs and python:
	mainPkg = "unleash"
	// modules:
	mainMod = "index" // the unleash module
)

//go:embed cmd/pulumi-resource-unleash/bridge-metadata.json
var metadata []byte

// Provider returns additional overlaid schema and metadata associated with the provider.
func Provider() tfbridge.ProviderInfo {
	// Create a Pulumi provider mapping
	prov := tfbridge.ProviderInfo{
		P: pfbridge.ShimProvider(unleashshim.New(version.Version)()),

		Name:    "unleash",
		Version: version.Version,
		// DisplayName is a way to be able to change the casing of the provider name when being
		// displayed on the Pulumi registry
		DisplayName: "Unleash",
		// Change this to your personal name (or a company name) that you would like to be shown in
		// the Pulumi Registry if this package is published there.
		Publisher: "Pulumiverse",
		// LogoURL is optional but useful to help identify your package in the Pulumi Registry
		// if this package is published there.
		//
		// You may host a logo on a domain you control or add an PNG logo (100x100) for your package
		// in your repository and use the raw content URL for that file as your logo URL.
		LogoURL: "",
		// PluginDownloadURL is an optional URL used to download the Provider
		// for use in Pulumi programs
		// e.g. https://github.com/org/pulumi-provider-name/releases/download/v${VERSION}/
		PluginDownloadURL: "",
		Description:       "A Pulumi package for creating and managing unleash cloud resources.",
		// category/cloud tag helps with categorizing the package in the Pulumi Registry.
		// For all available categories, see `Keywords` in
		// https://www.pulumi.com/docs/guides/pulumi-packages/schema/#package.
		Keywords:   []string{"unleash", "category/cloud"},
		License:    "Apache-2.0",
		Homepage:   "https://www.pulumi.com",
		Repository: "https://github.com/pulumiverse/pulumi-unleash",
		// The GitHub Org for the provider - defaults to `terraform-providers`. Note that this should
		// match the TF provider module's require directive, not any replace directives.
		GitHubOrg:    "Unleash",
		MetadataInfo: tfbridge.NewProviderMetadata(metadata),
		Config: map[string]*tfbridge.SchemaInfo{
			"base_url": {
				Default: &tfbridge.DefaultInfo{EnvVars: []string{"UNLEASH_URL"}},
			},
			"authorization": {
				Secret:  tfbridge.True(),
				Default: &tfbridge.DefaultInfo{EnvVars: []string{"UNLEASH_AUTH_TOKEN"}},
			},
			"max_concurrent_requests": {
				Default: &tfbridge.DefaultInfo{EnvVars: []string{"UNLEASH_MAX_CONCURRENT_REQUESTS"}},
			},
		},
		JavaScript: &tfbridge.JavaScriptInfo{
			// Publish under the Pulumiverse npm scope, not the official @pulumi namespace.
			PackageName: "@pulumiverse/unleash",
			// RespectSchemaVersion ensures the SDK is generated linking to the correct version of the provider.
			RespectSchemaVersion: true,
		},
		Python: &tfbridge.PythonInfo{
			// Publish under the Pulumiverse PyPI namespace, not the official pulumi_ prefix.
			PackageName: "pulumiverse_unleash",
			// RespectSchemaVersion ensures the SDK is generated linking to the correct version of the provider.
			RespectSchemaVersion: true,
			// Enable modern PyProject support in the generated Python SDK.
			PyProject: struct{ Enabled bool }{true},
		},
		Golang: &tfbridge.GolangInfo{
			// Set where the SDK is going to be published to.
			ImportBasePath: path.Join(
				"github.com/hagaym1/pulumi-unleash/sdk/",
				tfbridge.GetModuleMajorVersion(version.Version),
				"go",
				mainPkg,
			),
			// Opt in to all available code generation features.
			GenerateResourceContainerTypes: true,
			GenerateExtraInputTypes:        true,
			// RespectSchemaVersion ensures the SDK is generated linking to the correct version of the provider.
			RespectSchemaVersion: true,
		},
		Resources: map[string]*tfbridge.ResourceInfo{
			// Upstream "id" is a computed-only Int64Attribute; Pulumi ids must be
			// strings. Coercible types just need a SchemaInfo.Type override (bridge
			// pkg/pf/tfbridge/ids.go stringifies the value at runtime) - no ComputeID.
			"unleash_service_account": {
				Fields: map[string]*tfbridge.SchemaInfo{
					"id": {Type: "string"},
				},
			},
			"unleash_service_account_token": {
				Fields: map[string]*tfbridge.SchemaInfo{
					"id": {Type: "string"},
				},
			},
		},
		CSharp: &tfbridge.CSharpInfo{
			// Publish under the Pulumiverse NuGet namespace, not the official Pulumi.* namespace.
			RootNamespace: "Pulumiverse",
			// RespectSchemaVersion ensures the SDK is generated linking to the correct version of the provider.
			RespectSchemaVersion: true,
			// Use a wildcard import so NuGet will prefer the latest possible version.
			PackageReferences: map[string]string{
				"Pulumi": "3.*",
			},
		},
	}

	// MustComputeTokens maps all resources and datasources from the upstream provider into Pulumi.
	//
	// tokens.SingleModule puts every upstream item into your provider's main module.
	//
	// You shouldn't need to override anything, but if you do, use the [tfbridge.ProviderInfo.Resources]
	// and [tfbridge.ProviderInfo.DataSources].
	prov.MustComputeTokens(tokens.SingleModule("unleash_", mainMod,
		tokens.MakeStandard(mainPkg)))

	prov.SetAutonaming(255, "-")

	return prov
}
