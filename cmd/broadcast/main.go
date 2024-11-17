package main

import (
	"encoding/json"
	"log"

	mapset "github.com/deckarep/golang-set/v2"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Topology map[string][]string

type ResBody struct {
	Type     string `json:"type"`
	Messages *[]int `json:"messages,omitempty"`
}

type ReqBody struct {
	// Client fields
	Type  string `json:"type"`
	MsgId *int   `json:"msg_id,omitempty"`
	// InReplyTo int    `json:"in_reply_to,omitempty"`

	// Broadcast
	Topology Topology `json:"topology,omitempty"`
	Message  int      `json:"message,omitempty"`
}

type state struct {
	store mapset.Set[int]
	topo  Topology
	node  *maelstrom.Node
}

func (s *state) handleBroadcast(msg maelstrom.Message) error {
	var body ReqBody

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	if s.store.ContainsOne(body.Message) {
		return nil
	}

	s.store.Add(body.Message)

	for _, dst := range s.topo[s.node.ID()] {
		if dst == msg.Src {
			continue
		}

		s.node.Send(dst, ReqBody{
			Type:    "broadcast",
			Message: body.Message,
		})
	}

	if body.MsgId != nil {
		return s.node.Reply(msg, ResBody{
			Type: "broadcast_ok",
		})
	}

	return nil
}

func (s *state) handleRead(msg maelstrom.Message) error {
	copyStore := s.store.ToSlice()

	return s.node.Reply(msg, ResBody{
		Type:     "read_ok",
		Messages: &copyStore,
	})
}

func (s *state) handleTopo(msg maelstrom.Message) error {
	var body ReqBody
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	s.topo = body.Topology

	return s.node.Reply(msg, ReqBody{
		Type: "topology_ok",
	})
}

func main() {
	s := state{
		store: mapset.NewSet[int](),
		node:  maelstrom.NewNode(),
	}

	s.node.Handle("broadcast", s.handleBroadcast)
	s.node.Handle("read", s.handleRead)
	s.node.Handle("topology", s.handleTopo)

	if err := s.node.Run(); err != nil {
		log.Fatal(err)
	}
}
