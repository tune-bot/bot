# tune-bot discord
Discord bot to interface with tunebot and play user playlists in the voice channel. 

## Setup
Either use the installation script from the [infrastructure](https://www.github.com/tune-bot/infrastructure) repository, or take the following steps to run the bot locally. Both methods will require the Discord bot token which is only accessible from the Discord [Developer Portal](https://discord.com/developers), under `Applications > TuneBot > Bot > "Reset Token"`. The token can only be seen when it is first generated, so reinstalling the bot would require resetting the token. You must also ensure the [database](https://www.github.com/tune-bot/database) is running.

```
export DISCORD_TOKEN=<token>
git clone github.com/tune-bot/discord
cd discord
go run main.go
```