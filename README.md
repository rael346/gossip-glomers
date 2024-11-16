# Gossip Glomers: Distributed System Challenges from Fly.io
[source](https://fly.io/dist-sys/)

## Challenge 1: Echo
```sh
cd cmd/echo/ && go build
maelstrom test -w echo --bin ./echo --node-count 1 --time-limit 10
```

## Challenge 2: Unique Id Generation
```sh
cd cmd/idgen/ && go idgen
maelstrom test -w unique-ids --bin ./idgen --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition
```
