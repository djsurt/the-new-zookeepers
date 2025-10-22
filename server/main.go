package main

type NodeId uint64
type term uint64

type Node[LE any, LS any] struct {
	// latest term this node has seen
	currentTerm term
	// candidate voted for this term
	// TODO(amgg): probably give this a better type than just a ptr to make it nullable
	votedFor *NodeId
	log      Log[LE, LS]
	// leader-specific data if this node is a leader, or nil otherwise
	leaderData *LeaderData
}

type LeaderData struct {
	nextIndex  map[NodeId]uint
	matchIndex map[NodeId]uint
}

func (n Node[_, _]) IsLeader() bool { return n.leaderData != nil }
