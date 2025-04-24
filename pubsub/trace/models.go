package trace

type Event struct {
	TypeID        int         `json:"type_id,omitempty"`
	Type          string      `json:"type,omitempty"`
	PeerID        string      `json:"peer_id,omitempty"`
	Timestamp     int64       `json:"timestamp,omitempty"`
	TimestampNice string      `json:"timestamp_nice,omitempty"`
	Publish       *Publish    `json:"publish,omitempty"`
	Reject        *Reject     `json:"reject,omitempty"`
	Duplicate     *Duplicate  `json:"duplicate,omitempty"`
	Deliver       *Deliver    `json:"deliver,omitempty"`
	PeerAdd       *PeerAdd    `json:"peer_add,omitempty"`
	PeerRemove    *PeerRemove `json:"peer_remove,omitempty"`
	Join          *Join       `json:"join,omitempty"`
	Leave         *Leave      `json:"leave,omitempty"`
	ConnGraft     *ConnGraft  `json:"conn_graft,omitempty"`
	ConnPrune     *ConnPrune  `json:"conn_prune,omitempty"`
	RPCSend       *RPCSend    `json:"rpc_send,omitempty"`
	RPCReceive    *RPCReceive `json:"rpc_receive,omitempty"`
	RPCDrop       *RPCDrop    `json:"rpc_drop,omitempty"`
}

type Publish struct {
	MessageID string `json:"message_id,omitempty"`
	Topic     string `json:"topic,omitempty"`
}

type Reject struct {
	MessageID    string `json:"message_id,omitempty"`
	ReceivedFrom string `json:"received_from,omitempty"`
	Reason       string `json:"reason,omitempty"`
	Topic        string `json:"topic,omitempty"`
}

type Duplicate struct {
	MessageID    string `json:"message_id,omitempty"`
	ReceivedFrom string `json:"received_from,omitempty"`
	Topic        string `json:"topic,omitempty"`
}

type Deliver struct {
	MessageID    string `json:"message_id,omitempty"`
	ReceivedFrom string `json:"received_from,omitempty"`
	Topic        string `json:"topic,omitempty"`
}

type PeerAdd struct {
	Peer  string `json:"peer,omitempty"`
	Proto string `json:"proto,omitempty"`
}

type PeerRemove struct {
	Peer string `json:"peer,omitempty"`
}

type Join struct {
	Topic string `json:"topic,omitempty"`
}

type Leave struct {
	Topic string `json:"topic,omitempty"`
}

type ConnGraft struct {
	Peer  string `json:"peer,omitempty"`
	Topic string `json:"topic,omitempty"`
}

type ConnPrune struct {
	Peer  string `json:"peer,omitempty"`
	Topic string `json:"topic,omitempty"`
}

type RPCSend struct {
	SendTo string `json:"send_to,omitempty"`
	Meta   Meta   `json:"meta,omitempty"`
}

type RPCReceive struct {
	ReceivedFrom string `json:"received_from,omitempty"`
	Meta         Meta   `json:"meta,omitempty"`
}

type RPCDrop struct {
	SendTo string `json:"send_to,omitempty"`
	Meta   Meta   `json:"meta,omitempty"`
}

type Meta struct {
	Messages      []Message      `json:"messages,omitempty"`
	Subscriptions []Subscription `json:"subscription,omitempty"`
	// Control      Control        `json:"control,omitempty"`
}

type Message struct {
	MessageID string `json:"message_id,omitempty"`
	Topic     string `json:"topic,omitempty"`
}

type Subscription struct {
	Subscribed bool   `json:"subscribed,omitempty"`
	Topic      string `json:"topic,omitempty"`
}
