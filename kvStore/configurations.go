package main

// An initialization function
// that creates a map
// This contains configuration of Raft cluster.

func getConfiguration() map[int][]string {
	return map[int][]string{
		0: []string{"localhost", "15000"},
		1: []string{"localhost", "16000"},
		2: []string{"localhost", "17000"},
		3: []string{"localhost", "18000"},
		4: []string{"localhost", "19000"},
	}
}
