package main

import (
	"encoding/json"
	"log"
	"sync"

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

type topology map[string][]string

type TopoReqBody struct {
	Type string   `json:"type"`
	Topo topology `json:"topology"`
}

type state struct {
	store []int
	topo  topology
	mu    sync.RWMutex
	n     *maelstrom.Node
}

func (s *state) handleBroadcast(msg maelstrom.Message) error {
	var body BroadcastReqBody

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	s.mu.Lock()
	s.store = append(s.store, body.Msg)
	s.mu.Unlock()

	for _, dst := range s.topo[s.n.ID()] {
		if dst == msg.Src || dst == s.n.ID() {
			continue
		}

		s.n.RPC(dst, body, func(msg maelstrom.Message) error {
			return nil
		})
	}

	return s.n.Reply(msg, BroadcastResBody{
		Type: "broadcast_ok",
	})
}

func (s *state) handleRead(msg maelstrom.Message) error {
	s.mu.RLock()
	copyStore := make([]int, len(s.store))
	copy(copyStore, s.store)
	s.mu.RUnlock()

	return s.n.Reply(msg, ReadResBody{
		Type: "read_ok",
		Msgs: copyStore,
	})
}

func (s *state) handleTopo(msg maelstrom.Message) error {
	var body TopoReqBody
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	s.topo = body.Topo

	return s.n.Reply(msg, BroadcastResBody{
		Type: "topology_ok",
	})
}

func main() {
	s := state{
		store: []int{},
		n:     maelstrom.NewNode(),
	}

	s.n.Handle("broadcast", s.handleBroadcast)
	s.n.Handle("read", s.handleRead)
	s.n.Handle("topology", s.handleTopo)

	if err := s.n.Run(); err != nil {
		log.Fatal(err)
	}
}
