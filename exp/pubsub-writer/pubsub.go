package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	b7shost "github.com/blocklessnetwork/b7s/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"golang.org/x/exp/rand"
)

func subToTopic(h host.Host, topic string) (*pubsub.Topic, *pubsub.Subscription, error) {

	ps, err := pubsub.NewGossipSub(context.Background(), h)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create new gossipsub: %w", err)
	}

	th, err := ps.Join(topic)
	if err != nil {
		return nil, nil, fmt.Errorf("could not join topic: %w", err)
	}

	sub, err := th.Subscribe()
	if err != nil {
		return nil, nil, fmt.Errorf("could not subscribe to topic: %w", err)
	}

	return th, sub, nil
}

func publishMessages(h host.Host, topic *pubsub.Topic) error {

	for i := 0; ; i++ {

		msg := struct {
			Timestamp string `json:"timestamp"`
			Message   string `json:"message"`
			From      string `json:"from"`
			Rand      string `json:"rand"`
			Seq       uint   `json:"seq"`
		}{
			Timestamp: time.Now().String(),
			Message:   "new pubsub message!",
			From:      h.ID().String(),
			Rand:      getRandomString(12),
			Seq:       uint(i),
		}

		payload, _ := json.Marshal(msg)

		err := topic.Publish(context.Background(), payload)
		if err != nil {
			log.Error().Err(err).Uint("seq", uint(i)).Msg("could not publish message")
		}

		time.Sleep(publishTimeout)
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func getRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func publishMessagesb7s(h *b7shost.Host, topic *pubsub.Topic) error {
	return publishMessages(h.Host, topic)
}
