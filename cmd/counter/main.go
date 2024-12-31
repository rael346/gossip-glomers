package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type ResBody struct {
	Type  string `json:"type"`
	Value int    `json:"value,omitempty"`
}

type ReqBody struct {
	Type  string `json:"type"`
	Delta int    `json:"delta"`
}

type server struct {
	node *maelstrom.Node
	kv   *maelstrom.KV
	mu   sync.RWMutex
}

func (s *server) handleInit(msg maelstrom.Message) error {
	s.KVWrite(0)
	return nil
}

func (s *server) handleAdd(msg maelstrom.Message) error {
	var body ReqBody
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	s.mu.Lock()
	val, err := s.KVRead(s.node.ID())
	if err != nil {
		return err
	}

	err = s.KVWrite(val + body.Delta)
	if err != nil {
		return err
	}
	s.mu.Unlock()

	s.node.Reply(msg, ResBody{
		Type: "add_ok",
	})

	return nil
}

func (s *server) handleRead(msg maelstrom.Message) error {
	s.mu.RLock()
	// sum := 0
	// for _, id := range s.node.NodeIDs() {
	// 	val, err := s.KVRead(id)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	sum += val
	// }
	val, err := s.KVRead(s.node.ID())
	if err != nil {
		return err
	}
	s.mu.RUnlock()

	s.node.Reply(msg, ResBody{
		Type:  "read_ok",
		Value: val,
	})
	return nil
}

func (s *server) KVWrite(val int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return s.kv.Write(ctx, s.node.ID(), val)
}

func (s *server) KVRead(id string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return s.kv.ReadInt(ctx, s.node.ID())
}

func main() {
	node := maelstrom.NewNode()
	s := server{
		node: node,
		kv:   maelstrom.NewSeqKV(node),
	}

	s.node.Handle("init", s.handleInit)
	s.node.Handle("add", s.handleAdd)
	s.node.Handle("read", s.handleRead)

	if err := s.node.Run(); err != nil {
		log.Fatal(err)
	}
}
