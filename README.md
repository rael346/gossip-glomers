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

### 3d: Efficient Broadcast, Part I
```sh
make broadcast-d
```
- The strategy from 3c still works, but needs a different topology configuration to get the required performance
    - `--topology grid`:
        ```
        ...

        :net {:all {:send-count 99026,
                     :recv-count 99026,
                     :msg-count 99026,
                     :msgs-per-op 55.789295},
               :clients {:send-count 3650, :recv-count 3650, :msg-count 3650},
               :servers {:send-count 95376,
                         :recv-count 95376,
                         :msg-count 95376,
                         :msgs-per-op 53.732956},
               :valid? true},
        ...

        :stable-latencies {
            0 0, 
            0.5 450, 
            0.95 673, 
            0.99 739, 
            1 807
        },
        ```
    - `--topology line`
        ```
        ...
         :net {:all {:send-count 46774,
                     :recv-count 46774,
                     :msg-count 46774,
                     :msgs-per-op 25.51773},
               :clients {:send-count 3766, :recv-count 3766, :msg-count 3766},
               :servers {:send-count 43008,
                         :recv-count 43008,
                         :msg-count 43008,
                         :msgs-per-op 23.463175},
               :valid? true},

        ...

        :stable-latencies {0 0, 
                           0.5 1561, 
                           0.95 2266, 
                           0.99 2365, 
                           1 2423},
        ```
    - `--topology tree4`
        ```
        ... 

         :net {:all {:send-count 45308,
                     :recv-count 45308,
                     :msg-count 45308,
                     :msgs-per-op 26.280743},
               :clients {:send-count 3548, :recv-count 3548, :msg-count 3548},
               :servers {:send-count 41760,
                         :recv-count 41760,
                         :msg-count 41760,
                         :msgs-per-op 24.222738},
               :valid? true},
        ...

        :stable-latencies {0 0,
                           0.5 377,
                           0.95 494,
                           0.99 506,
                           1 521},
        ```

# 3e: Efficient Broadcast, Part II
```sh
make broadcast-e
```
- The main goal for this challenge is to lower the `msgs-per-op` in exchange for higher latency throughout the system. 
- The simplest solution is to batch broadcast the messages instead of single broadcasts like the previous challenge.
    - **Goroutine for each neighboring node**: In the previous challenge, I got away with just iterating through the neighbors and sending the messsages to them. But to decrease the latency, I put each of the neighbor's broadcast (the broadcast with the neighbor as the destination) into a goroutine so they don't need to wait on each other (because of `time.Sleep()`)
        - Using `context` to timeout the request is very useful in implementing a timed request
        - Using Go's channel as a queue for the messages is very convenient since it acts as a FIFO queue data structure and is thread-safe by default
    - **Background scheduler for sending batch broadcast**: sending a broadcast after every single received message can only work with single-message broadcast. Instead, have a separate goroutine that runs every 300ms to send every message in the message channel to the neighboring node. This also ensure that a single goroutine can read from the channel, which is less concurrency bug to deal with. 

## Side Notes
- [Having multiple binaries in a single project](https://ieftimov.com/posts/golang-package-multiple-binaries)
    - Mainly because I didn't want to scatter these challenges in different repos
- Multiline shell commands just means escaping the newline character (`\n`), hence the backslash `\` at the end of each shell 
- Using `omitempty` in Go for JSON marshal/unmarshal will remove the field if it is the default value for that field's type (0 for int, "" for string, etc). So to differentiate between a field that is actually empty from a field specifically sets as the default value, make that field a pointer and check for `nil` instead. 


