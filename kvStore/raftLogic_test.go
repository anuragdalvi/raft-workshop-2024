package main

import (
	"reflect"
	"testing"
)

func TestBasicRaftLogic(t *testing.T) {
	// Logic1 Generate the outgoing message
	rLogic1 := raftLogic{}
	cluster1 := []int{1, 2}
	rLogic1.newRaftLogic(1, cluster1)

	// Logic 2: Follower accepts this message from logic1 essentially bypassing the network for this test.
	rLogic2 := raftLogic{}
	cluster2 := []int{1, 2}
	rLogic2.newRaftLogic(2, cluster2)

	rLogic1.becomingLeader()
	rLogic1.clientNewCommand("set x 42")
	rLogic1.updateFollower(2, 0)
	logic1Entry1 := AppendEntriesRequest{
		RaftMessage: RaftMessage{1, 2},
		PrevIndex:   0,
		PrevTerm:    -1,
		Entries:     []LogEntry{{-1, " "}, {1, "set x 42"}}}

	// Assert Contents of outgoing message for Raft Logic case 1
	if reflect.DeepEqual(rLogic1.outgoing[0], logic1Entry1.Encode()) == false {
		t.Errorf("Test 1: Failed to update follower")
	}

	// To Call Encode and Decode.
	logicEntryReq := AppendEntriesRequest{}
	logicEntryRes := AppendEntriesResponse{}

	originalNextIndex := rLogic1.nextIndex[2]
	// Now,  take the outgoing message created above and hand-deliver it to a follower
	//rLogic2.handleAppendEntriesRequest(logicEntryReq.Decode(rLogic1.outgoing[0]))

	rLogic2.handleMessage(*logicEntryReq.Decode(rLogic1.outgoing[0]))
	logic2Entry := AppendEntriesResponse{
		RaftMessage: RaftMessage{2, 1},
		Success:     true,
		MatchIndex:  1,
	}

	if reflect.DeepEqual(rLogic2.outgoing[0], logic2Entry.Encode()) == false {
		t.Errorf("Test 2: Follower update not successful")
	}

	rLogic1.handleMessage(*logicEntryRes.Decode(rLogic2.outgoing[0]))
	//A successful response should have advanced the next_index for the follower

	if rLogic1.nextIndex[2] != originalNextIndex+1 {
		t.Errorf("Test 3: Failed to Increment next Index for the follower.")
	}
}
