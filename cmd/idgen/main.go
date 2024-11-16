package main

import (
	"encoding/json"
	"log"
	"math/rand/v2"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type MsgBody struct {
	Type string `json:"type"`
}

type ResBody struct {
	Type string `json:"type"`
	Id   int    `json:"id"`
}

func main() {
	n := maelstrom.NewNode()

	n.Handle("generate", func(msg maelstrom.Message) error {
		var body MsgBody

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		resBody := ResBody{
			Type: "generate_ok",
			Id:   rand.Int(),
		}

		return n.Reply(msg, resBody)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
