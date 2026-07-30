// Copyright 2024, Pulumi Corporation.  All rights reserved.
//go:build nodejs || all
// +build nodejs all

package examples

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcnetwork "github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	postgresNetworkAlias = "db"
	unleashAdminAPIToken = "*:*.unleash-insecure-admin-api-token"
	unleashExposedPort   = "4242/tcp"
	postgresExposedPort  = "5432/tcp"
)

func TestBasicTs(t *testing.T) {
	t.Skip("Skipping until the provider has been implemented")

	opts := getJSBaseOptions(t).With(integration.ProgramTestOptions{
		Dir: filepath.Join(getCwd(t), "basic-ts"),
	})

	integration.ProgramTest(t, &opts)
}

func TestOssTs(t *testing.T) {
	unleashURL := startUnleashForTest(t)

	test := getJSBaseOptions(t).With(integration.ProgramTestOptions{
		Dir: filepath.Join(getCwd(t), "oss-ts"),
		Env: []string{
			fmt.Sprintf("UNLEASH_URL=%s", unleashURL),
			fmt.Sprintf("UNLEASH_AUTH_TOKEN=%s", unleashAdminAPIToken),
		},
	})
	integration.ProgramTest(t, &test)
}

// startUnleashForTest brings up Postgres + Unleash in testcontainers, on a shared
// network, and returns the base URL of the Unleash server reachable from the host.
func startUnleashForTest(t *testing.T) string {
	t.Helper()
	ctx := context.Background()

	nw, err := tcnetwork.New(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, nw.Remove(context.Background()))
	})

	dbContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:16",
			ExposedPorts: []string{postgresExposedPort},
			Env: map[string]string{
				"POSTGRES_DB":               "unleash",
				"POSTGRES_HOST_AUTH_METHOD": "trust",
			},
			Networks: []string{nw.Name},
			NetworkAliases: map[string][]string{
				nw.Name: {postgresNetworkAlias},
			},
			WaitingFor: wait.ForExec([]string{
				"pg_isready", "--username=postgres", "--host=127.0.0.1", "--port=5432",
			}).WithStartupTimeout(1 * time.Minute),
		},
		Started: true,
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, dbContainer.Terminate(context.Background()))
	})

	unleashContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "unleashorg/unleash-server:latest",
			ExposedPorts: []string{unleashExposedPort},
			Env: map[string]string{
				"DATABASE_URL":          fmt.Sprintf("postgres://postgres:unleash@%s/unleash", postgresNetworkAlias),
				"DATABASE_SSL":          "false",
				"INIT_ADMIN_API_TOKENS": unleashAdminAPIToken,
			},
			Networks: []string{nw.Name},
			WaitingFor: wait.ForHTTP("/health").
				WithPort(unleashExposedPort).
				WithStartupTimeout(3 * time.Minute),
		},
		Started: true,
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, unleashContainer.Terminate(context.Background()))
	})

	host, err := unleashContainer.Host(ctx)
	require.NoError(t, err)
	mappedPort, err := unleashContainer.MappedPort(ctx, unleashExposedPort)
	require.NoError(t, err)

	return fmt.Sprintf("http://%s:%s", host, mappedPort.Port())
}
