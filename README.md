# Gossip Glomers: Distributed System Challenges from Fly.io
[source](https://fly.io/dist-sys/)

## Challenge 1: Echo
```sh
go build echo.go
./maelstrom/maelstrom test -w echo --bin ./echo --node-count 1 --time-limit 10
```
