package iavl

import (
	"fmt"
	"runtime"
)

type nodedbCache struct {
	firstVersion int64
	lastVersion  int64

	// dirtyRoots
	deletedRoots map[int64]struct{}
	dirtyRoots   map[int64][]byte

	// nodes
	deletedNodes map[string]struct{}
	dirtyNodes   map[string][]byte // node key -> node bytes

	// orphans, shared across different layers
	dirtyOrphans map[string][]byte // orphan key -> node hash

	// ignore fast node, may do not need if we commit by interval
}

func newNdbCache(firstVersion int64) *nodedbCache {
	c := &nodedbCache{
		firstVersion: firstVersion,
		lastVersion:  firstVersion,
		deletedRoots: make(map[int64]struct{}),
		dirtyRoots:   make(map[int64][]byte),
		deletedNodes: make(map[string]struct{}),
		dirtyNodes:   make(map[string][]byte),
		dirtyOrphans: make(map[string][]byte),
	}
	return c
}

func (l *nodedbCache) FirstVersion() int64 {
	trace()
	return l.firstVersion
}

func (l *nodedbCache) LastVersion() int64 {
	trace()
	return l.lastVersion
}

func (l *nodedbCache) Destroy() {
	trace()

	l.firstVersion = 0
	l.lastVersion = 0

	for k := range l.deletedNodes {
		delete(l.deletedNodes, k)
	}
	for k := range l.dirtyNodes {
		delete(l.dirtyNodes, k)
	}
	for k := range l.dirtyOrphans {
		delete(l.dirtyOrphans, k)
	}
}

func (l *nodedbCache) GetNode(key []byte) []byte {
	trace()

	if _, ok := l.deletedNodes[string(key)]; ok { // already deleted
		return nil
	}
	if v, ok := l.dirtyNodes[string(key)]; ok { // found it
		return v
	}
	return nil
}

func (l *nodedbCache) DeleteNode(nodeKey []byte) {
	trace()

	l.deletedNodes[string(nodeKey)] = struct{}{}
	delete(l.dirtyNodes, string(nodeKey))
}

func (l *nodedbCache) SaveNode(nodeKey, value []byte) {
	trace()

	l.dirtyNodes[string(nodeKey)] = value
}

func (l *nodedbCache) DeleteRoot(version int64) {
	trace()

	l.deletedRoots[version] = struct{}{}
	delete(l.dirtyRoots, version)
}

func (l *nodedbCache) DeleteRootsFrom(fromVersion int64) {
	trace()

	for i := fromVersion; i <= l.lastVersion; i++ {
		l.deletedRoots[i] = struct{}{}
		delete(l.dirtyRoots, i)
	}
}

// DeleteRootsRange deletes versions from an interval (not inclusive).
func (l *nodedbCache) DeleteRootsRange(fromVersion, toVersion int64) {
	trace()

	for i := fromVersion; i < toVersion; i++ {
		l.deletedRoots[i] = struct{}{}
		delete(l.dirtyRoots, i)
	}
}

func (l *nodedbCache) HasRoot(version int64) bool {
	trace()

	if _, ok := l.dirtyRoots[version]; ok {
		return true
	}

	return false
}

func (l *nodedbCache) GetRoot(version int64) []byte {
	trace()

	if v, ok := l.dirtyRoots[version]; ok {
		return v
	}

	return nil
}

func (l *nodedbCache) GetRoots() map[int64][]byte {
	trace()

	m := make(map[int64][]byte, 0)
	for k, v := range l.dirtyRoots {
		m[k] = v
	}
	return m
}

func (l *nodedbCache) SaveRoot(version int64, hash []byte) error {
	trace()

	if l.lastVersion < version {
		l.lastVersion = version
	}

	l.dirtyRoots[version] = hash

	return nil
}

func (l *nodedbCache) SaveOrphan(orphanKey, value []byte) {
	trace()

	l.dirtyOrphans[string(orphanKey)] = value
}

func (l *nodedbCache) GetOrphans() map[string][]byte {
	trace()

	return l.dirtyOrphans
}

func (l *nodedbCache) DeleteOrphan(orphanKey []byte) {
	trace()

	delete(l.dirtyOrphans, string(orphanKey))
}

// TODO: debug purpose, to be removed
func trace() {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	fmt.Printf("%s:%d %s\n", file, line, f.Name())
}
