package main

import (
	"encoding/json"
	"log"
	"sync"

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
	store map[int]struct{}
	topo  Topology
	mu    sync.RWMutex
	n     *maelstrom.Node
}

func (s *state) handleBroadcast(msg maelstrom.Message) error {
	var body ReqBody

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	s.mu.RLock()
	_, isMessageInStore := s.store[body.Message]
	s.mu.RUnlock()
	if isMessageInStore {
		return nil
	}

	s.mu.Lock()
	s.store[body.Message] = struct{}{}
	s.mu.Unlock()

	for _, dst := range s.topo[s.n.ID()] {
		if dst == msg.Src {
			continue
		}

		s.n.Send(dst, ReqBody{
			Type:    "broadcast",
			Message: body.Message,
		})
	}

	if body.MsgId != nil {
		return s.n.Reply(msg, ResBody{
			Type: "broadcast_ok",
		})
	}

	return nil
}

func (s *state) handleRead(msg maelstrom.Message) error {
	s.mu.RLock()
	copyStore := make([]int, 0, len(s.store))
	for k := range s.store {
		copyStore = append(copyStore, k)
	}
	s.mu.RUnlock()

	return s.n.Reply(msg, ResBody{
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

	return s.n.Reply(msg, ReqBody{
		Type: "topology_ok",
	})
}

func main() {
	s := state{
		store: map[int]struct{}{},
		n:     maelstrom.NewNode(),
	}

	s.n.Handle("broadcast", s.handleBroadcast)
	s.n.Handle("read", s.handleRead)
	s.n.Handle("topology", s.handleTopo)

	if err := s.n.Run(); err != nil {
		log.Fatal(err)
	}
}
