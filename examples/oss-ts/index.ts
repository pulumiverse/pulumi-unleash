import * as unleash from "@pulumi/unleash";

const provider = new unleash.Provider("test", {
    baseUrl: "http://localhost:4242",
    authorization: "*:*.unleash-insecure-admin-api-token",
});

const apiToken = new unleash.ApiToken("test", {
    tokenName: "pulumi-example",
    type: "client",
    projects: ["default"],
    environment: "development",
}, { provider });

export const tokenName = apiToken.tokenName;
export const environment = apiToken.environment;
