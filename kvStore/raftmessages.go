package main

import (
	"encoding/json"
	"fmt"
)

// Composition used for types of Raft Messages

type RaftMessage struct {
	Sender int `json:"sender"`
	Dest   int `json:"destination"`
}

func (u RaftMessage) Encode() string {
	b, err := json.Marshal(&u)
	if err != nil {
		fmt.Println(err)
		return "Err"
	}
	return string(b)
}

// #GO: Struct decoded representation of struct
func (u RaftMessage) Decode(s string) *RaftMessage {
	var raftMessage RaftMessage
	err := json.Unmarshal([]byte(s), &raftMessage)
	if err != nil {
		fmt.Println("Error Decoding Message")
	}
	return &raftMessage
}

// Struct Embedding
type AppendEntriesRequest struct {
	RaftMessage `json:"raftMessage"`
	PrevIndex   int        `json:"prevIndex"`
	PrevTerm    int        `json:"prevTerm"`
	Entries     []LogEntry `json:"entries"`
}

// #GO: String representation of a struct in Golang.
func (u AppendEntriesRequest) String() string {
	return fmt.Sprintf("AppendEntriesRequest("+
		"sender=%d,"+
		"dest=%d,"+
		"prevIndex=%d,"+
		"prevTerm=%d,"+
		// %v format 'verb' for a slice in golang.
		// #GO: https://pkg.go.dev/fmt
		"entries=%v \n", u.RaftMessage.Sender, u.RaftMessage.Dest, u.PrevIndex, u.PrevTerm, u.Entries)
}

// #GO: JSON representation of a struct.
func (u AppendEntriesRequest) Encode() string {
	b, err := json.Marshal(&u)
	if err != nil {
		fmt.Println(err)
		return "Err"
	}
	return string(b)
}

// #GO: Struct decoded representation of struct
func (u AppendEntriesRequest) Decode(s string) *AppendEntriesRequest {
	var appendRequest AppendEntriesRequest
	err := json.Unmarshal([]byte(s), &appendRequest)
	if err != nil {
		fmt.Println("Error Decoding Message")
	}
	return &appendRequest
}

// Struct Embedding
type AppendEntriesResponse struct {
	RaftMessage
	Success    bool `json:"success"`
	MatchIndex int  `json:"matchIndex"`
}

// String Representation for testing.
func (u AppendEntriesResponse) String() string {
	return fmt.Sprintf("AppendEntriesResponse("+
		"sender=%d,"+
		"dest=%d,"+
		"success=%t,"+
		"match_index=%d \n)", u.RaftMessage.Sender, u.RaftMessage.Dest, u.Success, u.MatchIndex)
}

// I know I am duplicating this code but I will address it later by trying to fix it using Type Assertion.
// #GO: JSON representation of a struct.
func (u AppendEntriesResponse) Encode() string {
	b, err := json.Marshal(&u)
	if err != nil {
		fmt.Println(err)
		return "Err"
	}
	return string(b)
}

// #GO: Struct decoded representation of struct
func (u AppendEntriesResponse) Decode(s string) *AppendEntriesResponse {
	var appendResponse AppendEntriesResponse
	err := json.Unmarshal([]byte(s), &appendResponse)
	if err != nil {
		fmt.Println("Error Decoding Message")
	}
	return &appendResponse
}

// Struct Embedding
type ApplicationNewCommand struct {
	RaftMessage `json:"raftMessage"`
	Message     string `json:"message"`
}

func (u ApplicationNewCommand) String() string {
	return fmt.Sprintf("ApplicationNewCommand("+
		"sender=%d,"+
		"dest=%d,"+
		"message=%s \n", u.RaftMessage.Sender, u.RaftMessage.Dest, u.Message)
}

// I know I am duplicating this code but I will address it later by trying to fix it using Type Assertion.
// #GO: JSON representation of a struct.
func (u ApplicationNewCommand) Encode() string {
	b, err := json.Marshal(&u)
	if err != nil {
		fmt.Println(err)
		return "Err"
	}
	return string(b)
}

// #GO: Struct decoded representation of struct
func (u ApplicationNewCommand) Decode(s string) ApplicationNewCommand {
	var appCommand ApplicationNewCommand
	err := json.Unmarshal([]byte(s), &appCommand)
	if err != nil {
		fmt.Println("Error Decoding Message")
	}
	return appCommand
}
