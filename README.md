# Gossip Glomers: Distributed System Challenges from Fly.io
[source](https://fly.io/dist-sys/)

## Challenge 1: Echo
```sh
go build -o build/echo cmd/echo/main.go
maelstrom test \
    -w echo \
    --bin ./build/echo \
    --node-count 1 \
    --time-limit 10 \
```

## Challenge 2: Unique Id Generation
```sh
go build -o build/idgen cmd/idgen/main.go
maelstrom test \
    -w unique-ids \
    --bin ./build/idgen \
    --time-limit 30 \
    --rate 1000 \
    --node-count 3 \
    --availability total \
    --nemesis partition \
```

## Challenge 3: Broadcast

### 3a: Single-Node Broadcast
```sh
go build -o build/broadcast cmd/broadcast/main.go
maelstrom test \
    -w broadcast \
    --bin ./build/broadcast \
    --node-count 1 \
    --time-limit 20 \
    --rate 10 \
```

## Side Notes
- [Having multiple binaries in a single project](https://ieftimov.com/posts/golang-package-multiple-binaries)
    - Mainly because I didn't want to scatter these challenges in different repos
- Multiline shell commands just means escaping the newline character (`\n`), hence the backslash `\` at the end of each shell 
