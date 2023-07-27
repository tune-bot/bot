# tune-bot bot
Discord bot to interface with tunebot and play user playlists in the voice channel. 

## How to Build Locally
This requires the Discord bot token which is only accessible from the Discord [Developer Portal](https://discord.com/developers), under `Applications > TuneBot > Bot > "Reset Token"`. The token can only be seen when it is first generated, and it available as the environment variable DISCORD_TOKEN. You must also ensure the [core](https://www.github.com/tune-bot/core) is running for local database development, or have a valid host and credentials to the production tunebot server.

```
export DISCORD_TOKEN=<token>
git clone github.com/tune-bot/bot
cd bot
go run .
```