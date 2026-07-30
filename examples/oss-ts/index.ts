import * as unleash from "@pulumiverse/unleash";

const baseUrl = process.env.UNLEASH_URL ?? "http://localhost:4242";
const authorization = process.env.UNLEASH_AUTH_TOKEN ?? "*:*.unleash-insecure-admin-api-token";

const provider = new unleash.Provider("test", {
    baseUrl,
    authorization,
});

const apiToken = new unleash.ApiToken("test", {
    tokenName: "pulumi-example",
    type: "client",
    projects: ["default"],
    environment: "development",
}, { provider });

export const tokenName = apiToken.tokenName;
export const environment = apiToken.environment;
