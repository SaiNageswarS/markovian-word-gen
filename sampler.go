package main

import (
	"math/rand"
	"time"
)

// Sampler provides weighted random sampling functionality
type Sampler[T comparable] struct {
	itemToIdxMap map[T]int  // map of item to index in the items slice/cumulative frequency slice
	items        []T        // slice of items
	cumFreq      []int      // cumulative frequency of items
	random       *rand.Rand // random number generator
}

// NewSampler creates a new Sampler with the specified capacity
func NewSampler[T comparable](cap int) *Sampler[T] {
	return &Sampler[T]{
		itemToIdxMap: make(map[T]int),
		items:        make([]T, 0, cap),
		cumFreq:      make([]int, cap),
		random:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Add adds an item with the specified frequency to the sampler
func (s *Sampler[T]) Add(item T, freq int) {
	idx, exists := s.itemToIdxMap[item]

	if !exists {
		// New item
		idx = len(s.items) // add the item to the end of the items slice
		s.items = append(s.items, item)
		s.itemToIdxMap[item] = idx

		var prev int
		if idx > 0 {
			prev = s.cumFreq[idx-1]
		}
		s.cumFreq[idx] = prev + freq
	} else {
		// Existing item - update cumulative frequencies for all items after the current item
		for i := idx; i < len(s.items); i++ {
			s.cumFreq[i] += freq
		}
	}
}

// Sample returns a randomly sampled item based on weights, or nil if empty
func (s *Sampler[T]) Sample() *T {
	if len(s.items) == 0 {
		return nil
	}

	total := s.cumFreq[len(s.items)-1]
	rng := s.random.Intn(total)
	rngIdx := s.bsearch(0, len(s.items)-1, rng)
	return &s.items[rngIdx]
}

// bsearch performs binary search to find the index for the given target
func (s *Sampler[T]) bsearch(l, r, target int) int {
	if l == r {
		return l
	}

	m := l + (r-l)/2
	if s.cumFreq[m] > target {
		// keep m; answer is in [l, m]
		return s.bsearch(l, m, target)
	} else {
		// cumFreq[m] <= target; answer must be > m
		return s.bsearch(m+1, r, target)
	}
}
