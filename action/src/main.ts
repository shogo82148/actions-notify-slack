import * as core from "@actions/core";
import { notify } from "./notify";

async function run(): Promise<void> {
  try {
    await notify({
      payload: core.getInput("payload"),
      teamId: core.getInput("team-id"),
      channelId: core.getInput("channel-id"),
    });
  } catch (error) {
    if (error instanceof Error) {
      core.setFailed(error);
    } else {
      core.setFailed(`${error}`);
    }
  }
}

void run();
