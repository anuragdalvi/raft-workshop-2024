package main

import (
	"reflect"
	"testing"
)

func TestAppendEntry(t *testing.T) {

	// append_entries() adds entries and return True or False

	log := LogEntry{}
	logEntry := log.newLogEntry()

	got := log.appendEntries(1, 1, []LogEntry{{2, "X"}}, &logEntry)

	if got != false {
		t.Errorf("Test 1: Appending entries with gaps in Term  ")
	}

	// Try an append with a bad previous_term (should fail)
	got = log.appendEntries(0, 1, []LogEntry{{2, "X"}}, &logEntry)

	if got != false {
		t.Errorf("Test 2: Appending entries with bad Previous Term  ")
	}

	// Add some new entries (should work)
	got = log.appendEntries(0, -1, []LogEntry{{1, "X"}, {1, "Y"}}, &logEntry)
	if got == false {
		t.Errorf("Test 3: Appending entries should work ")
	}

	if reflect.DeepEqual(logEntry[1:], []LogEntry{{1, "X"}, {1, "Y"}}) == false {
		t.Errorf("Test 4: Appended Entries are not equal to new Entries")
	}

	// Do the operation twice (idempotency).  Log should be the same afterwards
	got = log.appendEntries(0, -1, []LogEntry{{1, "X"}, {1, "Y"}}, &logEntry)
	if got == false {
		t.Errorf("Test 5: Test for Idempotancy")
	}

	if reflect.DeepEqual(logEntry[1:], []LogEntry{{1, "X"}, {1, "Y"}}) == false {
		t.Errorf("Test 6: Appended Entries are not equal to new Entries")
	}

	// Append empty entries.   Should work as long as prev_term and prev_index match up
	got = log.appendEntries(1, 1, []LogEntry{}, &logEntry)

	if got == false {
		t.Errorf("Test 7: Empty Entry Successful")
	}

	if reflect.DeepEqual(logEntry[1:], []LogEntry{{1, "X"}, {1, "Y"}}) == false {
		t.Errorf("Test 8: No new entries appended Empty Entry Successful")
	}

	// Append Entry that matches an existing entry
	got = log.appendEntries(0, -1, []LogEntry{{1, "X"}}, &logEntry)

	if got == false {
		t.Errorf("Test 9: Append an entry matching exisiting entry")
	}

	if reflect.DeepEqual(logEntry[1:], []LogEntry{{1, "X"}, {1, "Y"}}) == false {
		t.Errorf("Test 10: No new entries appended Empty Entry Successful")
	}

	// Add a single entry that does NOT match an existing entry. Should delete remaining entries
	got = log.appendEntries(0, -1, []LogEntry{{2, "A"}}, &logEntry)

	if got == false {
		t.Errorf("Test 11: Changed Exisiting Logs")
	}

	if reflect.DeepEqual(logEntry[1:], []LogEntry{{2, "A"}}) == false {
		t.Errorf("Test 12: Logs have changed.")
	}

	// A more complex deletion check if there is a mismatch in the middle
	newLog := log.createLogWithEntries([]LogEntry{
		{1, "A"},
		{1, "B"},
		{2, "C"},
		{2, "D"},
		{3, "E"},
	})

	got = log.appendEntries(1, 1, []LogEntry{{1, "B"}, {4, "F"}}, &newLog)
	if got == false {
		t.Errorf("Test 13: Append Complex entries")
	}

	if reflect.DeepEqual(newLog[1:], []LogEntry{{1, "A"}, {1, "B"}, {4, "F"}}) == false {
		t.Errorf("Test 12: Complex Entries Confirmed.")
	}

}

func TestFigureSeven(t *testing.T) {
	log := LogEntry{}

	createLog := func(terms []int) []LogEntry {

		var logTerm []LogEntry
		for _, term := range terms {
			logTerm = append(logTerm, LogEntry{term, " "})
		}
		return log.createLogWithEntries(logTerm)
	}

	//try adding a new log entry at index 11, term=8 on all

	//(a) Should fail, there's a gap in the log at position 10
	logA := log.newLogEntry()
	logA = createLog([]int{1, 1, 1, 4, 4, 5, 5, 6, 6})

	got := log.appendEntries(10, 6, []LogEntry{{8, " "}}, &logA)
	if got != false {
		t.Errorf("LargeInputTest 1: Missing Term ")
	}

	// (b) Should fail. Major gap
	logB := log.newLogEntry()
	logB = createLog([]int{1, 1, 1, 4})

	got = log.appendEntries(10, 6, []LogEntry{{8, " "}}, &logB)
	if got != false {
		t.Errorf("LargeInputTest 2: Major Gap")
	}

	// (c) Should work. But last entry will get replaced.
	logC := log.newLogEntry()
	logC = createLog([]int{1, 1, 1, 4, 4, 5, 5, 6, 6, 6, 6})

	got = log.appendEntries(10, 6, []LogEntry{{8, " "}}, &logC)
	if got == false {
		t.Errorf("LargeInputTest 4: Last Entry Replaced ")
	}

	// (d) Should work, two entries get replaced.
	logD := log.newLogEntry()
	logD = createLog([]int{1, 1, 1, 4, 4, 5, 5, 6, 6, 6, 7, 7})

	got = log.appendEntries(10, 6, []LogEntry{{8, " "}}, &logD)
	if got == false {
		t.Errorf("LargeInputTest 4: Two Entries Replaced ")
	}

	// (f) Should fail. Terms don't match up right
	logF := log.newLogEntry()
	logF = createLog([]int{1, 1, 1, 2, 2, 2, 3, 3, 3, 3, 3})

	got = log.appendEntries(10, 6, []LogEntry{{8, " "}}, &logF)
	if got != false {
		t.Errorf("LargeInputTest 5: Terms Donot Match Upright")
	}
}
