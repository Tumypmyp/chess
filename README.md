
# Golang Telegram games bot

With this bot you can play Tic-Tac-Toe with other Telegram users.


# Usage
Bot is running on [@TumypmypGamesBot][bot]

`/newgame @player2` - starts game with *@player2*

A game could be played through bot or in a group chat.

`/leaderboard` - returns list of best players

<p align="middle">
  <img src=files/usage1.jpg width="40%" />
  <img src=files/usage2.jpg width="40%" /> 
</p>

# Description

Goal for this pet-project is to create a chess game in Telegram, explore tools such are `golang`, `redis`, `microservices`.

# Structure

I am using microservice architecture. Currentlty there are 2 microservices: `bot_service` and `leaderboard_service`, which comunicate through gRPC.

# Install

The microservices are dockerized, so to run on your machine you can:

`clone` the repository:
```
git clone https://github.com/Tumypmyp/chess
cd chess
```
Update to your telegram bot token and run:
<!-- sed 's/word1/word2/g' input.file > output.file -->
<!-- docker build -t dependencies -f ./dependencies.Dockerfile . -->
```
sed 's/0123:TelegramTOKEN/new:bot_token_here/g' docker-compose-example.yaml > docker-compose.yaml
make
```

# Next steps

By doing this project with small steps, the goal shifted a little. The goal now is to make turn-based game wrapper using microservices.

* Divide the `bot_service` in smaller microservices.
* Consider other asynchronous communication patterns
* Update game interface


[bot]: https://t.me/TumypmypGamesBot
