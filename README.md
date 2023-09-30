# Notify Slack GitHub Action

Notify Slack GitHub Action is a simple GitHub Action to notify Slack.
The action uses OpenID Connect for authentication.
You don't need to manage any secrets, such as Incoming Webhook URLs or Bot User OAuth Token.

## Synopsis

1. Go to https://gha-notify.shogo82148.com/
2. Click the "Add to Slack" button
3. Invite @actions-notify-slack bot into your channel.
4. Grant the write permission to your repository by calling the slash command `/gha-notify ORG/REPO`.
5. Add the following step to your workflow.

```yaml
- uses: shogo82148/actions-notify-slack@v0
  with:
    team-id: T3G1HAY66 # replace it to your team id
    channel-id: C3GMGG162 # replace it to your channel id
    payload: '{"text": "hello world"}'
```
