package main

import "fmt"

func main() {
	rLogic1 := raftLogic{}
	cluster1 := []int{1, 2}
	rLogic1.newRaftLogic(1, cluster1)
	rLogic1.becomingLeader()
	rLogic1.clientNewCommand("set x 42")
	rLogic1.updateFollower(2, 0)
	fmt.Println(rLogic1.outgoing)
}
