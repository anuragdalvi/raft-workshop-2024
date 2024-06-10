package main

type LogEntry struct {
	Term    int
	Command string
}

// The Raft paper assumes 1-based indexing for everything.
// Initialize LogEntry with a dummy Entry
// To avoid edge cases of having no entries
func (*LogEntry) newLogEntry() []LogEntry {
	return []LogEntry{{-1, " "}}
}

func (*LogEntry) getLength(log *[]LogEntry) int {
	return len(*log)
}

// This operation is performed by a leader only.
func (*LogEntry) appendNewCommand(term int, command string, log *[]LogEntry) {
	*log = append(*log, LogEntry{term, command})
}

// prevIndex argument specifies the position in the log after which the new entries go
// The prevTerm argument specifies the term value of the log entry at position prevIndex.
// entries is a list of zero
func (l *LogEntry) appendEntries(prevIndex int, prevTerm int, entries []LogEntry, log *[]LogEntry) bool {

	// Log is not allowed to have gaps/holes.
	if prevIndex >= len(*log) {
		return false
	}
	//fmt.Println((*log)[prevIndex].Term)
	//fmt.Println(prevTerm)

	// log has to link together properly.
	// The previous term has to match up with the log
	if (*log)[prevIndex].Term != prevTerm {
		return false
	}

	// Tricky case that requires careful reading of paragraph 5.3 and figure 2.
	// "If an existing entry conflicts with a new one (same index
	// but different terms), delete the existing entry and all that
	// follow it (§5.3)"
	for i, entry := range entries {
		n := prevIndex + i + 1
		if n >= len(*log) {
			break
		}
		if (*log)[n].Term != entry.Term {
			*log = (*log)[:n]
			break
		}
	}

	initialIndex := prevIndex + 1
	var finalEntries []LogEntry
	// Check if entries are to be inserted in the middle
	finalIndex := initialIndex + len(entries)
	if finalIndex < len(*log) {
		finalEntries = (*log)[finalIndex:]
	}

	*log = append((*log)[:initialIndex], append(entries, finalEntries...)...)

	return true
}

// Helper function to create log entries.
func (l *LogEntry) createLogWithEntries(entries []LogEntry) []LogEntry {
	newLog := l.newLogEntry()
	newLog = append(newLog, entries...)
	return newLog
}
