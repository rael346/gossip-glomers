package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type BroadcastReqBody struct {
	Type string `json:"type"`
	Msg  int    `json:"message"`
}

type BroadcastResBody struct {
	Type string `json:"type"`
}

type ReadResBody struct {
	Type string `json:"type"`
	Msgs []int  `json:"messages"`
}

type TopoReqBody struct {
	Type     string              `json:"type"`
	Topology map[string][]string `json:"topology"`
}

var store []int = []int{}

func main() {
	n := maelstrom.NewNode()

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body BroadcastReqBody

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		store = append(store, body.Msg)

		resBody := BroadcastResBody{
			Type: "broadcast_ok",
		}

		return n.Reply(msg, resBody)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		resBody := ReadResBody{
			Type: "read_ok",
			Msgs: store,
		}

		return n.Reply(msg, resBody)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		var body TopoReqBody

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		resBody := BroadcastResBody{
			Type: "topology_ok",
		}

		return n.Reply(msg, resBody)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
