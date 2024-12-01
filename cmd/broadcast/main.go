package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

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
	Message  *int     `json:"message,omitempty"`
	Messages *[]int   `json:"messages,omitempty"`
}

type state struct {
	store    mapset.Set[int]
	newStore map[int]struct{}
	mu       sync.RWMutex

	topo     Topology
	neighbor mapset.Set[string]
	node     *maelstrom.Node
	queue    chan int
}

func (s *state) handleBroadcast(msg maelstrom.Message) error {
	var body ReqBody
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	s.node.Reply(msg, ResBody{
		Type: "broadcast_ok",
	})

	log.Printf("DEBUG: src %s", msg.Src)
	if body.Messages != nil {
		log.Printf("DEBUG: batch received %v", *body.Messages)
		s.mu.Lock()
		for _, newMsg := range *body.Messages {
			if _, ok := s.newStore[newMsg]; ok {
				continue
			}

			s.newStore[newMsg] = struct{}{}
			s.queue <- newMsg
		}
		s.mu.Unlock()
	}

	if body.Message != nil {
		log.Printf("DEBUG: single received %d", *body.Message)
		s.mu.Lock()
		newMsg := *body.Message
		if _, ok := s.newStore[newMsg]; ok {
			s.mu.Unlock()
			return nil
		}

		log.Printf("DEBUG: add new val %d", newMsg)
		s.newStore[newMsg] = struct{}{}
		s.queue <- newMsg
		s.mu.Unlock()
	}

	return nil
}

func (s *state) batchBroadcast() {
	msgs := make([]int, 0, len(s.queue))
	for range len(s.queue) {
		msgs = append(msgs, <-s.queue)
	}
	if len(msgs) == 0 {
		return
	}

	log.Printf("DEBUG: batch broadcast %v", msgs)

	queue := s.neighbor.Clone()
	for _, dst := range queue.ToSlice() {
		go func() {
			for queue.ContainsOne(dst) {
				if err := s.RPC(dst, ReqBody{
					Type:     "broadcast",
					Messages: &msgs,
				}); err != nil {
					continue
				}
				queue.Remove(dst)
			}
		}()
	}
}

func (s *state) RPC(dst string, body any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_, err := s.node.SyncRPC(ctx, dst, body)
	return err
}

func (s *state) handleRead(msg maelstrom.Message) error {
	// copyStore := s.store.ToSlice()

	copyStore := make([]int, 0, len(s.newStore))
	s.mu.RLock()
	for val := range s.newStore {
		copyStore = append(copyStore, val)
	}
	s.mu.RUnlock()

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

	s.neighbor = mapset.NewSet(body.Topology[s.node.ID()]...)

	return s.node.Reply(msg, ReqBody{
		Type: "topology_ok",
	})
}

func main() {
	s := state{
		node:     maelstrom.NewNode(),
		queue:    make(chan int, 100),
		newStore: map[int]struct{}{},
		store:    mapset.NewSet[int](),
	}

	s.node.Handle("broadcast", s.handleBroadcast)
	s.node.Handle("read", s.handleRead)
	s.node.Handle("topology", s.handleTopo)

	go func() {
		ticker := time.NewTicker(250 * time.Millisecond)

		for range ticker.C {
			s.batchBroadcast()
		}
	}()

	if err := s.node.Run(); err != nil {
		log.Fatal(err)
	}
}
