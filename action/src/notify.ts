import * as core from "@actions/core";
import * as http from "@actions/http-client";

function isIdTokenAvailable(): boolean {
  const token = process.env["ACTIONS_ID_TOKEN_REQUEST_TOKEN"];
  const url = process.env["ACTIONS_ID_TOKEN_REQUEST_URL"];
  return token && url ? true : false;
}

interface NotifyParams {
  payload: string;
  teamId: string;
  channelId: string;
}

export async function notify(params: NotifyParams): Promise<void> {
  const defaultEndpoint = "https://gha-notify.shogo82148.com";

  if (!isIdTokenAvailable()) {
    core.setFailed(
      `OIDC provider is not available. please enable it. see https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/about-security-hardening-with-openid-connect`,
    );
    return;
  }

  const headers: { [name: string]: string } = {};
  const token = await core.getIDToken(defaultEndpoint);
  headers["Authorization"] = `Bearer ${token}`;

  const payload = JSON.parse(params.payload);
  payload.team = params.teamId;
  payload.channel = params.channelId;

  const client = new http.HttpClient("actions-notify-slack");
  await client.postJson(defaultEndpoint + "/notify", payload, headers);
}
