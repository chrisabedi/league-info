#League-Info Discord Bot

This is a small discord bot to hit the League of Legends API and return relevant information to the Channel requested

Currently two commands recognized are

`!ping`

which will return !Pong

and

`!puuid <GameName>#<tagLine>`

which will return the PUUID for the given RiotID and Tagline provided

and is *only* set up to hit the americas region of league of legends api.


There is heavy set up to get the needed API keys for discord and League of Legends developer api, but once you are able to get a bot token created from discord, Invite the Bot to your server with correct bot permissions, and finally, a developer token from League of Legends, the bot should work to get a summoners PUUID (player unique Identifier) in order to get more sophisticated results like most recent matches and performance information


This is a work in progress and still needs a good amount of work but comes with a league api client and a discord bot interface.

The ideal state is to get the last matches history for the player GameName / taglines Unique PUUID, but the JSON is large so that is the next step for this particular project

If you'd like to make a PR and suggest changes, feel free, I'm doing this as a nice to have discord bot in order to allow my friends easy access to quickly view there most recent matches in discords UI.

