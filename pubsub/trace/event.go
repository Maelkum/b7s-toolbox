package trace

import (
	"encoding/hex"
	"time"

	pb "github.com/libp2p/go-libp2p-pubsub/pb"
	"github.com/libp2p/go-libp2p/core/peer"
)

func createEventRecord(ev *pb.TraceEvent) Event {

	out := Event{
		TypeID:        int(ev.GetType()),
		Type:          ev.GetType().String(),
		PeerID:        peerIDToString(ev.PeerID),
		Timestamp:     ev.GetTimestamp(),
		TimestampNice: time.Unix(0, ev.GetTimestamp()).String(),
	}

	if ev.PublishMessage != nil {
		rec := Publish{
			MessageID: hex.EncodeToString(ev.PublishMessage.MessageID),
			Topic:     valueOrZero(ev.PublishMessage.Topic),
		}
		out.Publish = &rec
	}

	if ev.RejectMessage != nil {

		rec := Reject{
			MessageID:    hex.EncodeToString(ev.RejectMessage.MessageID),
			ReceivedFrom: peerIDToString(ev.RejectMessage.ReceivedFrom),
			Reason:       valueOrZero(ev.RejectMessage.Reason),
			Topic:        valueOrZero(ev.RejectMessage.Topic),
		}
		out.Reject = &rec
	}

	if ev.DuplicateMessage != nil {

		rec := Duplicate{
			MessageID:    hex.EncodeToString(ev.DuplicateMessage.MessageID),
			ReceivedFrom: peerIDToString(ev.DuplicateMessage.ReceivedFrom),
			Topic:        valueOrZero(ev.DuplicateMessage.Topic),
		}

		out.Duplicate = &rec
	}

	if ev.DeliverMessage != nil {

		rec := Deliver{
			MessageID:    hex.EncodeToString(ev.DeliverMessage.MessageID),
			ReceivedFrom: peerIDToString(ev.DeliverMessage.ReceivedFrom),
			Topic:        valueOrZero(ev.DeliverMessage.Topic),
		}

		out.Deliver = &rec
	}

	if ev.AddPeer != nil {

		rec := PeerAdd{
			Peer:  peerIDToString(ev.AddPeer.PeerID),
			Proto: valueOrZero(ev.AddPeer.Proto),
		}

		out.PeerAdd = &rec
	}

	if ev.RemovePeer != nil {

		out.PeerRemove = &PeerRemove{Peer: peerIDToString(ev.RemovePeer.PeerID)}
	}

	if ev.Join != nil {

		rec := Join{
			Topic: valueOrZero(ev.Join.Topic),
		}

		out.Join = &rec
	}

	if ev.Leave != nil {
		out.Leave = &Leave{Topic: valueOrZero(ev.Leave.Topic)}
	}

	if ev.Graft != nil {

		rec := ConnGraft{
			Peer:  peerIDToString(ev.Graft.PeerID),
			Topic: valueOrZero(ev.Graft.Topic),
		}

		out.ConnGraft = &rec
	}

	if ev.Prune != nil {

		rec := ConnPrune{
			Peer:  peerIDToString(ev.Graft.PeerID),
			Topic: valueOrZero(ev.Graft.Topic),
		}

		out.ConnPrune = &rec
	}

	if ev.SendRPC != nil {

		rec := RPCSend{
			SendTo: peerIDToString(ev.SendRPC.SendTo),
			Meta:   createMetaObject(ev.SendRPC.Meta),
		}

		out.RPCSend = &rec
	}

	if ev.RecvRPC != nil {

		rec := RPCReceive{
			ReceivedFrom: peerIDToString(ev.RecvRPC.ReceivedFrom),
			Meta:         createMetaObject(ev.RecvRPC.Meta),
		}

		out.RPCReceive = &rec
	}

	if ev.DropRPC != nil {

		rec := RPCDrop{
			SendTo: peerIDToString(ev.DropRPC.SendTo),
			Meta:   createMetaObject(ev.DropRPC.Meta),
		}

		out.RPCDrop = &rec
	}

	return out
}

func valueOrZero[T any](v *T) T {
	if v == nil {
		var out T
		return out
	}

	return *v
}

// helper function for bytes that we know are peerID, to not parse it each and every time.
func peerIDToString(id []byte) string {
	peerID, _ := peer.IDFromBytes(id)
	return peerID.String()
}

func createMetaObject(meta *pb.TraceEvent_RPCMeta) Meta {

	msgs := make([]Message, len(meta.Messages))
	for i, msg := range meta.Messages {
		msgs[i] = Message{
			Topic:     valueOrZero(msg.Topic),
			MessageID: hex.EncodeToString(msg.MessageID),
		}
	}

	subs := make([]Subscription, len(meta.Subscription))
	for i, sub := range meta.Subscription {
		subs[i] = Subscription{
			Subscribed: valueOrZero(sub.Subscribe),
			Topic:      valueOrZero(sub.Topic),
		}
	}

	return Meta{
		Messages:      msgs,
		Subscriptions: subs,
	}

}
