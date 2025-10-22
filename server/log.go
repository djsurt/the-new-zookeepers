package main

import "slices"

type LogSquasher[Entry any, Snapshot any] func(prevSnapshot *Snapshot, entries []Entry) Snapshot

type Log[Entry any, Snapshot any] struct {
	squasher       LogSquasher[Entry, Snapshot]
	snapshot       *Snapshot
	realFirstIndex uint64
	values         []Entry
}

func (log *Log[Entry, Snapshot]) squashFirstN(n uint64) {
	toSquash := log.values[:n+1]
	newValues := log.values[n+1:]

	newSnapshot := log.squasher(log.snapshot, toSquash)

	log.snapshot = &newSnapshot
	log.values = newValues
	log.realFirstIndex += n
}

type LogEntryPredicate[Entry any] func(entry Entry) bool

func (log *Log[Entry, Snapshot]) SquashUntil(predicate LogEntryPredicate[Entry]) {
	firstNonSquashed := slices.IndexFunc(log.values, predicate)
	// none found
	if firstNonSquashed == -1 {
		return
	}
	// none to squash (b/c first element not squash candidate)
	if firstNonSquashed < 1 {
		return
	}

	log.squashFirstN(uint64(firstNonSquashed - 1))
}
