
# Golang Telegram games bot

A Telegram bot written in Golang, with which you can play turn-based games with other users. The bot consists of 3 microservices that communicate via gRPC. I use Redis as a database. The bot is running on AWS using docker compose.


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

I am using microservice architecture. Currentlty there are 3 microservices: `bot_service`, `player_service` `leaderboard_service`, which comunicate through gRPC.

<p align="middle">
  <img src=files/schema.jpg width="90%" />\
</p>


# Install

The microservices are dockerized, so to run on your machine you can:

`clone` the repository:
```
git clone https://github.com/Tumypmyp/chess
cd chess
```
Update to your telegram bot token and run:
<!-- sed 's/word1/word2/g' input.file > output.file -->-->
```
sed 's/0123:TelegramTOKEN/new:bot_token_here/g' docker-compose-example.yaml > docker-compose.yaml
make
```

# Next steps

By doing this project with small steps, the goal shifted a little. The goal now is to make turn-based game wrapper using microservices.

* Divide the `player_service` in smaller microservices.


[bot]: https://t.me/TumypmypGamesBot
