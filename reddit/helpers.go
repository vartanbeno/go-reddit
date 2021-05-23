package reddit

import (
	"sync"
)

// OrderedMaxSet is intended to be able to check if things have been seen while
// expiring older entries that are unlikely to be seen again.
// This is to avoid memory issues in long-running streams.
type OrderedMaxSet struct {
	MaxSize int
	set     map[string]struct{}
	keys    []string
	mutex   *sync.Mutex
}

// NewOrderedMaxSet instantiates an OrderedMaxSet and returns it for downstream use.
func NewOrderedMaxSet(maxSize int) OrderedMaxSet {
	var mutex = &sync.Mutex{}
	orderedMaxSet := OrderedMaxSet{
		MaxSize: maxSize,
		set:     map[string]struct{}{},
		keys:    []string{},
		mutex:   mutex,
	}

	return orderedMaxSet
}

// Add accepts a string and inserts it into an OrderedMaxSet
func (s *OrderedMaxSet) Add(v string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, ok := s.set[v]
	if !ok {
		s.keys = append(s.keys, v)
		s.set[v] = struct{}{}
	}
	if len(s.keys) > s.MaxSize {
		for _, id := range s.keys[:len(s.keys)-s.MaxSize] {
			delete(s.set, id)
		}
		s.keys = s.keys[(len(s.keys) - s.MaxSize):]

	}
}

// Delete accepts a string and deletes it from OrderedMaxSet
func (s *OrderedMaxSet) Delete(v string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.set, v)
}

// Len returns the number of elements in OrderedMaxSet
func (s *OrderedMaxSet) Len() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return len(s.set)
}

// Exists accepts a string and determines if it is present in OrderedMaxSet
func (s *OrderedMaxSet) Exists(v string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, ok := s.set[v]
	return ok
}
