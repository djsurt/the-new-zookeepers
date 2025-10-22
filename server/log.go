package main

type LogSquasher[Entry any, Snapshot any] func(prevSnapshot *Snapshot, entries []Entry) Snapshot

type LogEntry[V any] struct {
	value V
}

type Log[Entry any, Snapshot any] struct {
	squasher       LogSquasher[Entry, Snapshot]
	snapshot       *Snapshot
	realFirstIndex uint64
	values         []Entry
}

func (log *Log[_Entry, _Snapshot]) SquashFirstN(n uint64) {
	toSquash := log.values[:n+1]
	newValues := log.values[n+1:]

	newSnapshot := log.squasher(log.snapshot, toSquash)

	log.snapshot = &newSnapshot
	log.values = newValues
	log.realFirstIndex += n
}
