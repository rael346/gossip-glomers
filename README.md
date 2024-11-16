# Gossip Glomers: Distributed System Challenges from Fly.io
[source](https://fly.io/dist-sys/)

## Challenge 1: Echo
```sh
make echo
```

## Challenge 2: Unique Id Generation
```sh
make idgen
```

## Challenge 3: Broadcast

### 3a: Single-Node Broadcast
```sh
make broadcast-a
```

### 3b: Multi-Node Broadcast
```sh
make broadcast-b
```

## Side Notes
- [Having multiple binaries in a single project](https://ieftimov.com/posts/golang-package-multiple-binaries)
    - Mainly because I didn't want to scatter these challenges in different repos
- Multiline shell commands just means escaping the newline character (`\n`), hence the backslash `\` at the end of each shell 
