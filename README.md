
# Golang Telegram games bot

With this bot you can play Tic-Tac-Toe with other Telegram users.


# Usage
Bot is running on [@TumypmypGamesBot][bot]

```/newgame @player2``` - starts game with ```@player2```

A game could be played through bot or in a group chat.

<p align="middle">
  <img src=files/usage1.jpg width="40%" />
  <img src=files/usage2.jpg width="40%" /> 
</p>

```/leaderboard``` - returns list of best players

# Description

Idea for this pet project was to create a chess game in Telegram, and explore tools such are `golang`, `redis`, `microservices`.

# Structure

I am using microservice architecture. Currentlty there are 2 services: `bot_service` and `leaderboard_service`, which comunicate through gRPC.

# Next steps
 But by doing this project with small steps, the goal shifted a little. The goal now is much broader - to make turn-based game wrapper.



[bot]: https://t.me/TumypmypGamesBot
<!-- 
To start a server with docker on your machine run 
```
docker compose up
``` -->