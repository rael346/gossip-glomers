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
A few optimization:
1. To avoid message loops (nodes keep passing the message around), each node can check if the message is already stored 
2. During inter-broadcasting (sending message between nodes and not with clients) 
    a. A node doesn't need to send to its neighbor if the message came from that neighbor
    b. A node doesn't need to send `broadcast_ok` response if the broadcast request came from a neighbor instead of a client (can be checked by the field `msg_id` in message body)

### 3c: Fault Tolerant Broadcast
```sh
make broadcast-c
```
- The main idea is to retry internal `broadcast` until the destination node sends back a `broadcast_ok` response. This makes the 2b optimization above obsolete since every `broadcast` requires a response now.
- A simple implementation is for every `broadcast` received from a client, send it to every neighbor, and if the node doesn't receive a response after say 1 second, resend the message again 
    - Note: you will need to send the `broadcast_ok` response after adding the message to the node and before sending the internal `broadcast`s since the client will timeout after a while 
    - This approach is not ideal when using `Node.RPC()` since the callback is saved per `msg_id`. So if a `broadcast` never receive a response then that callback will never be deleted.

## Side Notes
- [Having multiple binaries in a single project](https://ieftimov.com/posts/golang-package-multiple-binaries)
    - Mainly because I didn't want to scatter these challenges in different repos
- Multiline shell commands just means escaping the newline character (`\n`), hence the backslash `\` at the end of each shell 
- Using `omitempty` in Go for JSON marshal/unmarshal will remove the field if it is the default value for that field's type (0 for int, "" for string, etc). So to differentiate between a field that is actually empty from a field specifically sets as the default value, make that field a pointer and check for `nil` instead. 
