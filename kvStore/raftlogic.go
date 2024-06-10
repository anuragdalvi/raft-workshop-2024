package main

import "fmt"

// This implements the logic for a single Raft server.

type raftLogic struct {
	nodeNum     int
	cluster     []int
	currentTerm int
	role        string
	lObject     LogEntry
	log         []LogEntry
	nextIndex   map[int]int
	outgoing    []string
}

func (r *raftLogic) newRaftLogic(nodeNum int, cluster []int) {
	r.nodeNum = nodeNum
	r.cluster = cluster
	r.currentTerm = 1
	r.role = "FOLLOWER"
	r.lObject = LogEntry{}
	r.log = r.lObject.newLogEntry()
	r.nextIndex = make(map[int]int)
	for _, s := range r.cluster {
		r.nextIndex[s] = 0
	}
}

// Creating a Dummy sending function for testing without the networking logic.
func (r *raftLogic) send(message string) []string {

	r.outgoing = append(r.outgoing, message)
	return r.outgoing
}

func (r *raftLogic) becomingLeader() {
	//Upon becoming leader, we assume that all followers look exactly like us.
	// "Initialized to leader last log index + 1" (Figure 2).
	for _, s := range r.cluster {
		r.nextIndex[s] = len(r.log)
	}
	r.role = "LEADER"
}

// Clients need to have the leader add new command to the end of their log.
func (r *raftLogic) clientNewCommand(command string) {
	if r.role == "LEADER" {
		r.lObject.appendNewCommand(r.currentTerm, command, &r.log)
	}
}

// Leader need to be able to update followers by sending an AppendEntries message
func (r *raftLogic) updateFollower(follower int, numEntries int) {
	if r.role == "LEADER" {
		// The leader is tracking the state of each following using a variable nextIndex.
		// nextIndex records the index at which the leader *thinks* the next log entry
		// should go on each follower.  This may or may not be accurate.   We'll find
		// out when the follower sends a response back.

		prevIndex := r.nextIndex[follower] - 1
		prevTerm := r.log[prevIndex].Term
		var entries []LogEntry

		// Potential for Off-ByOne error.
		//if numEntries == 0 {
		//	entries = r.log[prevIndex+1:]
		//} else {
		//	entries = r.log[prevIndex+1 : prevIndex+1+numEntries]
		//}
		entries = r.log[prevIndex:]

		appendEntryRequest := AppendEntriesRequest{
			RaftMessage: RaftMessage{Sender: r.nodeNum, Dest: follower},
			PrevIndex:   prevIndex,
			PrevTerm:    prevTerm,
			Entries:     entries,
		}

		r.send(appendEntryRequest.Encode())

	}
}

func (r *raftLogic) updateAllFollowers() {
	for _, node := range r.cluster {
		if node != r.nodeNum {
			r.updateFollower(node, 0)
		}
	}
}

// All the Raft servers operate in response to messages that they receive.
// A single point of entry for all messages.
func (r *raftLogic) handleMessage(message interface{}) {
	switch v := message.(type) {
	case AppendEntriesRequest:
		r.handleAppendEntriesRequest(&v)
	case AppendEntriesResponse:
		r.handleAppendEntriesResponse(&v)
	default:
		fmt.Println("Error")
	}

}

func (r *raftLogic) handleAppendEntriesRequest(msg *AppendEntriesRequest) {

	success := r.lObject.appendEntries(msg.PrevIndex, msg.PrevTerm, msg.Entries, &r.log)

	appendEntryResponse := AppendEntriesResponse{
		RaftMessage: RaftMessage{Sender: r.nodeNum, Dest: msg.RaftMessage.Sender},
		Success:     success,
		MatchIndex:  msg.PrevIndex + len(msg.Entries) - 1,
	}
	r.send(appendEntryResponse.Encode())
}

func (r *raftLogic) handleAppendEntriesResponse(msg *AppendEntriesResponse) {
	//fmt.Println(r.nextIndex[msg.Sender])
	if msg.Success {
		r.nextIndex[msg.Sender] = msg.MatchIndex + 1
	} else {
		r.nextIndex[msg.Sender] -= 1
	}
}
